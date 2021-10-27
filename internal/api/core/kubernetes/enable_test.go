package kubernetes_test

import (
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Enable", func() {
	BeforeEach(func() {
		setup()
		fakeKubeClient.GetReturnsOnCall(0, &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "ReplicaSet",
				"apiVersion": "apps/v1",
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"traffic.spinnaker.io/load-balancers": "[\"service test-service1\", \"service test-service2\"]",
					},
					"name":      "test-rs1",
					"namespace": "test-namespace",
					"uid":       "test-uid",
				},
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"labelKey1": "labelValue1",
								"labelKey2": "labelValue2",
							},
						},
					},
				},
			},
		}, nil)
		fakeKubeClient.GetReturnsOnCall(1, &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "Service",
				"apiVersion": "v1",
				"spec": map[string]interface{}{
					"selector": map[string]interface{}{
						"selectorKey1": "selectorValue1",
						"selectorKey2": "selectorValue2",
					},
				},
			},
		}, nil)
		fakeKubeClient.GetReturnsOnCall(2, &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "Service",
				"apiVersion": "v1",
				"spec": map[string]interface{}{
					"selector": map[string]interface{}{
						"selectorKey3": "selectorValue3",
						"selectorKey4": "selectorValue4",
					},
				},
			},
		}, nil)
		fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{
			Items: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":              "test-pod1",
							"namespace":         "test-namespace",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"annotations": map[string]interface{}{
								"moniker.spinnaker.io/application": "wrong-application",
							},
							"labels": map[string]interface{}{
								"labelKey1": "labelValue1",
								"labelKey2": "labelValue2",
							},
							"ownerReferences": []interface{}{
								map[string]interface{}{
									"name": "test-rs1",
									"kind": "replicaSet",
									"uid":  "test-uid",
								},
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
							"name":              "test-pod2",
							"namespace":         "test-namespace",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"annotations": map[string]interface{}{
								"moniker.spinnaker.io/application": "test-application",
							},
							"labels": map[string]interface{}{
								"labelKey1": "labelValue1",
								"labelKey2": "labelValue2",
							},
							"ownerReferences": []interface{}{
								map[string]interface{}{
									"name": "test-rs2",
									"kind": "replicaSet",
									"uid":  "test-uid1",
								},
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
							"name":              "test-pod3",
							"namespace":         "test-namespace",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"annotations": map[string]interface{}{
								"moniker.spinnaker.io/application": "test-application",
							},
							"labels": map[string]interface{}{
								"labelKey1": "labelValue1",
								"labelKey2": "labelValue2",
							},
							"ownerReferences": []interface{}{
								map[string]interface{}{
									"name": "test-rs1",
									"kind": "replicaSet",
									"uid":  "test-uid",
								},
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
							"name":              "test-pod4",
							"namespace":         "test-namespace",
							"creationTimestamp": "2020-02-13T14:12:03Z",
							"annotations": map[string]interface{}{
								"moniker.spinnaker.io/application": "test-application",
							},
							"labels": map[string]interface{}{
								"labelKey1": "labelValue1",
								"labelKey2": "labelValue2",
							},
							"ownerReferences": []interface{}{
								map[string]interface{}{
									"name": "test-rs1",
									"kind": "replicaSet",
									"uid":  "test-uid",
								},
							},
							"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
						},
					},
				},
			},
		}, nil)
		fakeKubeClient.PatchUsingStrategyReturnsOnCall(0, kubernetes.Metadata{}, nil, nil)
		fakeKubeClient.PatchUsingStrategyReturnsOnCall(1, kubernetes.Metadata{}, nil, nil)
		fakeKubeClient.PatchUsingStrategyReturnsOnCall(2, kubernetes.Metadata{}, nil, nil)
		fakeKubeClient.PatchUsingStrategyReturnsOnCall(3, kubernetes.Metadata{}, nil, nil)
	})

	JustBeforeEach(func() {
		kubernetesController.Enable(c, enableManifestRequest)
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("internal: error getting kubernetes provider test-account: error getting provider"))
		})
	})

	When("the manifest name is malformed", func() {
		BeforeEach(func() {
			enableManifestRequest.ManifestName = "malformed"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("manifest name must be in format '{kind} {name}'"))
		})
	})

	When("the provider is namespace scoped and the manifest kind is not valid", func() {
		BeforeEach(func() {
			ns := "test-ns"
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
				Name:      "test-name",
				Host:      "test-host",
				Namespace: &ns,
			}, nil)
			enableManifestRequest.ManifestName = "clusterRole my-role"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'clusterRole'"))
		})
	})

	When("getting the target manifest returns an error not found", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturnsOnCall(0, nil, k8serrors.NewNotFound(schema.GroupResource{Group: "",
				Resource: "fake resource"}, "fake resource not found"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusNotFound))
			Expect(c.Errors.Last().Error()).To(Equal("resource ReplicaSet test-rs-v001 does not exist"))
		})
	})

	When("getting the target manifest returns a generic error", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturnsOnCall(0, nil, errors.New("generic error"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error getting resource (kind: ReplicaSet, name: test-rs-v001, namespace: test-namespace): generic error"))
		})
	})

	When("the load balancers annotation is incorrectly formatted", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturnsOnCall(0, &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "ReplicaSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"traffic.spinnaker.io/load-balancers": "service test-service1, service test-service2\"]",
						},
						"name":      "test-name",
						"namespace": "test-namespace",
						"uid":       "test-uid",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"labelKey1": "labelValue1",
									"labelKey2": "labelValue2",
								},
							},
						},
					},
				},
			}, nil)
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("error unmarshaling annotation 'traffic.spinnaker.io/load-balancers' for resource (kind: ReplicaSet, name: test-name, namespace: test-namespace) into string slice: invalid character 's' looking for beginning of value"))
		})
	})

	When("listing the pods returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ListResourceWithContextReturns(nil, errors.New("error listing pods"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error listing pods"))
		})
	})

	Context("#getLoadBalancer", func() {
		When("the load balancer is incorrectly formatted", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturnsOnCall(0, &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "[\"bad-formatting\"]",
							},
							"name":      "test-name",
							"namespace": "test-namespace",
							"uid":       "test-uid",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"labels": map[string]interface{}{
										"labelKey1": "labelValue1",
										"labelKey2": "labelValue2",
									},
								},
							},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("failed to attach/detach to/from load balancer 'bad-formatting'. load balancers must be specified in the form '{kind} {name}', e.g. 'service my-service'"))
			})
		})

		When("the load balancer kind is not supported", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturnsOnCall(0, &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "[\"bad-kind my-service1\"]",
							},
							"name":      "test-name",
							"namespace": "test-namespace",
							"uid":       "test-uid",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"labels": map[string]interface{}{
										"labelKey1": "labelValue1",
										"labelKey2": "labelValue2",
									},
								},
							},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("no support for load balancing via bad-kind exists in Spinnaker"))
			})
		})

		When("getting the load balancer returns an error not found", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturnsOnCall(1, nil, k8serrors.NewNotFound(schema.GroupResource{Group: "",
					Resource: "fake resource"}, "fake resource not found"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("load balancer service test-service1 does not exist"))
			})
		})

		When("getting the load balancer returns a generic error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturnsOnCall(1, nil, errors.New("generic error"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error getting service test-service1: generic error"))
			})
		})
	})

	When("there are no selectors on the services", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturnsOnCall(1, &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Service",
					"apiVersion": "v1",
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{},
					},
				},
			}, nil)
			fakeKubeClient.GetReturnsOnCall(2, &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Service",
					"apiVersion": "v1",
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{},
					},
				},
			}, nil)
		})

		It("does not call patch", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			Expect(fakeKubeClient.PatchUsingStrategyCallCount()).To(Equal(0))
		})
	})

	When("patching the target returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.PatchUsingStrategyReturnsOnCall(0, kubernetes.Metadata{}, nil, errors.New("error patching"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error patching"))
		})
	})

	When("patching a pod returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.PatchUsingStrategyReturnsOnCall(1, kubernetes.Metadata{}, nil, errors.New("error patching pod"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error patching pod"))
		})
	})

	When("inserting the kubernetes resource returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error creating resource"))
		})
	})

	It("succeeds", func() {
		Expect(c.Writer.Status()).To(Equal(http.StatusOK))
		Expect(fakeKubeClient.PatchUsingStrategyCallCount()).To(Equal(8))
		kind1, name1, namespace1, patchBody1, patchType1 := fakeKubeClient.PatchUsingStrategyArgsForCall(0)
		kind2, name2, namespace2, patchBody2, patchType2 := fakeKubeClient.PatchUsingStrategyArgsForCall(1)
		kind3, name3, namespace3, patchBody3, patchType3 := fakeKubeClient.PatchUsingStrategyArgsForCall(2)
		kind4, name4, namespace4, patchBody4, patchType4 := fakeKubeClient.PatchUsingStrategyArgsForCall(3)
		kind5, name5, namespace5, patchBody5, patchType5 := fakeKubeClient.PatchUsingStrategyArgsForCall(4)
		kind6, name6, namespace6, patchBody6, patchType6 := fakeKubeClient.PatchUsingStrategyArgsForCall(5)
		kind7, name7, namespace7, patchBody7, patchType7 := fakeKubeClient.PatchUsingStrategyArgsForCall(6)
		kind8, name8, namespace8, patchBody8, patchType8 := fakeKubeClient.PatchUsingStrategyArgsForCall(7)
		Expect(kind1).To(Equal("ReplicaSet"))
		Expect(name1).To(Equal("test-rs1"))
		Expect(namespace1).To(Equal("test-namespace"))
		Expect(string(patchBody1)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/spec/template/metadata/labels/selectorKey1",
			  "value": "selectorValue1"
			},
			{
			  "op": "add",
			  "path": "/spec/template/metadata/labels/selectorKey2",
			  "value": "selectorValue2"
			}
		  ]`))
		Expect(patchType1).To(Equal(types.JSONPatchType))

		Expect(kind2).To(Equal("Pod"))
		Expect(name2).To(Equal("test-pod1"))
		Expect(namespace2).To(Equal("test-namespace"))
		Expect(string(patchBody2)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey1",
			  "value": "selectorValue1"
			},
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey2",
			  "value": "selectorValue2"
			}
		  ]`))
		Expect(patchType2).To(Equal(types.JSONPatchType))

		Expect(kind3).To(Equal("Pod"))
		Expect(name3).To(Equal("test-pod3"))
		Expect(namespace3).To(Equal("test-namespace"))
		Expect(string(patchBody3)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey1",
			  "value": "selectorValue1"
			},
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey2",
			  "value": "selectorValue2"
			}
		  ]`))
		Expect(patchType3).To(Equal(types.JSONPatchType))

		Expect(kind4).To(Equal("Pod"))
		Expect(name4).To(Equal("test-pod4"))
		Expect(namespace4).To(Equal("test-namespace"))
		Expect(string(patchBody4)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey1",
			  "value": "selectorValue1"
			},
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey2",
			  "value": "selectorValue2"
			}
		  ]`))
		Expect(patchType4).To(Equal(types.JSONPatchType))

		Expect(kind5).To(Equal("ReplicaSet"))
		Expect(name5).To(Equal("test-rs1"))
		Expect(namespace5).To(Equal("test-namespace"))
		Expect(string(patchBody5)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/spec/template/metadata/labels/selectorKey3",
			  "value": "selectorValue3"
			},
			{
			  "op": "add",
			  "path": "/spec/template/metadata/labels/selectorKey4",
			  "value": "selectorValue4"
			}
		  ]`))
		Expect(patchType5).To(Equal(types.JSONPatchType))

		Expect(kind6).To(Equal("Pod"))
		Expect(name6).To(Equal("test-pod1"))
		Expect(namespace6).To(Equal("test-namespace"))
		Expect(string(patchBody6)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey3",
			  "value": "selectorValue3"
			},
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey4",
			  "value": "selectorValue4"
			}
		  ]`))
		Expect(patchType6).To(Equal(types.JSONPatchType))

		Expect(kind7).To(Equal("Pod"))
		Expect(name7).To(Equal("test-pod3"))
		Expect(namespace7).To(Equal("test-namespace"))
		Expect(string(patchBody7)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey3",
			  "value": "selectorValue3"
			},
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey4",
			  "value": "selectorValue4"
			}
		  ]`))
		Expect(patchType7).To(Equal(types.JSONPatchType))

		Expect(kind8).To(Equal("Pod"))
		Expect(name8).To(Equal("test-pod4"))
		Expect(namespace8).To(Equal("test-namespace"))
		Expect(string(patchBody8)).To(MatchJSON(`[
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey3",
			  "value": "selectorValue3"
			},
			{
			  "op": "add",
			  "path": "/metadata/labels/selectorKey4",
			  "value": "selectorValue4"
			}
		  ]`))
		Expect(patchType8).To(Equal(types.JSONPatchType))
	})
})
