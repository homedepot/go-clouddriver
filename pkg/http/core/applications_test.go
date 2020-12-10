package core_test

import (
	// . "github.com/homedepot/go-clouddriver/pkg/http/v0"

	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
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
			fakeSQLClient.ListKubernetesClustersByFieldsReturns([]kubernetes.Resource{
				{
					AccountName:  "test-account1",
					Name:         "test-name1",
					Kind:         "test-kind1",
					SpinnakerApp: "test-spinnaker-app1",
				},
				{
					AccountName:  "test-account2",
					Name:         "test-name2",
					Kind:         "test-kind2",
					SpinnakerApp: "test-spinnaker-app2",
				},
				{
					AccountName:  "test-account3",
					Name:         "test-name3",
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
				fakeSQLClient.ListKubernetesClustersByFieldsReturns(nil, errors.New("error listing resources"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
			}, nil)
			fakeKubeClient.ListResourceReturnsOnCall(0, &unstructured.UnstructuredList{
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
					{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-deployment2",
								"namespace":         "test-namespace2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"labels": map[string]interface{}{
									"label1": "test-label2",
								},
								"uid": "cec15437-4e6a-11ea-9788-4201ac100006",
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListResourceReturnsOnCall(1, &unstructured.UnstructuredList{
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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

		When("getting the gcloud access token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
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
				fakeKubeClient.ListResourceReturnsOnCall(0, nil, errors.New("error listing deployments"))
				fakeKubeClient.ListResourceReturnsOnCall(2, nil, errors.New("error listing deployments"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing replicasets returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceReturnsOnCall(1, nil, errors.New("error listing replicaSets"))
				fakeKubeClient.ListResourceReturnsOnCall(3, nil, errors.New("error listing replicaSets"))
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
				// "account2",
			}, nil)
			fakeKubeClient.ListResourceReturnsOnCall(0, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Ingress",
							"apiVersion": "networking.k8s.io/v1beta1",
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
			fakeKubeClient.ListResourceReturnsOnCall(1, &unstructured.UnstructuredList{
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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

		When("getting the gcloud access token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
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
				fakeKubeClient.ListResourceReturnsOnCall(0, nil, errors.New("error listing ingresses"))
				fakeKubeClient.ListResourceReturnsOnCall(1, nil, errors.New("error listing ingresses"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing services returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceReturnsOnCall(1, nil, errors.New("error listing services"))
				fakeKubeClient.ListResourceReturnsOnCall(3, nil, errors.New("error listing services"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListLoadBalancers)
			})
		})
	})

	Describe("#ListClusters", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/clusters"
			createRequest(http.MethodGet)
			fakeSQLClient.ListKubernetesClustersByApplicationReturns([]kubernetes.Resource{
				{
					AccountName: "test-account1",
					Cluster:     "test-kind1 test-name1",
				},
				{
					AccountName: "test-account2",
					Cluster:     "test-kind2 test-name2",
				},
				{
					AccountName: "test-account2",
					Cluster:     "test-kind3 test-name3",
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

		When("listing clusters returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesClustersByApplicationReturns(nil, errors.New("error listing clusters"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing clusters"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("there is an empty cluster", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesClustersByApplicationReturns([]kubernetes.Resource{
					{
						AccountName: "test-account1",
						Cluster:     "test-kind1 test-name1",
					},
					{
						AccountName: "test-account2",
						Cluster:     "",
					},
					{
						AccountName: "test-account2",
						Cluster:     "test-kind3 test-name3",
					},
				}, nil)
			})

			It("is omitted in the response", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListClusters2)
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
			}, nil)
			fakeKubeClient.ListResourceReturnsOnCall(0, &unstructured.UnstructuredList{
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
			fakeKubeClient.ListResourceReturnsOnCall(1, &unstructured.UnstructuredList{
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
			fakeKubeClient.ListResourceReturnsOnCall(2, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "DaemonSet",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "test-ds1",
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
								"desiredNumberScheduled": 2,
								"currentNumberScheduled": 1,
								"numberReady":            1,
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListResourceReturnsOnCall(3, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "StatefulSet",
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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

		When("getting the gcloud access token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
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
				fakeKubeClient.ListResourceReturnsOnCall(0, nil, errors.New("error listing replicasets"))
				fakeKubeClient.ListResourceReturnsOnCall(2, nil, errors.New("error listing replicasets"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("listing pods returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceReturnsOnCall(1, nil, errors.New("error listing pods"))
				fakeKubeClient.ListResourceReturnsOnCall(3, nil, errors.New("error listing pods"))
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
			fakeKubeClient.ListResourceReturns(&unstructured.UnstructuredList{
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("illegal base64 data at input byte 0"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the gcloud access token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting token"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting resource"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing pods returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceReturns(nil, errors.New("error listing pods"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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

	Describe("#GetJob", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/applications/test-application/jobs/test-account/test-namespace/job test-job1"
			createRequest(http.MethodGet)
			fakeKubeClient.GetReturns(&unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Job",
					"apiVersion": "batch/v1",
					"metadata": map[string]interface{}{
						"name":              "test-job1",
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("illegal base64 data at input byte 0"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the gcloud access token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting token"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
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
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting resource"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetJob)
			})
		})
	})
})
