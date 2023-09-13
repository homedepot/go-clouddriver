package front50_test

import (
	"encoding/json"
	"net/http"

	. "github.com/homedepot/go-clouddriver/internal/front50"
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

	Describe("#Project", func() {
		JustBeforeEach(func() {
			response, err = client.Project("fakeProjectID")
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
						"id": "048de9a7-7b57-4097-8444-e44682d9dcfc",
						"name": "spinnaker",
						"email": "david_m_rogers@homedepot.com",
						"config": {
							"pipelineConfigs": [
								{
									"application": "billy",
									"pipelineConfigId": "b1bb2476-388b-47b9-8730-9617cccfe458"
								}
							],
							"applications": [
								"smoketests"
							],
							"clusters": [
								{
									"account": "gae_np-te-cd-tools-np",
									"stack": "*",
									"detail": "*",
									"applications": null
								},
								{
									"account": "gke_np-te-cd-tools_us-central1_smoketests_sandbox-us-central1-np",
									"stack": "*",
									"detail": "*",
									"applications": null
								},
								{
									"account": "gke_np-te-cd-tools_us-central1_sandbox-us-central1-np",
									"stack": "*",
									"detail": "*",
									"applications": null
								}
							]
						},
						"updateTs": 1635962823717,
						"createTs": 1587655303067,
						"lastModifiedBy": "DXR05"
					}`

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/v2/projects/fakeProjectID"),
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
