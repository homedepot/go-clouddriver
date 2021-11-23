package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("Annotation", func() {
	var (
		u           unstructured.Unstructured
		err         error
		application string
		version     SpinnakerVersion
		m           map[string]interface{}
	)

	Context("#AddSpinnakerAnnotations", func() {
		BeforeEach(func() {
			application = "test-application"
		})

		JustBeforeEach(func() {
			u, err = ToUnstructured(m)
			Expect(err).To(BeNil())
			AddSpinnakerAnnotations(&u, application)
		})

		When("the object is a deployment", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"namespace": "default",
						"name":      "test-name",
					},
				}
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/deployment"))
				Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("deployment test-name"))

				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/deployment"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("deployment test-name"))
			})

			Context("template annotations already exist", func() {
				BeforeEach(func() {
					m = map[string]interface{}{
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
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
					Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
					Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/deployment"))
					Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
					Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("deployment test-name"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/deployment"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("deployment test-name"))
				})
			})
		})

		When("the object is a replicaset", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind":       "ReplicaSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"namespace": "default",
						"name":      "test-name",
					},
				}
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/replicaSet"))
				Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("replicaSet test-name"))

				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/replicaSet"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("replicaSet test-name"))
			})

			Context("template annotations already exist", func() {
				BeforeEach(func() {
					m = map[string]interface{}{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
					Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
					Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/replicaSet"))
					Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
					Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("replicaSet test-name"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/replicaSet"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("replicaSet test-name"))
				})
			})
		})

		When("the object is a daemonset", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"kind":       "DaemonSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"namespace": "default",
						"name":      "test-name",
					},
				}
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/daemonSet"))
				Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("daemonSet test-name"))

				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
				Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/daemonSet"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("daemonSet test-name"))
			})

			Context("template annotations already exist", func() {
				BeforeEach(func() {
					m = map[string]interface{}{
						"kind":       "DaemonSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
					Expect(annotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
					Expect(annotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/daemonSet"))
					Expect(annotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
					Expect(annotations[AnnotationSpinnakerMonikerCluster]).To(Equal("daemonSet test-name"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactName]).To(Equal("test-name"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactType]).To(Equal("kubernetes/daemonSet"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerApplication]).To(Equal(application))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerCluster]).To(Equal("daemonSet test-name"))
				})
			})
		})

		Context("annotating 'artifact.spinnaker.io/location'", func() {
			When("the namespace is not set", func() {
				BeforeEach(func() {
					m = map[string]interface{}{
						"kind":       "DaemonSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"name": "test-name",
						},
					}
				})

				It("annotates the location annotation with empty string", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(Equal(""))
				})
			})

			When("the namespace is set", func() {
				BeforeEach(func() {
					m = map[string]interface{}{
						"kind":       "DaemonSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"name":      "test-name",
							"namespace": "test-namespace",
						},
					}
				})

				It("annotates the location annotation with the namespace", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(Equal("test-namespace"))
				})
			})
		})
	})

	Context("#AddSpinnakerVersionAnnotations", func() {
		JustBeforeEach(func() {
			AddSpinnakerVersionAnnotations(&u, version)
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
				version = SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, err = ToUnstructured(m)
				Expect(err).To(BeNil())
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
			})

			Context("template annotations already exist", func() {
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
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				})
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
				version = SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, _ = ToUnstructured(m)
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
			})

			Context("template annotations already exist", func() {
				BeforeEach(func() {
					m := map[string]interface{}{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				})
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
				version = SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, _ = ToUnstructured(m)
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
			})

			Context("template annotations already exist", func() {
				BeforeEach(func() {
					m := map[string]interface{}{
						"kind":       "DaemonSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"namespace": "default",
							"name":      "test-name",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				})
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
				version = SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})

			AfterEach(func() {
				u, _ = ToUnstructured(m)
			})

			It("adds the annotations", func() {
				annotations := u.GetAnnotations()
				Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
				Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
				Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
			})

			Context("template annotations already exist", func() {
				BeforeEach(func() {
					m := map[string]interface{}{
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
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"annotations": map[string]interface{}{
										"annotation1": "value1",
										"annotation2": "value2",
									},
									"namespace": "default",
									"name":      "test-name",
								},
							},
						},
					}
					u, err = ToUnstructured(m)
					Expect(err).To(BeNil())
					application = "test-application"
				})

				It("keeps the original annotations", func() {
					annotations := u.GetAnnotations()
					Expect(annotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(annotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))

					templateAnnotations, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "annotations")
					Expect(templateAnnotations["annotation1"]).To(Equal("value1"))
					Expect(templateAnnotations["annotation2"]).To(Equal("value2"))
					Expect(templateAnnotations[AnnotationSpinnakerArtifactVersion]).To(Equal("v002"))
					Expect(templateAnnotations[AnnotationSpinnakerMonikerSequence]).To(Equal("2"))
				})
			})
		})
	})

	Context("#SpinnakerMonikerApplication", func() {
		var (
			fakeResource unstructured.Unstructured
			value        string
		)

		JustBeforeEach(func() {
			value = SpinnakerMonikerApplication(fakeResource)
		})

		When("annotation is missing", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
					},
				}
			})

			It("returns empty string", func() {
				Expect(value).To(BeEmpty())
			})
		})

		When("annotation is set", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"moniker.spinnaker.io/application": "test-application",
							},
						},
					},
				}
			})

			It("returns annotation value", func() {
				Expect(value).To(Equal("test-application"))
			})
		})
	})

	Context("#AnnotationMatches", func() {
		var (
			fakeResource unstructured.Unstructured
			key          string
			valueToMatch string
			matches      bool
		)

		JustBeforeEach(func() {
			key = AnnotationSpinnakerMonikerStack
			matches = AnnotationMatches(fakeResource, key, valueToMatch)
		})

		Context("Value to match is specified", func() {
			BeforeEach(func() {
				valueToMatch = "test-value"
			})

			When("annotation value does not match", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/stack": "wrong-value",
								},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(matches).To(BeFalse())
				})
			})

			When("annotation value matches", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/stack": "test-value",
								},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(matches).To(BeTrue())
				})
			})
		})

		Context("Value to match is blank", func() {
			BeforeEach(func() {
				valueToMatch = ""
			})

			When("annotation value is specfied", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/stack": "test-value",
								},
							},
						},
					}
				})

				It("returns false", func() {
					Expect(matches).To(BeFalse())
				})
			})

			When("annotation value is blank", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations":                map[string]interface{}{},
								"moniker.spinnaker.io/stack": "",
							},
						},
					}
				})

				It("returns true", func() {
					Expect(matches).To(BeTrue())
				})
			})

			When("annotation is missing", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(matches).To(BeTrue())
				})
			})
		})

		Context("Value to match is the wildcard", func() {
			BeforeEach(func() {
				valueToMatch = "*"
			})

			When("annotation value is specfied", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/stack": "test-value",
								},
							},
						},
					}
				})

				It("returns false", func() {
					Expect(matches).To(BeTrue())
				})
			})

			When("annotation is not specified", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(matches).To(BeTrue())
				})
			})
		})
	})
})
