package core_test

import (
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Manifest", func() {
	Describe("#GetManifest", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/manifests/test-account/test-namespace/pod test-pod"
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

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("#GetManifestByTarget", func() {
		var target string

		BeforeEach(func() {
			setup()
			target = "newest"
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			uri = svr.URL + "/manifests/test-account/test-namespace/test-kind/cluster/test-application/deployment test-deployment/dynamic/" + target
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

		Context("target is second_newest", func() {
			BeforeEach(func() {
				target = "second_newest"
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
					validateResponse(payloadManifestCoordinates)
				})
			})
		})

		Context("target is oldest", func() {
			BeforeEach(func() {
				target = "oldest"
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
					Expect(ce.Message).To(Equal("requested target \"Oldest\" for cluster deployment test-deployment, but only one resource was found"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateResponse(payloadManifestCoordinates)
				})
			})
		})

		Context("the target is not supported", func() {
			BeforeEach(func() {
				target = "not_supported"
			})

			When("returns an error", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusNotImplemented))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Not Implemented"))
					Expect(ce.Message).To(Equal("requested target \"not_supported\" for cluster deployment test-deployment is not supported"))
					Expect(ce.Status).To(Equal(http.StatusNotImplemented))
				})
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadManifestCoordinates)
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
