package core_test

import (
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Instances", func() {
	Describe("#GetInstance", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/instances/test-account/test-namespace/pod test-pod"
			createRequest(http.MethodGet)
			fakeKubeClient.GetReturns(&unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Pod",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/cluster":     "test cluster",
							"moniker.spinnaker.io/application": "test-application",
						},
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
					"status": map[string]interface{}{
						"phase": "Running",
						"containerStatuses": []map[string]interface{}{
							{
								"name": "test-container-name",
								"state": map[string]interface{}{
									"running": map[string]interface{}{
										"startedAt": "2021-05-08T03:29:42Z",
									},
								},
							},
						},
					},
				}}, nil)
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

		When("getting the instance returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting instance"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting instance"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("the pod is in a down state", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(&unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"moniker.spinnaker.io/cluster":     "test cluster",
								"moniker.spinnaker.io/application": "test-application",
							},
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
						"status": map[string]interface{}{
							"phase": "Terminated",
							"containerStatuses": []map[string]interface{}{
								{
									"name": "test-container-name",
									"state": map[string]interface{}{
										"running": map[string]interface{}{
											"startedAt": "2021-05-08T03:29:42Z",
										},
									},
								},
							},
						},
					}}, nil)
			})

			It("returns the instance", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetDownInstance)
			})
		})

		When("the kind contains a .", func() {
			BeforeEach(func() {
				uri = svr.URL + "/instances/test-account/test-namespace/pod.api.k8s.io test-pod"
				createRequest(http.MethodGet)
			})

			It("splits the kind", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetInstance)
			})
		})

		When("it succeeds", func() {
			It("returns the instance", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetInstance)
			})
		})
	})

	Describe("#GetInstanceConsole", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/instances/test-account/test-namespace/pod test-pod/console?provider=kubernetes"
			createRequest(http.MethodGet)
			fakeKubeClient.GetReturns(&unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Pod",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/cluster":     "test cluster",
							"moniker.spinnaker.io/application": "test-application",
						},
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
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name": "test-container-name",
							},
						},
						"initContainers": []map[string]interface{}{
							{
								"name": "test-init-container-name",
							},
						},
					},
					"status": map[string]interface{}{
						"phase": "Running",
						"containerStatuses": []map[string]interface{}{
							{
								"name": "test-container-name",
								"state": map[string]interface{}{
									"running": map[string]interface{}{
										"startedAt": "2021-05-08T03:29:42Z",
									},
								},
							},
						},
					},
				}}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the provider is not kubernetes", func() {
			BeforeEach(func() {
				uri = svr.URL + "/instances/test-account/test-namespace/pod test-pod/console?provider=not-kubernetes"
				createRequest(http.MethodGet)
			})

			It("returns status not implemented", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotImplemented))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Not Implemented"))
				Expect(ce.Message).To(Equal("provider not-kubernetes console not implemented"))
				Expect(ce.Status).To(Equal(http.StatusNotImplemented))
			})
		})

		When("the kind is not a pod", func() {
			BeforeEach(func() {
				uri = svr.URL + "/instances/test-account/test-namespace/not-a-pod test-pod/console?provider=kubernetes"
				createRequest(http.MethodGet)
			})

			It("returns status not implemented", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotImplemented))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Not Implemented"))
				Expect(ce.Message).To(Equal("kind not-a-pod console not implemented"))
				Expect(ce.Status).To(Equal(http.StatusNotImplemented))
			})
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

		When("getting the instance returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting instance"))
			})

			It("returns status internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting instance"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("getting the logs returns an error", func() {
			BeforeEach(func() {
				fakeKubeClientset.PodLogsReturns("", errors.New("error getting logs"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(`{"output": []}`)
			})
		})

		When("the kind contains a .", func() {
			BeforeEach(func() {
				uri = svr.URL + "/instances/test-account/test-namespace/pod.api.k8s.io test-pod/console?provider=kubernetes"
				createRequest(http.MethodGet)
			})

			It("splits the kind", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetInstanceConsole)
			})
		})

		When("it succeeds", func() {
			It("returns the instance", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadGetInstanceConsole)
			})
		})
	})
})
