package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("Label", func() {
	var (
		u           unstructured.Unstructured
		application string
		err         error
		version     SpinnakerVersion
		m           map[string]interface{}
	)

	Context("#AddSpinnakerLabels", func() {
		JustBeforeEach(func() {
			err = AddSpinnakerLabels(&u, application)
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
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("adds the labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal(application))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
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
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("does not add the name label", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})

			Context("template labels already exist", func() {
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
										"key1": "value1",
										"key2": "value2",
									},
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
					Expect(templateLabels["key1"]).To(Equal("value1"))
					Expect(templateLabels["key2"]).To(Equal("value2"))
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-application"))
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
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("adds the labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal(application))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
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
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("does not add the name label", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})

			Context("template labels already exist", func() {
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
										"key1": "value1",
										"key2": "value2",
									},
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
					Expect(templateLabels["key1"]).To(Equal("value1"))
					Expect(templateLabels["key2"]).To(Equal("value2"))
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-application"))
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
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("adds the labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal(application))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
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
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("does not add the name label", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})

			Context("template labels already exist", func() {
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
										"key1": "value1",
										"key2": "value2",
									},
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original labels", func() {
					labels := u.GetLabels()
					Expect(labels[LabelKubernetesName]).To(Equal("test-already-here"))
					Expect(labels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))

					templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
					Expect(templateLabels["key1"]).To(Equal("value1"))
					Expect(templateLabels["key2"]).To(Equal("value2"))
					Expect(templateLabels[LabelKubernetesName]).To(Equal("test-application"))
					Expect(templateLabels[LabelKubernetesManagedBy]).To(Equal("spinnaker"))
				})
			})
		})
	})

	Context("#AddSpinnakerVersionLabels", func() {
		JustBeforeEach(func() {
			err = AddSpinnakerVersionLabels(&u, version)
		})

		When("kind is a deployment", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind": "deployment",
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
				}
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
				version = kubernetes.SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
			})

			It("adds the labels", func() {
				Labels := u.GetLabels()
				Expect(Labels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
				templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				Expect(templateLabels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
			})
		})

		When("kind is a replicaset", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind": "replicaset",
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
				}
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
				version = kubernetes.SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
			})

			It("adds the labels", func() {
				Labels := u.GetLabels()
				Expect(Labels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
				templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				Expect(templateLabels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
			})
		})

		When("kind is a daemonset", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind": "daemonset",
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
				}
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
				version = kubernetes.SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
			})

			It("adds the labels", func() {
				Labels := u.GetLabels()
				Expect(Labels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
				templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				Expect(templateLabels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
			})
		})

		When("kind is a statefulset", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind": "statefulset",
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
				}
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
				version = kubernetes.SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
			})

			It("adds the labels", func() {
				Labels := u.GetLabels()
				Expect(Labels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
				templateLabels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				Expect(templateLabels[LabelSpinnakerMonikerSequence]).To(Equal("2"))
			})
		})
	})
})
