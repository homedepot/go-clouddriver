package kubernetes_test

import (
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes/kubernetesfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	kc                              kubernetes.Controller
	fakeUnstructuredList            *unstructured.UnstructuredList
	currentVersion                  string
	isVersioned                     bool
	updatedVersion, expectedVersion kubernetes.SpinnakerVersion
	err                             error
	fakeDeployment                  kubernetes.Deployment
	fakePod                         kubernetes.Pod
)

var _ = Describe("Version", func() {
	Context("#GetCurrentVersion", func() {
		BeforeEach(func() {
			kc = kubernetes.NewController()
			fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
		})

		When("called with empty resources list", func() {
			BeforeEach(func() {
				currentVersion = kc.GetCurrentVersion(fakeUnstructuredList, "test-kind", "test-name")
			})

			It("returns 0 as the current version", func() {
				Expect(currentVersion).To(Equal("-1"))
			})
		})
		When("The highest version number in the cluster is 4", func() {
			BeforeEach(func() {
				fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "0",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster":  "pod fakeName",
									"moniker.spinnaker.io/sequence": "0",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-14T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "4",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster":  "pod fakeName",
									"moniker.spinnaker.io/sequence": "4",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
				}
				currentVersion = kc.GetCurrentVersion(fakeUnstructuredList, "pod", "fakeName")
			})

			It("return 4 as the current version", func() {
				Expect(currentVersion).To(Equal("4"))
			})
		})
		When("#FilterOnClusterAnnotation returns 0 items", func() {
			BeforeEach(func() {
				FakeManifestFilter := kubernetesfakes.FakeManifestFilter{}
				FakeManifestFilter.FilterOnClusterAnnotationReturns([]unstructured.Unstructured{})
				fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{{
					Object: map[string]interface{}{
						"kind": "fakeKind",
						"metadata": map[string]interface{}{
							"name":              "fakeName",
							"namespace":         "test-namespace2",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"labels": map[string]interface{}{
								"label1": "test-label1",
							},
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/versioned": "true",
							},
							"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
						},
					},
				},
				}}
				currentVersion = kc.GetCurrentVersion(fakeUnstructuredList, "test-kind", "test-name")
			})

			It("returns 0 as the current version", func() {
				Expect(currentVersion).To(Equal("-1"))
			})
		})
		When("#FilterOnLabel returns 0 items", func() {
			BeforeEach(func() {
				fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{{
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
				},
				}}
				FakeManifestFilter := kubernetesfakes.FakeManifestFilter{}
				FakeManifestFilter.FilterOnLabelReturns([]unstructured.Unstructured{})
				currentVersion = kc.GetCurrentVersion(fakeUnstructuredList, "test-kind", "test-name")
			})

			It("returns 0 as the current version", func() {
				Expect(currentVersion).To(Equal("-1"))
			})
		})
	})

	Context("#IsVersioned", func() {
		When("#GetAnnotations returns strategy.spinnaker.io/versioned annotaion", func() {
			When("strategy.spinnaker.io/versioned annotaion is true", func() {
				BeforeEach(func() {
					fakeResource := unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "fakeKind",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/versioned": "true",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					}
					isVersioned = kc.IsVersioned(&fakeResource)
				})
				It("returns true", func() {
					Expect(isVersioned).To(Equal(true))
				})
			})
			When("strategy.spinnaker.io/versioned annotaion is false", func() {
				BeforeEach(func() {
					fakeResource := unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "fakeKind",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/versioned": "false",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					}
					isVersioned = kc.IsVersioned(&fakeResource)
				})
				It("returns false", func() {
					Expect(isVersioned).To(Equal(false))
				})
			})
			When("the resource kind is Pod", func() {
				BeforeEach(func() {
					fakeResource := unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "Pod",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					}
					isVersioned = kc.IsVersioned(&fakeResource)
				})
				It("returns true", func() {
					Expect(isVersioned).To(Equal(true))
				})
			})
			When("the resource kind is statefulSet", func() {
				BeforeEach(func() {
					fakeResource := unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "statefulSet",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					}
					isVersioned = kc.IsVersioned(&fakeResource)
				})
				It("returns false", func() {
					Expect(isVersioned).To(Equal(false))
				})
			})
		})
	})

	Context("#IncrementVersion", func() {
		When("current version is 1", func() {
			BeforeEach(func() {
				kc = kubernetes.NewController()
				updatedVersion = kc.IncrementVersion("1")
				expectedVersion = kubernetes.SpinnakerVersion{
					Long:  "v002",
					Short: "2",
				}
			})
			It("returns expected version", func() {
				Expect(updatedVersion).To(Equal(expectedVersion))
			})
		})
		When("current version is 999", func() {
			BeforeEach(func() {
				kc = kubernetes.NewController()
				updatedVersion = kc.IncrementVersion("999")
				expectedVersion = kubernetes.SpinnakerVersion{
					Long:  "v000",
					Short: "0",
				}
			})
			It("returns expected version", func() {
				Expect(updatedVersion).To(Equal(expectedVersion))
			})
		})
	})

	Context("#versionVolumes", func() {
		When("manifest kind is deployment and volume type is configMap", func() {
			BeforeEach(func() {
				fakeResource := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"volumes": []map[string]interface{}{
										{
											"configMap": map[string]interface{}{
												"name": "test-config-map",
											},
											"name": "test-config-map",
										},
									},
								},
							},
						},
					},
				}

				requiredArtifacts := []clouddriver.TaskCreatedArtifact{
					{
						Name:      "test-config-map",
						Type:      "kubernetes/configMap",
						Reference: "test-config-map-v001",
					},
				}

				err = kc.VersionVolumes(fakeResource, requiredArtifacts)
				Expect(err).To(BeNil())
				fakeDeployment = kubernetes.NewDeployment(fakeResource.Object)
			})
			It("updates the volume name to contain the reference", func() {
				volumes := fakeDeployment.GetSpec().Template.Spec.Volumes
				Expect(volumes[0].ConfigMap.Name).To(Equal("test-config-map-v001"))
			})
		})
		When("manifest kind is deployment and volume type is secret", func() {
			BeforeEach(func() {
				fakeResource := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"volumes": []map[string]interface{}{
										{
											"secret": map[string]interface{}{
												"secretName": "test-secret",
											},
											"name": "test-secret",
										},
									},
								},
							},
						},
					},
				}

				requiredArtifacts := []clouddriver.TaskCreatedArtifact{
					{
						Name:      "test-secret",
						Type:      "kubernetes/secret",
						Reference: "test-secret-v001",
					},
				}

				err = kc.VersionVolumes(fakeResource, requiredArtifacts)
				Expect(err).To(BeNil())
				fakeDeployment = kubernetes.NewDeployment(fakeResource.Object)
			})
			It("updates the volume name to contain the reference", func() {
				volumes := fakeDeployment.GetSpec().Template.Spec.Volumes
				Expect(volumes[0].Secret.SecretName).To(Equal("test-secret-v001"))
			})
		})
		When("manifest kind is pod and volume type is configMap", func() {
			BeforeEach(func() {
				fakeResource := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "apps/v1",
						"spec": map[string]interface{}{
							"volumes": []map[string]interface{}{
								{
									"configMap": map[string]interface{}{
										"name": "test-config-map",
									},
									"name": "test-config-map",
								},
							},
						},
					},
				}

				requiredArtifacts := []clouddriver.TaskCreatedArtifact{
					{
						Name:      "test-config-map",
						Type:      "kubernetes/configMap",
						Reference: "test-config-map-v001",
					},
				}

				err = kc.VersionVolumes(fakeResource, requiredArtifacts)
				Expect(err).To(BeNil())
				fakePod = kubernetes.NewPod(fakeResource.Object)
			})
			It("updates the volume name to contain the reference", func() {
				volumes := fakePod.GetSpec().Volumes
				Expect(volumes[0].ConfigMap.Name).To(Equal("test-config-map-v001"))
			})
		})
		When("manifest kind is pod and volume type is secret", func() {
			BeforeEach(func() {
				fakeResource := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "apps/v1",
						"spec": map[string]interface{}{
							"volumes": []map[string]interface{}{
								{
									"Secret": map[string]interface{}{
										"secretName": "test-secret",
									},
									"name": "test-secret",
								},
							},
						},
					},
				}

				requiredArtifacts := []clouddriver.TaskCreatedArtifact{
					{
						Name:      "test-secret",
						Type:      "kubernetes/secret",
						Reference: "test-secret-v001",
					},
				}

				err = kc.VersionVolumes(fakeResource, requiredArtifacts)
				Expect(err).To(BeNil())
				fakePod = kubernetes.NewPod(fakeResource.Object)
			})
			It("updates the volume name to contain the reference", func() {
				volumes := fakePod.GetSpec().Volumes
				Expect(volumes[0].Secret.SecretName).To(Equal("test-secret-v001"))
			})
		})
	})
})
