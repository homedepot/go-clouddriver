package v1_test

import (
	"bytes"
	"errors"
	"net/http"

	// . "github.com/billiford/go-clouddriver/pkg/http/v1"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"

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
})
