package core_test

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artifacts", func() {
	Describe("#ListArtifactCredentials", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/credentials"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadArtifactCredentials)
			})
		})
	})

	Describe("#ListHelmArtifactAccountNames", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/account/helm-stable/names"
			createRequest(http.MethodGet)
			fakeHelmClient.GetIndexReturns(helm.Index{
				Entries: map[string][]helm.Resource{
					"prometheus-operator": {},
					"minecraft":           {},
				},
			}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the index returns an error", func() {
			BeforeEach(func() {
				fakeHelmClient.GetIndexReturns(helm.Index{}, errors.New("error getting index"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting index"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListHelmArtifactAccountNames)
			})
		})
	})

	Describe("#ListHelmArtifactAccountVersions", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/account/helm-stable/versions?artifactName=minecraft"
			createRequest(http.MethodGet)
			fakeHelmClient.GetIndexReturns(helm.Index{
				Entries: map[string][]helm.Resource{
					"prometheus-operator": {
						{
							Version: "2.0.0",
						},
						{
							Version: "2.1.0",
						},
					},
					"minecraft": {
						{
							Version: "1.0.0",
						},
						{
							Version: "1.1.0",
						},
					},
				}}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the index returns an error", func() {
			BeforeEach(func() {
				fakeHelmClient.GetIndexReturns(helm.Index{}, errors.New("error getting index"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting index"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListHelmArtifactAccountVersions)
			})
		})
	})

	Describe("#GetArtifact", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/fetch/"
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the request contains bad data", func() {
			BeforeEach(func() {
				setup()
				body.Write([]byte(";[]---"))
				createRequest(http.MethodPut)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Bad Request"))
				Expect(ce.Message).To(Equal("invalid character ';' looking for beginning of value"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the artifact is type helm/chart", func() {
			BeforeEach(func() {
				body.Write([]byte(payloadRequestFetchHelmArtifact))
				createRequest(http.MethodPut)
				fakeHelmClient.GetChartReturns([]byte("some-binary-data"), nil)
			})

			When("getting the chart returns an error", func() {
				BeforeEach(func() {
					fakeHelmClient.GetChartReturns(nil, errors.New("error getting chart"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error getting chart"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("some-binary-data")
				})
			})
		})

		Context("when the artifact is type embedded/base64", func() {
			BeforeEach(func() {
				body.Write([]byte(payloadRequestFetchBase64Artifact))
				createRequest(http.MethodPut)
			})

			When("the reference contains bad data", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(payloadRequestFetchBase64ArtifactBadReference))
					createRequest(http.MethodPut)
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("illegal base64 data at input byte 3"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("helloworld\n")
				})
			})
		})
	})
})
