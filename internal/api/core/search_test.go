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

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns an empty response", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearchEmptyResponse)
				Expect(fakeKubeClient.ListResourcesByKindAndNamespaceWithContextCallCount()).To(BeZero())
			})
		})

		Context("when the provider is namespace scoped", func() {
			var provider kubernetes.Provider

			BeforeEach(func() {
				d := "default"
				provider.Namespace = &d
				fakeSQLClient.GetKubernetesProviderReturns(provider, nil)
			})

			When("the namespace is incorrect", func() {
				BeforeEach(func() {
					d := "different-namespace"
					provider.Namespace = &d
					fakeSQLClient.GetKubernetesProviderReturns(provider, nil)
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

		It("succeeds", func() {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			validateResponse(payloadSearch)
		})
	})
})
