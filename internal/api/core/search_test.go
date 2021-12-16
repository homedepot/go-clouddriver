package core_test

import (
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Search", func() {
	Describe("#Search", func() {
		accountsHeader := ""

		BeforeEach(func() {
			setup()
			uri = svr.URL + "/search?pageSize=500&q=default&type=pod"
			accountsHeader = "account1"
			provider := kubernetes.Provider{
				Name: "account1",
			}
			fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{provider}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			createRequest(http.MethodGet)
			req.Header.Add("X-Spinnaker-Accounts", accountsHeader)
			doRequest()
		})

		When("kind and namespace are not provided", func() {
			BeforeEach(func() {
				uri = svr.URL + "/search?pageSize=500"
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("must provide query params 'q' to specify the namespace and 'type' to specify the kind"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("unsupported kinds", func() {
			When("kind is applications", func() {
				BeforeEach(func() {
					uri = svr.URL + "/search?pageSize=500&q=default&type=applications"
				})

				It("returns the default response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchDefault)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})

			When("kind is instances", func() {
				BeforeEach(func() {
					uri = svr.URL + "/search?pageSize=500&q=default&type=instances"
				})

				It("returns the default response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchDefault)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})

			When("kind is loadBalancers", func() {
				BeforeEach(func() {
					uri = svr.URL + "/search?pageSize=500&q=default&type=loadBalancers"
				})

				It("returns the default response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchDefault)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})

			When("kind is projects", func() {
				BeforeEach(func() {
					uri = svr.URL + "/search?pageSize=500&q=default&type=projects"
				})

				It("returns the default response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchDefault)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})

			When("kind is securityGroups", func() {
				BeforeEach(func() {
					uri = svr.URL + "/search?pageSize=500&q=default&type=securityGroups"
				})

				It("returns the default response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchDefault)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})
		})

		When("grabbing all providers returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns(nil, errors.New("error listing providers"))
			})

			It("returns internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("internal: error listing kubernetes providers: error listing providers"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when the provider is namespace scoped", func() {
			var provider kubernetes.Provider

			BeforeEach(func() {
				d := "default"
				provider.Name = "account1"
				provider.Namespace = &d
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{provider}, nil)
			})

			When("the namespace is incorrect", func() {
				BeforeEach(func() {
					d := "different-namespace"
					provider.Namespace = &d
					fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{provider}, nil)
				})

				It("returns an empty response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchEmptyResponse)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})

			When("the kind is cluster-scoped", func() {
				BeforeEach(func() {
					uri = svr.URL + "/search?pageSize=500&q=default&type=clusterRole"
				})

				It("returns an empty response", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadSearchEmptyResponse)
					Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
				})
			})

			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearch)
			})
		})

		When("there is an error listing resources", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourcesByKindAndNamespaceWithContextReturns(nil, errors.New("error listing resources"))
			})

			It("returns an empty response", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearchEmptyResponse)
			})
		})

		When("listing providers returns accounts the user does not have access to", func() {
			BeforeEach(func() {
				providers := []kubernetes.Provider{
					{
						Name: "account1",
					},
					{
						Name: "account2",
					},
					{
						Name: "account3",
					},
				}
				fakeSQLClient.ListKubernetesProvidersReturns(providers, nil)
			})

			It("filters the accounts", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearch)
			})
		})

		It("succeeds", func() {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			validateResponse(payloadSearch)
			Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(Equal(1))
		})
	})
})
