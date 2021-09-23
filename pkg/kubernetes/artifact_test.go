package kubernetes_test

import (
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Artifact", func() {
	var (
		err       error
		resource  *unstructured.Unstructured
		artifacts map[string]clouddriver.TaskCreatedArtifact
	)

	BeforeEach(func() {
		artifacts = map[string]clouddriver.TaskCreatedArtifact{
			"gcr.io/test-project/test-container-image": {
				Name:      "gcr.io/test-project/test-container-image",
				Type:      "docker/image",
				Reference: "gcr.io/test-project/test-container-image:v1.0.0",
			},
		}
	})

	JustBeforeEach(func() {
		err = ReplaceDockerImageArtifacts(resource, artifacts)
	})

	Context("#ReplaceDockerImageArtifacts", func() {
		When("no artifacts are passed in", func() {
			BeforeEach(func() {
				artifacts = nil
			})

			It("does not fail", func() {
				Expect(err).To(BeNil())
			})
		})

		When("there is an error getting the nested slice", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{},
						},
					},
				}
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(".spec.containers accessor error: [] is of the type []map[string]interface {}, expected []interface{}"))
			})
		})

		When("the kind does not have containers", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Secret",
						"apiVersion": "v1",
					},
				}
			})

			It("does not fail", func() {
				Expect(err).To(BeNil())
			})
		})

		When("the container is not a map", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"spec": map[string]interface{}{
							"containers": []interface{}{
								"string",
							},
						},
					},
				}
			})

			It("does not fail", func() {
				Expect(err).To(BeNil())
			})
		})

		When("the kind is pod", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "test-container-name",
									"image": "gcr.io/test-project/test-container-image",
								},
							},
						},
					},
				}
			})

			It("updates the container image to contain the reference", func() {
				Expect(err).To(BeNil())
				p := NewPod(resource.Object)
				containers := p.Object().Spec.Containers
				Expect(containers[0].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
			})
		})

		When("the kind is deployment", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
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
					},
				}
			})

			It("updates the container image to contain the reference", func() {
				Expect(err).To(BeNil())
				d := NewDeployment(resource.Object)
				containers := d.Object().Spec.Template.Spec.Containers
				Expect(containers[0].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
			})
		})
	})
})
