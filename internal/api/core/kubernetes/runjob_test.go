package kubernetes_test

import (
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunJob", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		kubernetesController.RunJob(c, runJobRequest)
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("internal: error getting kubernetes provider test-account: error getting provider"))
		})
	})

	When("getting the unstructured manifest returns an error", func() {
		BeforeEach(func() {
			runJobRequest.Manifest = map[string]interface{}{}
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("Object 'Kind' is missing in '{}'"))
		})
	})

	When("applying the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{}, errors.New("error applying manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error applying manifest"))
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

	When("the manifest contains a docker image artifact", func() {
		BeforeEach(func() {
			runJobRequest.RequiredArtifacts = []clouddriver.Artifact{
				{
					Reference: "gcr.io/test-project/test-container-image:v1.0.0",
					Name:      "gcr.io/test-project/test-container-image",
					Type:      artifact.TypeDockerImage,
				},
			}
		})

		It("replaces the artifact reference", func() {
			u, _ := fakeKubeClient.ApplyWithNamespaceOverrideArgsForCall(0)
			j := kubernetes.NewJob(u.Object)
			containers := j.Object().Spec.Template.Spec.Containers
			Expect(containers).To(HaveLen(1))
			Expect(containers[0].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
		})

		It("generates the name correctly", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			u, namespace := fakeKubeClient.ApplyWithNamespaceOverrideArgsForCall(0)
			name := u.GetName()
			Expect(string(namespace)).To(Equal("default"))
			Expect(name).To(HavePrefix("test-"))
			Expect(name).To(HaveLen(10))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
		})

		It("succeeds, using providers namespace", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			u, namespace := fakeKubeClient.ApplyWithNamespaceOverrideArgsForCall(0)
			name := u.GetName()
			Expect(string(namespace)).To(Equal("provider-namespace"))
			Expect(name).To(HavePrefix("test-"))
			Expect(name).To(HaveLen(10))
		})
	})

	Context("annotating 'artifact.spinnaker.io/location'", func() {
		When("the namespace is not set", func() {
			BeforeEach(func() {
				runJobRequest.Manifest = map[string]interface{}{
					"kind":       "Job",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"generateName": "test-",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container-name",
										"image": "gcr.io/test-project/test-container-image",
									},
								},
							},
						},
					},
				}
			})

			It("annotates the object accordingly", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u, _ := fakeKubeClient.ApplyWithNamespaceOverrideArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
			})
		})

		When("the namespace is set", func() {
			BeforeEach(func() {
				runJobRequest.Manifest = map[string]interface{}{
					"kind":       "Job",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"generateName": "test-",
						"namespace":    "test-namespace",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container-name",
										"image": "gcr.io/test-project/test-container-image",
									},
								},
							},
						},
					},
				}
			})

			It("annotates the object accordingly", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u, _ := fakeKubeClient.ApplyWithNamespaceOverrideArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(Equal("test-namespace"))
			})
		})
	})
})
