package kubernetes_test

import (
	"encoding/json"
	"errors"
	"net/http"

	ops "github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Patch", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		kubernetesController.Patch(c, patchManifestRequest)
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("internal: error getting kubernetes provider spin-cluster-account: error getting provider"))
		})
	})

	Context("when the manifest contains a docker image artifact", func() {
		BeforeEach(func() {
			patchManifestRequest.AllArtifacts = []clouddriver.Artifact{
				{
					Reference: "gcr.io/test-project/test-container-image:v1.0.0",
					Name:      "gcr.io/test-project/test-container-image",
					Type:      artifact.TypeDockerImage,
				},
			}
		})

		When("the patch body cannot be unmarshalled into map[string]interface{}", func() {
			BeforeEach(func() {
				patchManifestRequest.PatchBody = json.RawMessage(
					`[
					{
						"op": "replace",
						"path": "/spec/template/spec/containers/0/image",
						"value": "test/docker/redis@1.0.0"
					}
				]`,
				)
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("json: cannot unmarshal array into Go value of type map[string]interface {}"))
			})
		})

		It("replaces the artifact reference", func() {
			_, _, _, b, _ := fakeKubeClient.PatchUsingStrategyArgsForCall(0)
			Expect(string(b)).To(Equal("{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"image\":\"gcr.io/test-project/test-container-image:v1.0.0\",\"name\":\"test-container-name\"}]}}}}"))
		})
	})

	When("patching the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.PatchUsingStrategyReturns(kubernetes.Metadata{}, nil, errors.New("error patching manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error patching manifest"))
		})
	})

	When("creating the resource returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error creating resource"))
		})

	})

	Context("merge strategies", func() {
		Context("strategic patch type", func() {
			Context("validate request body can unmarshal", func() {
				var body string
				BeforeEach(func() {
					body = `{
										"app": "test-app",
										"cluster": "deployment patch-test",
										"criteria": "newest",
										"kind": "deployment",
										"manifestName": "deployment patch-test",
										"source": "text",
										"mode": "dynamic",
										"patchBody": {
											"spec": {
												"template": {
													"spec": {
														"containers": [
															{
																"image": "test/docker-hub/redis:6.0.9",
																"imagePullPolicy": "IfNotPresent",
																"name": "patch-demo-ctr-2"
															}
														]
													}
												}
											}
										},
										"cloudProvider": "kubernetes",
										"options": {
											"mergeStrategy": "strategic",
											"record": true
										},
										"location": "test-location",
										"account": "test-account",
										"requiredArtifacts": []
								}`
				})

				It("unmarshals", func() {
					pb := ops.PatchManifestRequest{}
					err := json.Unmarshal([]byte(body), &pb)
					Expect(err).To(BeNil())
				})
			})

			BeforeEach(func() {
				patchManifestRequest.Options.MergeStrategy = "strategic"
			})

			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})

		Context("json patch type", func() {
			Context("validate request body can unmarshal", func() {
				var body string
				BeforeEach(func() {
					body = `{
										"app": "test-app",
										"cluster": "deployment patch-test",
										"criteria": "newest",
										"kind": "deployment",
										"manifestName": "deployment patch-test",
										"source": "text",
										"mode": "dynamic",
										"patchBody": [
											{
												"op": "replace",
												"path": "/spec/template/spec/containers/0/image",
												"value": "test/docker/redis@1.0.0"
											}
										],
										"cloudProvider": "kubernetes",
										"options": {
											"mergeStrategy": "json",
											"record": true
										},
										"manifests": [
											{
												"op": "replace",
												"path": "/spec/template/spec/containers/0/image",
												"value": "test/docker-hub/redis@latest"
											}
										],
										"location": "test-location",
										"account": "test-account",
										"requiredArtifacts": []
									}`
				})

				It("unmarshals", func() {
					pb := ops.PatchManifestRequest{}
					err := json.Unmarshal([]byte(body), &pb)
					Expect(err).To(BeNil())
				})
			})

			When("it succeeds", func() {
				BeforeEach(func() {
					patchManifestRequest.Options.MergeStrategy = "json"
				})

				It("succeeds", func() {
					Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				})
			})
		})

		Context("merge patch type", func() {
			Context("validate request body can unmarshal", func() {
				var body string
				BeforeEach(func() {
					body = `{
										"app": "test-app",
										"cluster": "deployment patch-test",
										"criteria": "newest",
										"kind": "deployment",
										"manifestName": "deployment patch-test",
										"source": "text",
										"mode": "dynamic",
										"patchBody": {
											"spec": {
												"template": {
													"spec": {
														"containers": [
															{
																"image": "test/docker-hub/redis:6.0.9",
																"imagePullPolicy": "IfNotPresent",
																"name": "patch-demo-ctr-2"
															}
														]
													}
												}
											}
										},
										"cloudProvider": "kubernetes",
										"options": {
											"mergeStrategy": "json",
											"record": true
										},
										"manifests": [
											{
												"op": "replace",
												"path": "/spec/template/spec/containers/0/image",
												"value": "test/docker-hub/redis@latest"
											}
										],
										"location": "test-location",
										"account": "test-account",
										"requiredArtifacts": []
									}`
				})

				It("unmarshals", func() {
					pb := ops.PatchManifestRequest{}
					err := json.Unmarshal([]byte(body), &pb)
					Expect(err).To(BeNil())
				})
			})

			BeforeEach(func() {
				patchManifestRequest.Options.MergeStrategy = "merge"
			})

			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})

		Context("unknown patch type", func() {
			BeforeEach(func() {
				patchManifestRequest.Options.MergeStrategy = "unknown"
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("invalid merge strategy unknown"))
			})
		})

		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			kind, name, namespace, _, strategy := fakeKubeClient.PatchUsingStrategyArgsForCall(0)
			Expect(string(kind)).To(Equal("deployment"))
			Expect(string(name)).To(Equal("test-deployment"))
			Expect(string(namespace)).To(Equal(""))
			Expect(string(strategy)).To(Equal("application/strategic-merge-patch+json"))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
		})

		When("the kind is not supported", func() {
			BeforeEach(func() {
				patchManifestRequest.ManifestName = "namespace someNamespace"
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'namespace'"))
			})
		})

		When("the kind is supported", func() {
			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, _, strategy := fakeKubeClient.PatchUsingStrategyArgsForCall(0)
				Expect(string(kind)).To(Equal("deployment"))
				Expect(string(name)).To(Equal("test-deployment"))
				Expect(string(namespace)).To(Equal("provider-namespace"))
				Expect(string(strategy)).To(Equal("application/strategic-merge-patch+json"))
			})
		})
	})
})
