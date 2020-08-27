package core_test

import (
	// . "github.com/billiford/go-clouddriver/pkg/http/v0"

	"errors"
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	Describe("#GetTask", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/task/task-id"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing the resources returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns(nil, errors.New("error listing resources"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Bad Request"))
				Expect(ce.Message).To(Equal("error listing resources"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("no resources are returned", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{}, nil)
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("setting the kube config returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.SetDynamicClientForConfigReturns(errors.New("bad config"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("bad config"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting provider"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("there is an error decoding the provider CA data", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					CAData: "@#$%",
				}, nil)
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("illegal base64 data at input byte 0"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the manifest returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting resource"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting resource"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
