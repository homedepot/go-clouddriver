package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
)

var _ = Describe("Label", func() {
	var (
		u           *unstructured.Unstructured
		application string
		err         error
		kc          Controller
	)

	Context("#AddSpinnakerLabels", func() {
		BeforeEach(func() {
			kc = NewController()
		})

		JustBeforeEach(func() {
			err = kc.AddSpinnakerLabels(u, application)
		})

		When("the object is a deployment", func() {
			Context("the name label does not exist", func() {
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

				It("adds the labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal(application))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					d := NewDeployment(u.Object).Object()
					templateLabels := d.Spec.Template.ObjectMeta.Labels
					Expect(templateLabels[LabelKubernetesName]).To(Equal(application))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})

			Context("the name label exists", func() {
				BeforeEach(func() {
					m := map[string]interface{}{
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
							"labels": map[string]interface{}{
								"app.kubernetes.io/name": "test-already-here",
							},
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"namespace": "default",
									"name":      "test-name",
									"labels": map[string]interface{}{
										"app.kubernetes.io/name": "test-already-here",
									},
								},
							},
						},
					}
					u, err = kc.ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("does not add the name label", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					d := NewDeployment(u.Object).Object()
					templateLabels := d.Spec.Template.ObjectMeta.Labels
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})
		})

		When("the object is a replicaset", func() {
			Context("the name label does not exist", func() {
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

				It("adds the labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal(application))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					d := NewReplicaSet(u.Object).Object()
					templateLabels := d.Spec.Template.ObjectMeta.Labels
					Expect(templateLabels[LabelKubernetesName]).To(Equal(application))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})

			Context("the name label exists", func() {
				BeforeEach(func() {
					m := map[string]interface{}{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
							"labels": map[string]interface{}{
								"app.kubernetes.io/name": "test-already-here",
							},
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"namespace": "default",
									"name":      "test-name",
									"labels": map[string]interface{}{
										"app.kubernetes.io/name": "test-already-here",
									},
								},
							},
						},
					}
					u, err = kc.ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("does not add the name label", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					d := NewReplicaSet(u.Object).Object()
					templateLabels := d.Spec.Template.ObjectMeta.Labels
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})
		})

		When("the object is a daemonset", func() {
			Context("the name label does not exist", func() {
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

				It("adds the labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal(application))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					d := NewDaemonSet(u.Object).Object()
					templateLabels := d.Spec.Template.ObjectMeta.Labels
					Expect(templateLabels[LabelKubernetesName]).To(Equal(application))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})

			Context("the name label exists", func() {
				BeforeEach(func() {
					m := map[string]interface{}{
						"kind":       "DaemonSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
							"labels": map[string]interface{}{
								"app.kubernetes.io/name": "test-already-here",
							},
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"namespace": "default",
									"name":      "test-name",
									"labels": map[string]interface{}{
										"app.kubernetes.io/name": "test-already-here",
									},
								},
							},
						},
					}
					u, err = kc.ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("does not add the name label", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					d := NewDaemonSet(u.Object).Object()
					templateLabels := d.Spec.Template.ObjectMeta.Labels
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})
		})
	})

	Context("#AddSpinnakerVersionLabels", func() {
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
				err = kc.AddSpinnakerVersionLabels(&fakeResource, fakeVersion)

			})
			It("expect error not to have occured", func() {
				Expect(err).To(BeNil())
			})
			// It("#Label was called with LabelSpinnakerMonikerSequence and 2 version number", func() {

			// })
		})
	})
})
