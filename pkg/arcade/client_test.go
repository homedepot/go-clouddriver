package arcade_test

import (
	"net/http"

	. "github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server
		client Client
		err    error
		token  string
	)

	BeforeEach(func() {
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
			token, err = client.Token()
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

		When("it succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyHeaderKV("api-key", "test-api-key"),
					ghttp.RespondWith(http.StatusOK, `{"token":"some.bearer.token"}`),
				))
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(token).To(Equal("some.bearer.token"))
			})
		})
	})

	Describe("#Instance", func() {
		var c *gin.Context
		var c2 Client

		BeforeEach(func() {
			c = &gin.Context{}
			c.Set(ClientInstanceKey, client)
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				c2 = Instance(c)
			})

			It("succeeds", func() {
				Expect(c2).ToNot(BeNil())
			})
		})
	})
})
