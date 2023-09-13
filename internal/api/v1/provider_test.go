package v1_test

import (
	"bytes"
	"errors"
	"net/http"

	// . "github.com/homedepot/go-clouddriver/internal/api/v1"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	Describe("#CreateKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("provider not found"))
			fakeArcadeClient.TokenReturns("fake-token", nil)
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

		When("the token provider is invalid", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("unsupported token provider"))
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadErrorGettingToken)
			})
		})

		When("the a write permission group is not a read permission group", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersMissingReadGroup))
				createRequest(http.MethodPost)
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadErrorMissingReadGroup)
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

		When("the namespace is empty string", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersEmptyNamespace))
				createRequest(http.MethodPost)
			})

			It("returns ok and the namespace is nil and namespace array is nil", func() {
				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				validateResponse(payloadKubernetesProviderCreatedNoNamespace)
			})
		})

		When("namespace and namespaces are empty", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersEmptyNamespace))
				createRequest(http.MethodPost)
			})

			It("returns ok and namespace is nil and namespaces is an empty array", func() {
				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				validateResponse(payloadKubernetesProviderCreatedNoNamespace)
			})
		})

		When("the deprecated namespace property is provided", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProviderDeprecatedNamespace))
				createRequest(http.MethodPost)
			})

			It("returns ok with the namespace as part of the namespaces property", func() {
				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				validateResponse(payloadKubernetesProviderCreatedWithDeprecatedNamespace)
			})
		})

		When("namespaces are provided", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersMultipleNamespaces))
				createRequest(http.MethodPost)
			})

			It("returns ok and the namespaces", func() {
				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				validateResponse(payloadRequestKubernetesProvidersMultipleNamespaces)
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

		When("the token provider is invalid", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("unsupported token provider"))
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadErrorGettingToken)
			})
		})

		When("the a write permission group is not a read permission group", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesProvidersMissingReadGroup))
				createRequest(http.MethodPut)
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				validateResponse(payloadErrorMissingReadGroup)
			})
		})

		When("deleting the kubernetes provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.DeleteKubernetesProviderReturns(errors.New("error deleting provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesProviderDeleteGenericError)
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

			When("the namespace is empty string", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(payloadRequestKubernetesProvidersEmptyNamespace))
					createRequest(http.MethodPut)
				})

				It("returns ok and the namespace is nil", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadKubernetesProviderCreatedNoNamespace)
				})
			})

			When("the namespace is set on a provider with no namespaces (updating deprecated field)", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(payloadRequestKubernetesProviderDeprecatedNamespace))
					createRequest(http.MethodPut)
				})

				It("returns ok and the namespace is nil and namespaces including the provided namespace value", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadKubernetesProviderCreatedWithDeprecatedNamespace)
				})
			})

			When("the namespaces field is updated", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(payloadRequestKubernetesProvidersMultipleNamespaces))
					createRequest(http.MethodPut)
				})

				It("returns ok and the namespaces are updated to the ones passed in", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadRequestKubernetesProvidersMultipleNamespaces)
				})
			})

			When("the namespace is set on a provider with existing namespaces", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(payloadRequestKubernetesUpdateProvidersExistingNamespaces))
					createRequest(http.MethodPut)
				})

				It("returns ok and the namespaces are updated to the one passed in for the namespace field", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadKubernetesProviderCreatedWithDeprecatedNamespace)
				})
			})
		})
	})

	Describe("#GetKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			testProvider := kubernetes.Provider{
				Name:       "test-name",
				Host:       "test-host",
				Namespaces: []string{"ns1", "ns2"},
				CAData:     "dGVzdC1jYS1kYXRhCg==",
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

	Describe("#ListKubernetesProvider", func() {
		BeforeEach(func() {
			setup()
			namespace := "test-namespace"
			testProviders := []kubernetes.Provider{
				{
					Name:   "test-name1",
					Host:   "test-host1",
					CAData: "dGVzdC1jYS1kYXRhCg==",
					Permissions: kubernetes.ProviderPermissions{
						Read:  []string{"gg_test1"},
						Write: []string{"gg_test1"},
					},
				},
				{
					Name:      "test-name2",
					Host:      "test-host2",
					CAData:    "dGVzdC1jYS1kYXRhCg==",
					Namespace: &namespace,
					Permissions: kubernetes.ProviderPermissions{
						Read:  []string{"gg_test2"},
						Write: []string{"gg_test2"},
					},
				},
				{
					Name:       "test-name3",
					Host:       "test-host3",
					CAData:     "dGVzdC1jYS1kYXRhCg==",
					Namespaces: []string{namespace},
					Permissions: kubernetes.ProviderPermissions{
						Read:  []string{"gg_test3"},
						Write: []string{"gg_test3"},
					},
				},
			}

			fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns(testProviders, nil)
			uri = svr.URL + "/v1/kubernetes/providers"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting providers returns a generic error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns([]kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesProviderGetGenericError)
			})
		})

		When("no records found", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersAndPermissionsReturns([]kubernetes.Provider{}, nil)
			})

			It("returns empty list", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse("[]")
			})
		})

		When("it succeeds", func() {
			It("returns ok and the provider", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListKubernetesProviders)
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
