package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	fakeUnstructuredList            *unstructured.UnstructuredList
	currentVersion                  string
	isVersioned                     bool
	updatedVersion, expectedVersion SpinnakerVersion
	err                             error
	fakeDeployment                  *Deployment
	fakePod                         *Pod
)

var _ = Describe("Version", func() {
	Context("#GetCurrentVersion", func() {
		BeforeEach(func() {
			fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
		})

		When("called with empty resources list", func() {
			BeforeEach(func() {
				currentVersion = GetCurrentVersion(fakeUnstructuredList, "test-kind", "test-name")
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
				currentVersion = GetCurrentVersion(fakeUnstructuredList, "pod", "fakeName")
			})

			It("return 4 as the current version", func() {
				Expect(currentVersion).To(Equal("4"))
			})
		})

		When("#FilterOnClusterAnnotation returns 0 items", func() {
			BeforeEach(func() {
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
				currentVersion = GetCurrentVersion(fakeUnstructuredList, "test-kind", "test-name")
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
				currentVersion = GetCurrentVersion(fakeUnstructuredList, "test-kind", "test-name")
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
					isVersioned = IsVersioned(fakeResource)
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
					isVersioned = IsVersioned(fakeResource)
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
					isVersioned = IsVersioned(fakeResource)
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
					isVersioned = IsVersioned(fakeResource)
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
				updatedVersion = IncrementVersion("1")
				expectedVersion = SpinnakerVersion{
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
				updatedVersion = IncrementVersion("999")
				expectedVersion = SpinnakerVersion{
					Long:  "v000",
					Short: "0",
				}
			})

			It("returns expected version", func() {
				Expect(updatedVersion).To(Equal(expectedVersion))
			})
		})
	})

	Context("#NameWithoutVersion", func() {
		When("name is not versioned", func() {
			It("returns the name", func() {
				Expect(NameWithoutVersion("")).To(Equal(""))
				Expect(NameWithoutVersion("test-name")).To(Equal("test-name"))
				Expect(NameWithoutVersion("test-v001-name")).To(Equal("test-v001-name"))
				// Must be 3 or more digits
				Expect(NameWithoutVersion("test-name-v1")).To(Equal("test-name-v1"))
				Expect(NameWithoutVersion("test-name-v99")).To(Equal("test-name-v99"))
				// Must be lowercase 'v'
				Expect(NameWithoutVersion("test-name-V001")).To(Equal("test-name-V001"))
			})
		})

		When("name is versioned", func() {
			It("returns the name with the version removed", func() {
				Expect(NameWithoutVersion("test-name-v000")).To(Equal("test-name"))
				Expect(NameWithoutVersion("test-name-v999")).To(Equal("test-name"))
				Expect(NameWithoutVersion("test-name-v1000")).To(Equal("test-name"))
			})
		})
	})

})
