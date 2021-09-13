package core_test

import (
	"errors"
	"net/http"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
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
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error listing resources"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("no resources are returned", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{}, nil)
			})

			It("returns not found error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		When("creating the kube client returns an error", func() {
			BeforeEach(func() {
				fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("illegal base64 data at input byte 0"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting a token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting token"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("the task type is cleanup", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
					{
						AccountName: "test-account-name",
						TaskType:    clouddriver.TaskTypeCleanup,
					},
				}, nil)
			})

			It("does not call make calls to the server", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(Equal(0))
			})
		})

		When("the task type is noop", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
					{
						AccountName: "test-account-name",
						TaskType:    clouddriver.TaskTypeNoOp,
					},
				}, nil)
			})

			It("does not call make calls to the server", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(Equal(0))
			})
		})

		Context("when the task type is delete", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
					{
						AccountName: "test-account-name",
						TaskType:    clouddriver.TaskTypeDelete,
					},
				}, nil)
			})

			When("the server returns a not found error", func() {
				BeforeEach(func() {
					fakeKubeClient.GetReturns(nil, errors.New(`horizontalpodautoscalers.autoscaling "php-apache1-v008" not found`))
				})

				It("ignores the not found error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeKubeClient.GetCallCount()).To(Equal(1))
				})
			})

			When("the server returns a generic error", func() {
				BeforeEach(func() {
					fakeKubeClient.GetReturns(nil, errors.New(`generic error`))
				})

				It("ignores the not found error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					Expect(fakeKubeClient.GetCallCount()).To(Equal(1))
				})
			})
		})

		When("getting the manifest returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting resource"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
