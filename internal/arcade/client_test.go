package arcade_test

import (
	"net/http"
	"time"

	. "github.com/homedepot/go-clouddriver/internal/arcade"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server   *ghttp.Server
		client   Client
		err      error
		token    string
		provider string
	)

	BeforeEach(func() {
		provider = "google"
		server = ghttp.NewServer()
		client = NewClient(server.URL())
		client.WithAPIKey("test-api-key")
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("#NewDefaultClient", func() {
		BeforeEach(func() {
			client = NewDefaultClient()
		})

		It("succeeds", func() {
		})
	})

	Describe("#Token", func() {
		JustBeforeEach(func() {
			token, err = client.Token(provider)
		})

		When("the uri is invalid", func() {
			BeforeEach(func() {
				client = NewClient("::haha")
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("the server is not reachable", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("the response is not 2XX", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error getting token: 500 Internal Server Error"))
			})
		})

		When("the server returns bad data", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.RespondWith(http.StatusOK, ";{["),
				)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid character ';' looking for beginning of value"))
			})
		})

		When("provider is rancher", func() {
			BeforeEach(func() {
				provider = "rancher"
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyHeaderKV("api-key", "test-api-key"),
					ghttp.VerifyRequest(http.MethodGet, "/tokens", "provider=rancher"),
					ghttp.RespondWith(http.StatusOK, `{"token":"some.bearer.token"}`),
				))
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(token).To(Equal("some.bearer.token"))
			})
		})

		When("provider is rancher and the 60 second short expiry has passed", func() {
			BeforeEach(func() {
				provider = "rancher"
				client = NewClient(server.URL())
				client.WithAPIKey("test-api-key")
				client.WithShortExpiration(2)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyHeaderKV("api-key", "test-api-key"),
					ghttp.VerifyRequest(http.MethodGet, "/tokens", "provider=rancher"),
					ghttp.RespondWith(http.StatusOK, `{"token":"some.bearer.token"}`),
				))

				// call to test the short expiration
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyHeaderKV("api-key", "test-api-key"),
					ghttp.VerifyRequest(http.MethodGet, "/tokens", "provider=rancher"),
					ghttp.RespondWith(http.StatusOK, `{"token":"new.bearer.token"}`),
				))
			})
			It("returns a new token", func() {
				Expect(err).To(BeNil())
				Expect(token).To(Equal("some.bearer.token"))

				time.Sleep(3 * time.Second)

				token, err = client.Token(provider)
				Expect(err).To(BeNil())
				Expect(token).To(Equal("new.bearer.token"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyHeaderKV("api-key", "test-api-key"),
					ghttp.VerifyRequest(http.MethodGet, "/tokens", "provider=google"),
					ghttp.RespondWith(http.StatusOK, `{"token":"some.bearer.token"}`),
				))
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(token).To(Equal("some.bearer.token"))
			})
		})
	})
})
