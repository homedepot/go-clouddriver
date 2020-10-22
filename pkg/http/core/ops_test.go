package core_test

import (
	// . "github.com/billiford/go-clouddriver/pkg/http/v0"

	"bytes"
	"errors"
	"net/http"

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

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Bad Request"))
				Expect(ce.Message).To(Equal("invalid character 'd' looking for beginning of value"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("the request contains no operations", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte("[]"))
				createRequest(http.MethodPost)
			})

			It("returns status ok", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("deploying a manifest returns an error", func() {
			BeforeEach(func() {
				fakeAction.RunReturns(errors.New("error deploying manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error deploying manifest"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("delete manifest returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsDeleteManifest))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error deleting manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error deleting manifest"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("scaling the manifest returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsScaleManifest))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error scaling manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error scaling manifest"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("cleaning up artifacts returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsCleanupArtifacts))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error cleaning up artifacts"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error cleaning up artifacts"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("a rolling restart returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsRollingRestartManifest))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error rolling restart"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error rolling restart"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("a run job returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsRunJob))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error running job"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error running job"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("undo rollout returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsUndoRolloutManifest))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error undoing rollout"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error undoing rollout"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("patch manifest returns an error", func() {
			BeforeEach(func() {
				body = &bytes.Buffer{}
				body.Write([]byte(payloadRequestKubernetesOpsPatchManifest))
				createRequest(http.MethodPost)
				fakeAction.RunReturns(errors.New("error patching manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error patching manifest"))
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
