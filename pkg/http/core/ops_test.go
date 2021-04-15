package core_test

import (
	// . "github.com/homedepot/go-clouddriver/pkg/http/v0"

	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	kube "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes", func() {
	Describe("#CreateKubernetesDeployment", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/kubernetes/ops"
			body.Write([]byte(payloadRequestKubernetesOpsDeployManifest))
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

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("invalid character 'd' looking for beginning of value"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("deploying a manifest returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("delete manifest returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsDeleteManifest))
				createRequest(http.MethodPost)
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("scaling the manifest returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsScaleManifest))
				createRequest(http.MethodPost)
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("cleaning up artifacts returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsCleanupArtifacts))
				createRequest(http.MethodPost)
				fakeKubeController.ToUnstructuredReturns(nil, errors.New("error converting to unstructured"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error converting to unstructured"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("a rolling restart returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsRollingRestartManifest))
				createRequest(http.MethodPost)
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("a run job returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsRunJob))
				createRequest(http.MethodPost)
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("undo rollout returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsUndoRolloutManifest))
				createRequest(http.MethodPost)
				req.Header.Set("X-Spinnaker-Application", "test-app")
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("patch manifest returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsPatchManifest))
				createRequest(http.MethodPost)
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting kubernetes provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting kubernetes provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				or := kube.OperationsResponse{}
				b, _ := ioutil.ReadAll(res.Body)
				json.Unmarshal(b, &or)
				uuidLen := 36
				Expect(or.ID).To(HaveLen(uuidLen))
				Expect(or.ResourceURI).To(HavePrefix("/task"))
			})
		})
	})
})
