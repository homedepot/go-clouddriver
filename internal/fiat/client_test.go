package fiat_test

import (
	"encoding/json"
	"net/http"

	. "github.com/homedepot/go-clouddriver/internal/fiat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server                     *ghttp.Server
		client                     Client
		err                        error
		response, expectedResponse Response
		fakeResponse               string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL())
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

	Describe("#Authorize", func() {
		JustBeforeEach(func() {
			response, err = client.Authorize("fakeAccount")
		})

		When("the uri is invalid", func() {
			BeforeEach(func() {
				client = NewClient("::haha")
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(response).To(Equal(Response{}))
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
				Expect(err.Error()).To(Equal("user authorization error: 500 Internal Server Error"))
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
				fakeResponse = `{
					"name" : "test_group",
					"accounts" : [ {
					  "name" : "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
					  "authorizations" : [ "READ", "WRITE", "EXECUTE", "CREATE" ]
					}, {
					  "name" : "spin-cluster-account",
					  "authorizations" : [ "READ", "WRITE", "EXECUTE", "CREATE" ]
					} ],
					"serviceAccounts" : [ {
					  "name" : "gg_cloud_gcp_spinnaker_admins_member",
					  "memberOf" : [ "gg_cloud_gcp_spinnaker_admins" ]
					} ],
					"roles" : [ {
					  "name" : "gg_cloud_gcp_spinnaker_admins",
					  "source" : "EXTERNAL"
					}, {
					  "name" : "test_group",
					  "source" : "EXTERNAL"
					} ],
					"buildServices" : [ ],
					"extensionResources" : { },
					"admin" : false,
					"legacyFallback" : false,
					"allowAccessToUnknownApplications" : false
				  }
				  `

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/authorize/fakeAccount"),
					ghttp.RespondWith(http.StatusOK, fakeResponse),
				))
			})

			It("succeeds", func() {
				byt := []byte(fakeResponse)
				err := json.Unmarshal(byt, &expectedResponse)
				Expect(err).To(BeNil())
				Expect(expectedResponse).To(Equal(response))
			})
		})
	})
})
