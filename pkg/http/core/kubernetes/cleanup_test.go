package kubernetes_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
)

var _ = Describe("CleanupArtifacts", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		CleanupArtifacts(c, cleanupArtifactsRequest)
	})

	When("getting the unstructured manifest returns an error", func() {
		BeforeEach(func() {
			cleanupArtifactsRequest.Manifests = []map[string]interface{}{{}}
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("Object 'Kind' is missing in '{}'"))
		})
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("error getting provider"))
		})
	})

	When("there is an error decoding the CA data for the kubernetes provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{CAData: "{}{}{}{}"}, nil)
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("illegal base64 data at input byte 0"))
		})
	})

	When("getting the gcloud access token returns an error", func() {
		BeforeEach(func() {
			fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error getting token"))
		})
	})

	When("creating the kube client returns an error", func() {
		BeforeEach(func() {
			fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("bad config"))
		})
	})

	When("getting the gvr returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{}, errors.New("error getting gvr"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error getting gvr"))
		})
	})

	Context("when annotation 'strategy.spinnaker.io/max-version-history' is set", func() {
		BeforeEach(func() {
			cleanupArtifactsRequest.Manifests = []map[string]interface{}{
				{
					"kind": "ReplicaSet",
					"metadata": map[string]interface{}{
						"name":              "test-name-v002",
						"namespace":         "test-namespace",
						"creationTimestamp": "2020-02-13T14:12:03Z",
						"annotations": map[string]interface{}{
							"strategy.spinnaker.io/max-version-history": "2",
							"moniker.spinnaker.io/cluster":              "replicaSet test-name",
						},
					},
				},
			}
			ul := &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"name":              "test-name-v002",
								"namespace":         "test-namespace",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/max-version-history": "2",
									"moniker.spinnaker.io/cluster":              "replicaSet test-name",
								},
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"name":              "test-name-v001",
								"namespace":         "test-namespace",
								"creationTimestamp": "2020-02-13T13:12:03Z",
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/max-version-history": "2",
									"moniker.spinnaker.io/cluster":              "replicaSet test-name",
								},
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"name":              "test-name-v000",
								"namespace":         "test-namespace",
								"creationTimestamp": "2020-02-13T12:12:03Z",
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/max-version-history": "2",
									"moniker.spinnaker.io/cluster":              "replicaSet test-name",
								},
							},
						},
					},
				},
			}
			fakeKubeClient.ListResourcesByKindAndNamespaceReturns(ul, nil)
		})

		When("the annotation 'moniker.spinnaker.io/cluster' is not set", func() {
			BeforeEach(func() {
				cleanupArtifactsRequest.Manifests = []map[string]interface{}{
					{
						"kind": "ReplicaSet",
						"metadata": map[string]interface{}{
							"name":              "test-name-v002",
							"namespace":         "test-namespace",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/max-version-history": "2",
								// "moniker.spinnaker.io/cluster":              "replicaSet test-name",
							},
						},
					},
				}
			})

			It("does not list resources", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourcesByKindAndNamespaceCallCount()).To(BeZero())
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(BeZero())
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(1))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeCleanup))
			})
		})

		When("listing resources returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourcesByKindAndNamespaceReturns(nil, errors.New("error listing resources"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error listing resources to cleanup for max version history (kind: ReplicaSet, name: test-name-v002, namespace: test-namespace): error listing resources"))
			})
		})

		When("max version history is greater than the number of artifacts in the cluster", func() {
			BeforeEach(func() {
				ul := &unstructured.UnstructuredList{
					Items: []unstructured.Unstructured{
						{
							Object: map[string]interface{}{
								"kind": "ReplicaSet",
								"metadata": map[string]interface{}{
									"name":              "test-name-v000",
									"namespace":         "test-namespace",
									"creationTimestamp": "2020-02-13T14:12:03Z",
									"annotations": map[string]interface{}{
										"strategy.spinnaker.io/max-version-history": "2",
										"moniker.spinnaker.io/cluster":              "replicaSet test-name",
									},
								},
							},
						},
					},
				}
				fakeKubeClient.ListResourcesByKindAndNamespaceReturns(ul, nil)
			})

			It("does not delete anything", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(BeZero())
			})
		})

		When("max version history is equal to the number of artifacts in the cluster", func() {
			BeforeEach(func() {
				ul := &unstructured.UnstructuredList{
					Items: []unstructured.Unstructured{
						{
							Object: map[string]interface{}{
								"kind": "ReplicaSet",
								"metadata": map[string]interface{}{
									"name":              "test-name-v001",
									"namespace":         "test-namespace",
									"creationTimestamp": "2020-02-13T14:12:03Z",
									"annotations": map[string]interface{}{
										"strategy.spinnaker.io/max-version-history": "2",
										"moniker.spinnaker.io/cluster":              "replicaSet test-name",
									},
								},
							},
						},
						{
							Object: map[string]interface{}{
								"kind": "ReplicaSet",
								"metadata": map[string]interface{}{
									"name":              "test-name-v000",
									"namespace":         "test-namespace",
									"creationTimestamp": "2020-02-13T13:12:03Z",
									"annotations": map[string]interface{}{
										"strategy.spinnaker.io/max-version-history": "2",
										"moniker.spinnaker.io/cluster":              "replicaSet test-name",
									},
								},
							},
						},
					},
				}
				fakeKubeClient.ListResourcesByKindAndNamespaceReturns(ul, nil)
			})

			It("does not delete anything", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(BeZero())
			})
		})

		When("the cluster does not match", func() {
			BeforeEach(func() {
				ul := &unstructured.UnstructuredList{
					Items: []unstructured.Unstructured{
						{
							Object: map[string]interface{}{
								"kind": "ReplicaSet",
								"metadata": map[string]interface{}{
									"name":              "test-name-v002",
									"namespace":         "test-namespace",
									"creationTimestamp": "2020-02-13T14:12:03Z",
									"annotations": map[string]interface{}{
										"strategy.spinnaker.io/max-version-history": "2",
										"moniker.spinnaker.io/cluster":              "replicaSet test-name",
									},
								},
							},
						},
						{
							Object: map[string]interface{}{
								"kind": "ReplicaSet",
								"metadata": map[string]interface{}{
									"name":              "test-name-v001",
									"namespace":         "test-namespace",
									"creationTimestamp": "2020-02-13T13:12:03Z",
									"annotations": map[string]interface{}{
										"strategy.spinnaker.io/max-version-history": "2",
										"moniker.spinnaker.io/cluster":              "replicaSet test-name",
									},
								},
							},
						},
						{
							Object: map[string]interface{}{
								"kind": "ReplicaSet",
								"metadata": map[string]interface{}{
									"name":              "test-name-v000",
									"namespace":         "test-namespace",
									"creationTimestamp": "2020-02-13T12:12:03Z",
									"annotations": map[string]interface{}{
										"strategy.spinnaker.io/max-version-history": "2",
										"moniker.spinnaker.io/cluster":              "replicaSet wrong-cluster",
									},
								},
							},
						},
					},
				}
				fakeKubeClient.ListResourcesByKindAndNamespaceReturns(ul, nil)
			})

			It("does not delete anything", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(BeZero())
			})
		})

		When("deleting a resource returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceReturns(errors.New("error deleting resource"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error deleting resource to cleanup for max version history (kind: ReplicaSet, name: test-name-v000, namespace: test-namespace): error deleting resource"))
			})
		})

		When("it deletes the resources", func() {
			It("deletes the oldest resource", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(Equal(1))
				kind, name, namespace, _ := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("ReplicaSet"))
				Expect(name).To(Equal("test-name-v000"))
				Expect(namespace).To(Equal("test-namespace"))
			})
		})
	})

	When("creating the resource returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error creating resource"))
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
			Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeCleanup))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
		})

		It("succeeds,using provider's namespace", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
			Expect(string(kr.Namespace)).To(Equal("provider-namespace"))
		})
	})
})
