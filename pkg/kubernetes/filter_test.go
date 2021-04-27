package kubernetes_test

import (
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	filteredResourcesArray, fakeResourcesArray []unstructured.Unstructured
)

var _ = Describe("ManifestFilter", func() {
	Context("#FilterOnClusterAnnotation", func() {
		BeforeEach(func() {
			fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
		})

		When("called with empty list", func() {
			BeforeEach(func() {
				manifestFilter := kubernetes.NewManifestFilter([]unstructured.Unstructured{})
				filteredResourcesArray = manifestFilter.FilterOnClusterAnnotation("test-cluster")
			})

			It("returns an empty list", func() {
				Expect(filteredResourcesArray).To(Equal([]unstructured.Unstructured{}))
			})
		})

		When("called with a list of 4 items 2 of which are part of pod fakeName cluster", func() {
			BeforeEach(func() {
				fakeResourcesArray = []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "3",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster": "pod fakeName",
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
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "3",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster": "pod fakeName",
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
								"name":              "anotherFakeName",
								"namespace":         "test-namespace2",
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "3",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster": "pod anotherFakeName",
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
								"name":              "anotherFakeName",
								"namespace":         "test-namespace2",
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "3",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster": "pod anotherFakeName",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				}
				manifestFilter := kubernetes.NewManifestFilter(fakeResourcesArray)
				filteredResourcesArray = manifestFilter.FilterOnClusterAnnotation("pod fakeName")
			})

			It("returns a of 2 items", func() {
				Expect(len(filteredResourcesArray)).To(Equal(2))
			})
		})
	})

	Context("#FilterOnLabel", func() {
		BeforeEach(func() {
			fakeUnstructuredList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
		})

		When("called with empty list", func() {
			BeforeEach(func() {
				manifestFilter := kubernetes.NewManifestFilter([]unstructured.Unstructured{})
				filteredResourcesArray = manifestFilter.FilterOnLabel("firstFakeLabel")
			})

			It("returns an empty list", func() {
				Expect(filteredResourcesArray).To(Equal([]unstructured.Unstructured{}))
			})
		})

		When("called with a list of 4 items 2 of which have firstFakeLabel label", func() {
			BeforeEach(func() {
				fakeResourcesArray = []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "fakeName",
								"namespace":         "test-namespace2",
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"firstFakeLabel":                "test-label",
									"moniker.spinnaker.io/sequence": "3",
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
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"firstFakeLabel":                "test-label",
									"moniker.spinnaker.io/sequence": "3",
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
								"name":              "anotherFakeName",
								"namespace":         "test-namespace2",
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "3",
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
								"name":              "anotherFakeName",
								"namespace":         "test-namespace2",
								"resourceVersion":   "3",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "3",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				}
				manifestFilter := kubernetes.NewManifestFilter(fakeResourcesArray)
				filteredResourcesArray = manifestFilter.FilterOnLabel("firstFakeLabel")
			})

			It("returns a list of 2 items", func() {
				Expect(len(filteredResourcesArray)).To(Equal(2))
			})
		})
	})
})
