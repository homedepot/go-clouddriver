package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
)

var _ = Describe("Annotation", func() {
	var (
		u           *unstructured.Unstructured
		application string
		err         error
		kc          Controller
	)

	Context("#AddSpinnakerAnnotations", func() {
		BeforeEach(func() {
			kc = NewController()
		})

		JustBeforeEach(func() {
			err = kc.AddSpinnakerAnnotations(u, application)
		})

		When("the object is a deployment", func() {
			BeforeEach(func() {
				m := map[string]interface{}{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"namespace": "default",
						"name":      "test-name",
					},
				}
				u, err = kc.ToUnstructured(m)
				Expect(err).To(BeNil())
				application = "test-application"
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/deployment"))
				Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("deployment test-name"))

				d := NewDeployment(u.Object).Object()
				templateAnnotations := d.Spec.Template.ObjectMeta.Annotations
				Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/deployment"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("deployment test-name"))
			})
		})

		When("the object is a replicaset", func() {
			BeforeEach(func() {
				m := map[string]interface{}{
					"kind":       "ReplicaSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"namespace": "default",
						"name":      "test-name",
					},
				}
				u, err = kc.ToUnstructured(m)
				Expect(err).To(BeNil())
				application = "test-application"
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/replicaset"))
				Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("replicaset test-name"))

				d := NewReplicaSet(u.Object).Object()
				templateAnnotations := d.Spec.Template.ObjectMeta.Annotations
				Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/replicaset"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("replicaset test-name"))
			})
		})

		When("the object is a daemonset", func() {
			BeforeEach(func() {
				m := map[string]interface{}{
					"kind":       "DaemonSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"namespace": "default",
						"name":      "test-name",
					},
				}
				u, err = kc.ToUnstructured(m)
				Expect(err).To(BeNil())
				application = "test-application"
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/daemonset"))
				Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("daemonset test-name"))

				d := NewDaemonSet(u.Object).Object()
				templateAnnotations := d.Spec.Template.ObjectMeta.Annotations
				Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/daemonset"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("daemonset test-name"))
			})
		})
	})

	Context("#AddSpinnakerVersionAnnotations", func() {
		When("kind is a replicaset", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "test-kind",
						"metadata": map[string]interface{}{
							"name":              "test-name",
							"namespace":         "test-namespace2",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"labels": map[string]interface{}{
								"label1": "test-label1",
							},
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/versioned": "true",
								"moniker.spinnaker.io/cluster":    "test-kind test-name",
							},
							"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
						},
					},
				}

				fakeVersion = kubernetes.SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}

				kc = kubernetes.NewController()
				err = kc.AddSpinnakerVersionAnnotations(&fakeResource, fakeVersion)

			})
			It("expect error not to have occured", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
