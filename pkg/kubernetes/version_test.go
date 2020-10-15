package kubernetes_test

import (
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	kc                kubernetes.Controller
	fakeResourcesList *unstructured.UnstructuredList
	currentVersion    string
)

var _ = Describe("Pkg/Kubernetes/Version", func() {
	Context("GetCurrentVersion", func() {
		BeforeEach(func() {
			kc = kubernetes.NewController()
		})

		When("Called with empty resources list", func() {
			BeforeEach(func() {
				fakeResourcesList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
				currentVersion = kc.GetCurrentVersion(fakeResourcesList, "test-kind", "test-name")
			})

			It("return 0 as the current version", func() {
				Expect(currentVersion).To(Equal("0"))
			})
		})

		When("The higest version number in the cluster is 4", func() {
			BeforeEach(func() {
				fakeResourcesList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{
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
								"resourceVersion":   "4",
								"creationTimestamp": "2020-02-14T14:12:03Z",
								"labels": map[string]interface{}{
									"label1":                        "test-label1",
									"moniker.spinnaker.io/sequence": "4",
								},
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/cluster": "pod fakeName",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
				}
				currentVersion = kc.GetCurrentVersion(fakeResourcesList, "pod", "fakeName")
			})

			It("return 4 as the current version", func() {
				Expect(currentVersion).To(Equal("4"))
			})
		})
	})
})
