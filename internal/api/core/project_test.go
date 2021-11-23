package core_test

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/front50"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Project", func() {
	Describe("#ListProjectClusters", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/projects/test-project/clusters"
			createRequest(http.MethodGet)
			fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
			fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
				{
					Name: "test-account-1",
				},
				{
					Name: "test-account-2",
				},
			}, nil)
			fakeFront50Client.ProjectReturns(front50.Response{
				Config: front50.Config{
					PipelineConfigs: nil,
					Applications: []string{
						"test-application-1",
						"test-application-2",
					},
					Clusters: []front50.Cluster{
						{
							Account:      "test-account-1",
							Applications: []string{"test-application-1"},
							Detail:       "test-detail",
							Stack:        "test-stack",
						},
						{
							Account:      "test-account-1",
							Applications: []string{"test-application-1"},
							Detail:       "",
							Stack:        "",
						},
						{
							Account:      "test-account-1",
							Applications: []string{"test-application-1", "test-application-2"},
							Detail:       "*",
							Stack:        "*",
						},
						{
							Account:      "invalid-account",
							Applications: []string{"test-application-1"},
							Detail:       "*",
							Stack:        "*",
						},
					},
				},
			}, nil)
			fakeKubeClient.ListResourceWithContextReturnsOnCall(0, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "DaemonSet",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "test-ds-1",
								"namespace":         "test-namespace-0",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/application": "wrong-application",
									"moniker.spinnaker.io/detail":      "test-detail",
									"moniker.spinnaker.io/stack":       "test-stack",
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
										},
									},
								},
							},
							"status": map[string]interface{}{
								"desiredNumberScheduled": 2,
								"numberReady":            1,
							},
						},
					},
					{
						Object: map[string]interface{}{
							"kind":       "DaemonSet",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":              "test-ds-2",
								"namespace":         "test-namespace-1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/application": "test-application-1",
									"moniker.spinnaker.io/detail":      "test-detail",
									"moniker.spinnaker.io/stack":       "test-stack",
								},
							},
							"spec": map[string]interface{}{
								"replicas": 1,
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []map[string]interface{}{
											{
												"image": "test-image-1",
											},
										},
									},
								},
							},
							"status": map[string]interface{}{
								"desiredNumberScheduled": 4,
								"numberReady":            2,
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListResourceWithContextReturnsOnCall(1, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-rs-1",
								"namespace":         "test-namespace-1",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/application": "test-application-1",
								},
							},
							"spec": map[string]interface{}{
								"replicas": 1,
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []map[string]interface{}{
											{
												"image": "test-image-1",
											},
											{
												"image": "test-image-2",
											},
										},
									},
								},
							},
							"status": map[string]interface{}{
								"replicas":      8,
								"readyReplicas": 4,
							},
						},
					},
				},
			}, nil)
			fakeKubeClient.ListResourceWithContextReturnsOnCall(2, &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "StatefulSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":              "test-sts-1",
								"namespace":         "test-namespace-2",
								"creationTimestamp": "2020-02-13T14:12:03Z",
								"annotations": map[string]interface{}{
									"moniker.spinnaker.io/application": "test-application-1",
									"moniker.spinnaker.io/detail":      "another-detail",
									"moniker.spinnaker.io/stack":       "another-stack",
								},
							},
							"spec": map[string]interface{}{
								"replicas": 1,
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []map[string]interface{}{
											{
												"image": "test-image-3",
											},
										},
									},
								},
							},
							"status": map[string]interface{}{
								"replicas":      16,
								"readyReplicas": 8,
							},
						},
					},
				},
			}, nil)
			log.SetOutput(ioutil.Discard)
		})

		JustBeforeEach(func() {
			doRequest()
		})

		AfterEach(func() {
			teardown()
		})

		When("getting front50 project retuns an error", func() {
			BeforeEach(func() {
				fakeFront50Client.ProjectReturns(front50.Response{}, errors.New("error getting front50 project"))
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting front50 project"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("front50 project has no clusters", func() {
			BeforeEach(func() {
				fakeFront50Client.ProjectReturns(front50.Response{
					Config: front50.Config{
						PipelineConfigs: nil,
						Applications:    []string{"test-application1"},
						Clusters:        []front50.Cluster{},
					},
				}, nil)
			})

			It("retuns an empty list", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`[]`)
			})
		})

		When("none of the project's accounts are in the database", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
					{
						Name: "fake-account",
					},
				}, nil)
			})

			It("retuns an empty list", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`[]`)
			})
		})

		When("listing the providers errors", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{}, errors.New("error listing providers"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error listing providers"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("no resources match project config", func() {
			BeforeEach(func() {
				fakeFront50Client.ProjectReturns(front50.Response{
					Config: front50.Config{
						PipelineConfigs: nil,
						Applications: []string{
							"fake-application",
						},
						Clusters: []front50.Cluster{
							{
								Account: "test-account-1",
								Detail:  "*",
								Stack:   "*",
							},
						},
					},
				}, nil)
			})

			It("succeeds, returning applications with empty project clusters", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListProjectClustersNoMatches)
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListProjectClusters)
			})
		})
	})
})
