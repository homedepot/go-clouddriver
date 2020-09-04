package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/billiford/go-clouddriver/pkg/kubernetes"
)

var _ = Describe("Annotation", func() {
	var (
		u           *unstructured.Unstructured
		application string
		err         error
		kc          Controller
	)

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
})
