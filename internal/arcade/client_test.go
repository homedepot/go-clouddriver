package arcade_test

import (
	"net/http"

	. "github.com/homedepot/go-clouddriver/internal/arcade"
	. "github.com/onsi/ginkgo/v2"
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
