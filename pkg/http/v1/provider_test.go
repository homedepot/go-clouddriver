package v1_test

import (
	"bytes"
	"errors"
	"net/http"

	// . "github.com/homedepot/go-clouddriver/pkg/http/v1"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/jinzhu/gorm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	Describe("#CreateKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("provider not found"))
			uri = svr.URL + "/v1/kubernetes/providers"
			body.Write([]byte(payloadRequestKubernetesProviders))
			createRequest(http.MethodPost)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the request body is bad data", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte("dasdf[]dsf;;"))
				createRequest(http.MethodPost)
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadBadRequest)
			})
		})

		When("the ca data in the request is bad", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersBadCAData))
				createRequest(http.MethodPost)
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadErrorDecodingBase64)
			})
		})

		When("the provider already exists", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, nil)
			})

			It("returns status conflict", func() {
				Expect(res.StatusCode).To(Equal(http.StatusConflict))
				validateResponse(payloadConflictRequest)
			})
		})

		When("creating the kubernetes provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.CreateKubernetesProviderReturns(errors.New("error creating provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadErrorCreatingProvider)
			})
		})

		When("it succeeds", func() {
			It("returns status created", func() {
				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				validateResponse(payloadKubernetesProviderCreated)
			})
		})
	})

	Describe("#CreateOrReplaceKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/v1/kubernetes/providers"
			body.Write([]byte(payloadRequestKubernetesProviders))
			createRequest(http.MethodPut)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the request body is bad data", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte("dasdf[]dsf;;"))
				createRequest(http.MethodPut)
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadBadRequest)
			})
		})

		When("the ca data in the request is bad", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersBadCAData))
				createRequest(http.MethodPut)
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadErrorDecodingBase64)
			})
		})

		When("creating the kubernetes provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.CreateKubernetesProviderReturns(errors.New("error creating provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadErrorCreatingProvider)
			})
		})

		When("it succeeds", func() {
			When("the provider does not exist", func() {
				It("returns ok and the provider is created", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadKubernetesProviderCreated)
				})
			})

			When("the provider already exists", func() {
				BeforeEach(func() {
					fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, nil)
				})

				It("returns ok and the provider is replaced", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadKubernetesProviderCreated)
				})
			})
		})
	})

	Describe("#GetKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			testProvider := kubernetes.Provider{
				Name:   "test-name",
				Host:   "test-host",
				CAData: "dGVzdC1jYS1kYXRhCg==",
				Permissions: kubernetes.ProviderPermissions{
					Read:  []string{"gg_test"},
					Write: []string{"gg_test"},
				},
			}

			fakeSQLClient.GetKubernetesProviderAndPermissionsReturns(testProvider, nil)
			uri = svr.URL + "/v1/kubernetes/providers/test-name"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the record is not found", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderAndPermissionsReturns(kubernetes.Provider{}, nil)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotFound))
				validateResponse(payloadKubernetesProviderNotFound)
			})
		})

		When("getting the provider returns a generic error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderAndPermissionsReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesProviderGetGenericError)
			})
		})

		When("it succeeds", func() {
			It("returns ok and the provider", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadKubernetesProviderCreated)
			})
		})
	})

	Describe("#DeleteKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/v1/kubernetes/providers/test-name"
			createRequest(http.MethodDelete)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the record is not found", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, gorm.ErrRecordNotFound)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotFound))
				validateResponse(payloadKubernetesProviderNotFound)
			})
		})

		When("getting the provider returns a generic error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesProviderGetGenericError)
			})
		})

		When("deleting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.DeleteKubernetesProviderReturns(errors.New("error deleting provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesProviderDeleteGenericError)
			})
		})

		When("it succeeds", func() {
			It("returns status no content", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})
})
