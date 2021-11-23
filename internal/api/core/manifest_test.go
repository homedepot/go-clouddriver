package core_test

import (
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Manifest", func() {
	Describe("#GetManifest", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/manifests/test-account/test-namespace/pod test-pod?includeEvents=false"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("internal: error getting kubernetes provider test-account: error getting provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("getting the manifest returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting manifest"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting manifest"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the manifest returns null values", func() {
			BeforeEach(func() {
				uri = svr.URL + "/manifests/test-account/test-namespace/clusterRole test-cluster-role?includeEvents=false"
				createRequest(http.MethodGet)
				fakeKubeClient.GetReturns(&unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "rbac.authorization.k8s.io/v1",
						"kind":       "ClusterRole",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"artifact.spinnaker.io/location":   "",
								"artifact.spinnaker.io/name":       "test-cluster-role",
								"artifact.spinnaker.io/type":       "kubernetes/clusterRole",
								"artifact.spinnaker.io/version":    "",
								"moniker.spinnaker.io/application": "test-application",
								"moniker.spinnaker.io/cluster":     "clusterRole test-cluster-role",
							},
							"creationTimestamp": "2021-10-20T15:29:26Z",
							"labels": map[string]interface{}{
								"app.kubernetes.io/managed-by": "spinnaker",
								"app.kubernetes.io/name":       "test-application",
							},
							"name":            "test-cluster-role",
							"resourceVersion": "53990465",
							"uid":             "d1f1ab80-1320-4e2d-8d12-893c326af416",
						},
						"rules": nil,
					},
				}, nil)
			})

			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadManifestClusterRoleNoRules)
			})
		})

		Context("include events", func() {
			BeforeEach(func() {
				setup()
				uri = svr.URL + "/manifests/test-account/test-namespace/pod test-pod"
				createRequest(http.MethodGet)
				events := []v1.Event{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "test-kind",
							APIVersion: "test-api-version",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:         "test-event-name",
							GenerateName: "test-event-generate-name",
							Namespace:    "test-event-namespace",
						},
						InvolvedObject: v1.ObjectReference{
							Kind:      "test-kind",
							Namespace: "test-namespace",
							Name:      "test-pod",
						},
						Reason:  "test reason",
						Message: "test message",
						Count:   1,
					},
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "test-kind",
							APIVersion: "test-api-version",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:         "test-event-name2",
							GenerateName: "test-event-generate-name",
							Namespace:    "test-event-namespace",
						},
						InvolvedObject: v1.ObjectReference{
							Kind:      "test-kind",
							Namespace: "test-namespace",
							Name:      "test-pod",
						},
						Reason:  "test reason",
						Message: "test message",
						Count:   2,
					},
				}
				fakeKubeClientset.EventsReturns(events, nil)
			})

			When("getting the events returns an error", func() {
				BeforeEach(func() {
					fakeKubeClientset.EventsReturns(nil, errors.New("error getting events"))
				})

				It("fails silently and returns the manifest", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadManifestNoEvents)
				})
			})

			When("getting the events succeeds", func() {
				It("returns the events", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadManifestIncludeEvents)
				})
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadManifestNoEvents)
			})
		})
	})

	Describe("#GetManifestByCriteria", func() {
		var criteria string

		BeforeEach(func() {
			setup()
			fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-application",
								},
								"creationTimestamp": "2021-01-14T12:28:23Z",
								"name":              "rs-wrong-application",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(1),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"creationTimestamp": "2021-03-14T12:28:23Z",
								"name":              "rs-second-newest",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(2),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"creationTimestamp": "2021-02-14T12:28:23Z",
								"name":              "rs-oldest",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(2),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"creationTimestamp": "2021-04-14T12:28:23Z",
								"name":              "rs-newest",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(2),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"creationTimestamp": "2021-02-20T12:28:23Z",
								"name":              "rs-smallest",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(1),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"creationTimestamp": "2021-02-20T12:28:23Z",
								"name":              "rs-largest",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(4),
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"creationTimestamp": "2021-02-20T12:28:23Z",
								"name":              "rs-second-largest",
								"namespace":         "test-namespace",
							},
							"spec": map[string]interface{}{
								"replicas": int64(3),
							},
						},
					},
				},
			}, nil)
			criteria = "newest"
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			uri = svr.URL + "/manifests/test-account/test-namespace/test-kind/cluster/test-application/deployment test-deployment/dynamic/" + criteria
			createRequest(http.MethodGet)
			doRequest()
		})

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("internal: error getting kubernetes provider test-account: error getting provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("getting the gvr returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{}, errors.New("error getting gvr"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting gvr"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing resources returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(nil, errors.New("error listing resources"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing resources"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("there are no resources found", func() {
			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{}, nil)
			})

			It("returns status not found", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotFound))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Not Found"))
				Expect(ce.Message).To(Equal("no resources found for cluster deployment test-deployment"))
				Expect(ce.Status).To(Equal(http.StatusNotFound))
			})
		})

		Context("criteria is second_newest", func() {
			BeforeEach(func() {
				criteria = "second_newest"
			})

			When("there are less than two resources returned", func() {
				BeforeEach(func() {
					fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{
						Items: []unstructured.Unstructured{
							{
								Object: map[string]interface{}{
									"metadata": map[string]interface{}{
										"annotations": map[string]interface{}{
											kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
											kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-application",
										},
									},
								},
							},
							{
								Object: map[string]interface{}{
									"metadata": map[string]interface{}{
										"annotations": map[string]interface{}{
											kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
											kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
										},
									},
								},
							},
						},
					}, nil)
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal("requested target \"Second Newest\" for cluster deployment test-deployment, but only one resource was found"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(`{
            "kind": "test-kind",
            "name": "rs-second-newest",
            "namespace": "test-namespace"
          }`)
				})
			})
		})

		Context("criteria is oldest", func() {
			BeforeEach(func() {
				criteria = "oldest"
			})

			When("there is one resource returned", func() {
				BeforeEach(func() {
					fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{
						Items: []unstructured.Unstructured{
							{
								Object: map[string]interface{}{
									"metadata": map[string]interface{}{
										"annotations": map[string]interface{}{
											kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
											kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-application",
										},
										"name":      "test-name",
										"namespace": "test-namespace",
									},
								},
							},
							{
								Object: map[string]interface{}{
									"metadata": map[string]interface{}{
										"annotations": map[string]interface{}{
											kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
											kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
										},
										"name":      "test-name",
										"namespace": "test-namespace",
									},
								},
							},
						},
					}, nil)
				})

				It("returns this resource", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadManifestCoordinates)
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(`{
            "kind": "test-kind",
            "name": "rs-oldest",
            "namespace": "test-namespace"
          }`)
				})
			})
		})

		Context("criteria is smallest", func() {
			BeforeEach(func() {
				criteria = "smallest"
			})

			When("there is one resource returned", func() {
				BeforeEach(func() {
					fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{
						Items: []unstructured.Unstructured{
							{
								Object: map[string]interface{}{
									"metadata": map[string]interface{}{
										"annotations": map[string]interface{}{
											kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
											kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-application",
										},
										"name":      "test-name",
										"namespace": "test-namespace",
									},
								},
							},
							{
								Object: map[string]interface{}{
									"metadata": map[string]interface{}{
										"annotations": map[string]interface{}{
											kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
											kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
										},
										"name":      "test-name",
										"namespace": "test-namespace",
									},
								},
							},
						},
					}, nil)
				})

				It("returns this resource", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadManifestCoordinates)
				})
			})

			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`{
            "kind": "test-kind",
            "name": "rs-smallest",
            "namespace": "test-namespace"
          }`)
			})
		})

		Context("the criteria is not supported", func() {
			BeforeEach(func() {
				criteria = "not_supported"
			})

			When("returns an error", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal("unknown criteria: not_supported"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})
		})

		When("criteria is newest", func() {
			BeforeEach(func() {
				criteria = "newest"
			})

			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`{
            "kind": "test-kind",
            "name": "rs-newest",
            "namespace": "test-namespace"
          }`)
			})
		})

		When("criteria is largest", func() {
			BeforeEach(func() {
				criteria = "largest"
			})

			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`{
            "kind": "test-kind",
            "name": "rs-largest",
            "namespace": "test-namespace"
          }`)
			})
		})
	})

	Describe("#ListManifestsByCluster", func() {
		BeforeEach(func() {
			setup()
			fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "replicaSet test-cluster",
									kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-application",
								},
								"name":      "rs1-v000",
								"namespace": "test-namespace",
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "replicaSet test-cluster",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"name":      "rs2-v000",
								"namespace": "test-namespace",
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									kubernetes.AnnotationSpinnakerMonikerCluster:     "replicaSet test-cluster",
									kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
								},
								"name":      "rs2-v001",
								"namespace": "test-namespace",
							},
						},
					},
				},
			}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			uri = svr.URL + "/manifests/test-account/test-namespace/test-kind/cluster/test-application/replicaSet test-cluster"
			createRequest(http.MethodGet)
			doRequest()
		})

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("internal: error getting kubernetes provider test-account: error getting provider"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("getting the gvr returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{}, errors.New("error getting gvr"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting gvr"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing resources returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(nil, errors.New("error listing resources"))
			})

			It("returns an empty list", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`[]`)
			})
		})

		When("there are no resources found", func() {
			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{}, nil)
			})

			It("returns an empty list", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`[]`)
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadManifestCoordinatesList)
			})
		})
	})
})
