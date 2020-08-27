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
	)

	JustBeforeEach(func() {
		err = AddSpinnakerLabels(u, application)
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
			u, err = ToUnstructured(m)
			Expect(err).To(BeNil())
			application = "test-application"
		})

		It("adds the labels", func() {
			labels := u.GetLabels()
			Expect(labels[LabelKubernetesSpinnakerApp]).To(Equal(application))
			Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

			d := NewDeployment(u.Object).Object()
			templateLabels := d.Spec.Template.ObjectMeta.Labels
			Expect(templateLabels[LabelKubernetesSpinnakerApp]).To(Equal(application))
			Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
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
			u, err = ToUnstructured(m)
			Expect(err).To(BeNil())
			application = "test-application"
		})

		It("adds the labels", func() {
			labels := u.GetLabels()
			Expect(labels[LabelKubernetesSpinnakerApp]).To(Equal(application))
			Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

			d := NewDeployment(u.Object).Object()
			templateLabels := d.Spec.Template.ObjectMeta.Labels
			Expect(templateLabels[LabelKubernetesSpinnakerApp]).To(Equal(application))
			Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
		})
	})
})
