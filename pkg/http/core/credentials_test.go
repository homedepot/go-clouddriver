package core_test

import (
	// . "github.com/homedepot/go-clouddriver/pkg/http/v0"

	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
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
				fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns([]kubernetes.Provider{
					{
						Name:        "provider1",
						Host:        "host1",
						CAData:      "dGVzdAo=",
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
						CAData:      "dGVzdAo=",
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
					fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns(nil, errors.New("error listing providers"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing providers"))
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
				fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns([]kubernetes.Provider{
					{
						Name:        "provider1",
						Host:        "host1",
						CAData:      "dGVzdAo=",
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
						CAData:      "dGVzdAo=",
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
				fakeKubeClient.ListByGVRWithContextReturns(&unstructured.UnstructuredList{
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
					fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns(nil, errors.New("error listing providers"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(Equal("Internal Server Error"))
					Expect(ce.Message).To(Equal("error listing providers"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("decoding the ca data returns an error", func() {
				BeforeEach(func() {
					fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns([]kubernetes.Provider{
						{
							Name:        "provider1",
							Host:        "host1",
							CAData:      "{}",
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
							CAData:      "{}",
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

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("getting the gcloud token returns an error", func() {
				BeforeEach(func() {
					fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeArcadeClient.TokenCallCount()).To(Equal(2))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("creating the kube client returns an error", func() {
				BeforeEach(func() {
					fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeKubeController.NewClientCallCount()).To(Equal(2))
					validateResponse(payloadCredentialsExpandTrueNoNamespaces)
				})
			})

			When("listing namespaces returns an error", func() {
				BeforeEach(func() {
					fakeKubeClient.ListByGVRWithContextReturns(nil, errors.New("error listing"))
				})

				It("continues", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeKubeClient.ListByGVRWithContextCallCount()).To(Equal(2))
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
					CAData:      "dGVzdAo=",
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
					CAData:      "dGVzdAo=",
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
