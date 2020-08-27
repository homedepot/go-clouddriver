package core_test

import (
	// . "github.com/billiford/go-clouddriver/pkg/http/v0"

	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Credential", func() {
	Describe("#ListCredentials", func() {
		Context("expand query param is not set", func() {
			BeforeEach(func() {
				setup()
				uri = svr.URL + "/credentials"
				createRequest(http.MethodGet)
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
					{
						Name:        "provider1",
						Host:        "host1",
						CAData:      "caData1",
						BearerToken: "some.bearer.token",
						Permissions: kubernetes.ProviderPermissions{
							Read: []string{
								"gg_test",
							},
							Write: []string{
								"gg_test",
							},
						},
					},
					{
						Name:        "provider2",
						Host:        "host2",
						CAData:      "caData2",
						BearerToken: "some.bearer.token2",
						Permissions: kubernetes.ProviderPermissions{
							Read: []string{
								"gg_test2",
							},
							Write: []string{
								"gg_test2",
							},
						},
					},
				}, nil)
			})

			AfterEach(func() {
				teardown()
			})

			JustBeforeEach(func() {
				doRequest()
			})

			When("listing providers returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListKubernetesProvidersReturns(nil, errors.New("error listing providers"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing providers"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("listing read groups returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListReadGroupsByAccountNameReturns(nil, errors.New("error listing read groups"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing read groups"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("listing write groups returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListWriteGroupsByAccountNameReturns(nil, errors.New("error listing write groups"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing write groups"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentials)
				})
			})
		})

		Context("expand query param is set", func() {
			BeforeEach(func() {
				setup()
				uri = svr.URL + "/credentials?expand=true"
				createRequest(http.MethodGet)
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
					{
						Name:        "provider1",
						Host:        "host1",
						CAData:      "caData1",
						BearerToken: "some.bearer.token",
						Permissions: kubernetes.ProviderPermissions{
							Read: []string{
								"gg_test",
							},
							Write: []string{
								"gg_test",
							},
						},
					},
					{
						Name:        "provider2",
						Host:        "host2",
						CAData:      "caData2",
						BearerToken: "some.bearer.token2",
						Permissions: kubernetes.ProviderPermissions{
							Read: []string{
								"gg_test2",
							},
							Write: []string{
								"gg_test2",
							},
						},
					},
				}, nil)
				fakeKubeClient.ListReturns(&unstructured.UnstructuredList{
					Items: []unstructured.Unstructured{
						{
							Object: map[string]interface{}{
								"metadata": map[string]interface{}{
									"name": "namespace1",
								},
							},
						},
						{
							Object: map[string]interface{}{
								"metadata": map[string]interface{}{
									"name": "namespace2",
								},
							},
						},
					},
				}, nil)
				log.SetOutput(ioutil.Discard)
			})

			AfterEach(func() {
				teardown()
			})

			JustBeforeEach(func() {
				doRequest()
			})

			When("listing providers returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListKubernetesProvidersReturns(nil, errors.New("error listing providers"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing providers"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("listing read groups returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListReadGroupsByAccountNameReturns(nil, errors.New("error listing read groups"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing read groups"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("listing write groups returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListWriteGroupsByAccountNameReturns(nil, errors.New("error listing write groups"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing write groups"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("getting the kubernetes provider returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("decoding the ca data returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
						CAData: "{}",
					}, nil)
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("setting the dynamic client returns an error", func() {
				BeforeEach(func() {
					fakeKubeClient.SetDynamicClientForConfigReturns(errors.New("error setting the client"))
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("listing namespaces returns an error", func() {
				BeforeEach(func() {
					fakeKubeClient.ListReturns(nil, errors.New("error listing"))
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentialsExpandTrue)
				})
			})
		})
	})

	Describe("#GetAccountCredentials", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/credentials/test-account"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
				{
					Name:        "provider1",
					Host:        "host1",
					CAData:      "caData1",
					BearerToken: "some.bearer.token",
					Permissions: kubernetes.ProviderPermissions{
						Read: []string{
							"gg_test",
						},
						Write: []string{
							"gg_test",
						},
					},
				},
				{
					Name:        "provider2",
					Host:        "host2",
					CAData:      "caData2",
					BearerToken: "some.bearer.token2",
					Permissions: kubernetes.ProviderPermissions{
						Read: []string{
							"gg_test2",
						},
						Write: []string{
							"gg_test2",
						},
					},
				},
			}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing read groups returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListReadGroupsByAccountNameReturns(nil, errors.New("error listing read groups"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing read groups"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing write groups returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListWriteGroupsByAccountNameReturns(nil, errors.New("error listing write groups"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing write groups"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetAccountCredentials)
			})
		})
	})
})
