package core_test

import (
	// . "github.com/billiford/go-clouddriver/pkg/http/v0"

	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Application", func() {
	Describe("#ListApplications", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesResourcesByFieldsReturns([]kubernetes.Resource{
				{
					AccountName:  "test-account1",
					ID:           "test-id1",
					TaskID:       "test-task-id1",
					APIGroup:     "test-api-group1",
					Name:         "test-name1",
					Namespace:    "test-namespace1",
					Resource:     "test-resource1",
					Version:      "test-version1",
					Kind:         "test-kind1",
					SpinnakerApp: "test-spinnaker-app1",
				},
				{
					AccountName:  "test-account2",
					ID:           "test-id2",
					TaskID:       "test-task-id2",
					APIGroup:     "test-api-group2",
					Name:         "test-name2",
					Namespace:    "test-namespace2",
					Resource:     "test-resource2",
					Version:      "test-version2",
					Kind:         "test-kind2",
					SpinnakerApp: "test-spinnaker-app2",
				},
				{
					AccountName:  "test-account3",
					ID:           "test-id3",
					TaskID:       "test-task-id3",
					APIGroup:     "test-api-group3",
					Name:         "test-name3",
					Namespace:    "test-namespace3",
					Resource:     "test-resource3",
					Version:      "test-version3",
					Kind:         "test-kind3",
					SpinnakerApp: "test-spinnaker-app2",
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing resources by fields returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByFieldsReturns(nil, errors.New("error listing resources"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing resources"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadApplications)
			})
		})
	})

	Describe("#ListServerGroupManagers", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/serverGroupManagers"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesAccountsBySpinnakerAppReturns([]string{
				"account1",
				"account2",
			}, nil)
			fakeKubeClient.ListReturnsOnCall(0, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-deployment1",
								"namespace":         "test-namespace1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(1, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-rs1",
								"namespace": "test-namespace1",
								"annotations": map[string]interface{}{
									"artifact.spinnaker.io/name":        "test-deployment1",
									"artifact.spinnaker.io/type":        "kubernetes/deployment",
									"deployment.kubernetes.io/revision": "236",
								},
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(2, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-deployment2",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-12T14:11:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"uid": "bec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(3, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-rs2",
								"namespace": "test-namespace2",
								"annotations": map[string]interface{}{
									"artifact.spinnaker.io/name":        "test-deployment2",
									"artifact.spinnaker.io/type":        "kubernetes/deployment",
									"deployment.kubernetes.io/revision": "19",
								},
							},
						},
					},
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing kubernetes accounts by spinnaker app returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesAccountsBySpinnakerAppReturns(nil, errors.New("error listing accounts"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing accounts"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the kubernetes provider for an account errors", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("the ca data is bad", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					CAData: "{}",
				}, nil)
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("creating the kube client returns an error", func() {
			BeforeEach(func() {
				fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing deployments returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturnsOnCall(0, nil, errors.New("error listing deployments"))
				fakeKubeClient.ListReturnsOnCall(1, nil, errors.New("error listing deployments"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing replicasets returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturnsOnCall(1, nil, errors.New("error listing replicaSets"))
				fakeKubeClient.ListReturnsOnCall(3, nil, errors.New("error listing replicaSets"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadServerGroupManagers)
			})
		})
	})

	Describe("#ListLoadBalancers", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/loadBalancers"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesAccountsBySpinnakerAppReturns([]string{
				"account1",
				"account2",
			}, nil)
			fakeKubeClient.ListReturnsOnCall(0, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Ingress",
							"apiVersion": "v1beta1",
							"metadata": map[string]interface{}{
								"name":              "test-ingress1",
								"namespace":         "test-namespace1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(1, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service1",
								"namespace": "test-namespace1",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(2, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Ingress",
							"apiVersion": "v1beta1",
							"metadata": map[string]interface{}{
								"name":              "test-ingress2",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(3, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service1",
								"namespace": "test-namespace1",
							},
						},
					},
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing kubernetes accounts by spinnaker app returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesAccountsBySpinnakerAppReturns(nil, errors.New("error listing accounts"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing accounts"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the kubernetes provider for an account errors", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("the ca data is bad", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					CAData: "{}",
				}, nil)
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("creating the kube client returns an error", func() {
			BeforeEach(func() {
				fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing ingresses returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturnsOnCall(0, nil, errors.New("error listing ingresses"))
				fakeKubeClient.ListReturnsOnCall(1, nil, errors.New("error listing ingresses"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing services returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturnsOnCall(1, nil, errors.New("error listing services"))
				fakeKubeClient.ListReturnsOnCall(3, nil, errors.New("error listing services"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadLoadBalancers)
			})
		})
	})

	Describe("#ListClusters", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/clusters"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesResourcesByFieldsReturns([]kubernetes.Resource{
				{
					AccountName:  "test-account1",
					ID:           "test-id1",
					TaskID:       "test-task-id1",
					APIGroup:     "test-api-group1",
					Name:         "test-name1",
					Namespace:    "test-namespace1",
					Resource:     "test-resource1",
					Version:      "test-version1",
					Kind:         "test-kind1",
					SpinnakerApp: "test-spinnaker-app1",
				},
				{
					AccountName:  "test-account2",
					ID:           "test-id2",
					TaskID:       "test-task-id2",
					APIGroup:     "test-api-group2",
					Name:         "test-name2",
					Namespace:    "test-namespace2",
					Resource:     "test-resource2",
					Version:      "test-version2",
					Kind:         "test-kind2",
					SpinnakerApp: "test-spinnaker-app2",
				},
				{
					AccountName:  "test-account2",
					ID:           "test-id3",
					TaskID:       "test-task-id3",
					APIGroup:     "test-api-group3",
					Name:         "test-name3",
					Namespace:    "test-namespace3",
					Resource:     "test-resource3",
					Version:      "test-version3",
					Kind:         "test-kind3",
					SpinnakerApp: "test-spinnaker-app2",
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing resources by fields returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByFieldsReturns(nil, errors.New("error listing resources"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing resources"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListClusters)
			})
		})
	})

	Describe("#ListServerGroups", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/serverGroups"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesAccountsBySpinnakerAppReturns([]string{
				"account1",
				"account2",
			}, nil)
			fakeKubeClient.ListReturnsOnCall(0, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-rs1",
								"namespace":         "test-namespace1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"artifact.spinnaker.io/name":        "test-deployment1",
									"artifact.spinnaker.io/type":        "kubernetes/deployment",
									"artifact.spinnaker.io/location":    "test-namespace1",
									"moniker.spinnaker.io/application":  "test-deployment1",
									"moniker.spinnaker.io/cluster":      "deployment test-deployment1",
									"deployment.kubernetes.io/revision": "19",
								},
							},
							"spec": map[string]interface{}{
								"replicas": 1,
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []map[string]interface{}{
											{
												"image": "test-image1",
											},
											{
												"image": "test-image2",
											},
										},
									},
								},
							},
							"status": map[string]interface{}{
								"replicas":      1,
								"readyReplicas": 0,
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(1, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "test-pod1",
								"namespace":         "test-namespace1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"ownerReferences": []map[string]interface{}{
									{
										"name": "test-rs1",
									},
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(2, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-rs2",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"artifact.spinnaker.io/name":        "test-deployment2",
									"artifact.spinnaker.io/type":        "kubernetes/deployment",
									"artifact.spinnaker.io/location":    "test-namespace2",
									"moniker.spinnaker.io/application":  "test-deployment2",
									"moniker.spinnaker.io/cluster":      "deployment test-deployment1",
									"deployment.kubernetes.io/revision": "19",
								},
							},
							"spec": map[string]interface{}{
								"replicas": 1,
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []map[string]interface{}{
											{
												"image": "test-image3",
											},
											{
												"image": "test-image4",
											},
										},
									},
								},
							},
							"status": map[string]interface{}{
								"replicas":      1,
								"readyReplicas": 0,
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListReturnsOnCall(3, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "test-pod2",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"ownerReferences": []map[string]interface{}{
									{
										"name": "test-rs2",
									},
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing kubernetes accounts by spinnaker app returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesAccountsBySpinnakerAppReturns(nil, errors.New("error listing accounts"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing accounts"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the kubernetes provider for an account errors", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("the ca data is bad", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					CAData: "{}",
				}, nil)
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("creating the kube client returns an error", func() {
			BeforeEach(func() {
				fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing replicasets returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturnsOnCall(0, nil, errors.New("error listing replicasets"))
				fakeKubeClient.ListReturnsOnCall(1, nil, errors.New("error listing replicasets"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing pods returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturnsOnCall(1, nil, errors.New("error listing pods"))
				fakeKubeClient.ListReturnsOnCall(3, nil, errors.New("error listing pods"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListServerGroups)
			})
		})
	})

	Describe("#GetServerGroup", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/serverGroups/test-account/test-namespace/replicaSet test-rs1"
			createRequest(http.MethodGet)
			fakeKubeClient.ListReturns(&unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "test-pod1",
								"namespace":         "test-namespace1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"ownerReferences": []map[string]interface{}{
									{
										"name": "test-rs1",
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
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label1",
								},
								"ownerReferences": []map[string]interface{}{
									{
										"name": "test-rs1",
									},
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.GetReturns(&unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "ReplicaSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"name":              "test-rs1",
						"namespace":         "test-namespace1",
						"creationTimestamp": "2020-02-13T14:12:03Z",
						"annotations": map[string]interface{}{
							"artifact.spinnaker.io/name":        "test-deployment2",
							"artifact.spinnaker.io/type":        "kubernetes/deployment",
							"artifact.spinnaker.io/location":    "test-namespace2",
							"moniker.spinnaker.io/application":  "test-deployment2",
							"moniker.spinnaker.io/cluster":      "deployment test-deployment1",
							"deployment.kubernetes.io/revision": "19",
						},
					},
					"spec": map[string]interface{}{
						"replicas": 1,
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []map[string]interface{}{
									{
										"image": "test-image3",
									},
									{
										"image": "test-image4",
									},
								},
							},
						},
					},
					"status": map[string]interface{}{
						"replicas":      1,
						"readyReplicas": 0,
					},
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
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

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting provider"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("decoding the ca data returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					CAData: "{}",
				}, nil)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("illegal base64 data at input byte 0"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("creating the kube client returns an error", func() {
			BeforeEach(func() {
				fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("bad config"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the resource returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting resource"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting resource"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing pods returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListReturns(nil, errors.New("error listing pods"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing pods"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetServerGroup)
			})
		})
	})
})
