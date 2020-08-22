package core_test

import (
	// . "github.com/billiford/go-clouddriver/pkg/http/v0"

	"bytes"
	"errors"
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes", func() {
	Describe("#CreateKubernetesDeployment", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/kubernetes/ops"
			body.Write([]byte(payloadRequestKubernetesOps))
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
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Bad Request"))
				Expect(ce.Message).To(Equal("invalid character 'd' looking for beginning of value"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
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

		When("setting the kube config returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.WithConfigReturns(errors.New("bad config"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("bad config"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("applying the manifest returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(nil, kubernetes.Metadata{}, errors.New("error applying manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error applying manifest"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("creating the kubernetes resource returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error creating resource"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("#GetManifest", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/manifests/test-account/test-namespace/pod test-pod"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
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

		When("setting the kube config returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.WithConfigReturns(errors.New("bad config"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("bad config"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the manifest returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting manifest"))
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
