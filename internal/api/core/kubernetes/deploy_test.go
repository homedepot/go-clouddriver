package kubernetes_test

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Deploy", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		kubernetesController.Deploy(c, deployManifestRequest)
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("internal: error getting kubernetes provider : error getting provider"))
		})
	})

	When("converting the manifests to unstructured returns an error", func() {
		BeforeEach(func() {
			deployManifestRequest.Manifests = []map[string]interface{}{{}}
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("kubernetes: unable to convert manifest to unstructured: " +
				"Object 'Kind' is missing in '{}'"))
		})
	})

	Context("the kind is a list", func() {
		var fakeUnstructured unstructured.Unstructured
		var manifests []map[string]interface{}

		BeforeEach(func() {
			manifests = []map[string]interface{}{
				{
					"kind":       "ServiceMonitor",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "100",
						},
						"name": "test-list-name",
					},
				},
				{
					"kind":       "ServiceMonitor",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "100",
						},
						"name": "test-list-name2",
					},
				},
			}
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				items := make([]interface{}, 0, 2)
				for _, i := range manifests {
					items = append(items, i)
				}
				fakeUnstructured = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "list",
						"apiVersion": "test-api-version",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
								kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
								"deployment.kubernetes.io/revision":        "100",
							},
							"name": "test-name",
						},
						"items": items,
					},
				}
				deployManifestRequest.Manifests = []map[string]interface{}{
					fakeUnstructured.Object,
				}
			})

			It("merges the list items", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ApplyCallCount()).To(Equal(2))
			})
		})
	})

	When("the kind is a job and generateName is set", func() {
		BeforeEach(func() {
			deployManifestRequest = DeployManifestRequest{
				Manifests: []map[string]interface{}{
					{
						"kind":       "Job",
						"apiVersion": "v1/batch",
						"metadata": map[string]interface{}{
							"generateName": "test-",
						},
					},
				},
			}
		})

		It("generates a unique name for the job", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			u := fakeKubeClient.ApplyArgsForCall(0)
			Expect(u.GetKind()).To(Equal("Job"))
			Expect(u.GetName()).ToNot(BeEmpty())
			Expect(u.GetName()).To(HavePrefix("test-"))
			Expect(u.GetName()).To(HaveLen(10))
		})
	})

	When("the manifest is versioned", func() {
		When("listing resources by kind and namespace returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourcesByKindAndNamespaceReturns(nil, errors.New("ListResourcesByKindAndNamespaceReturns fake error"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("ListResourcesByKindAndNamespaceReturns fake error"))
			})
		})

		When("listing resources by kind and namespace returns an empty list", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourcesByKindAndNamespaceReturns(&unstructured.UnstructuredList{}, nil)
			})

			It("calls GetCurrentVersion with an empty list", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})
	})

	When("the manifest contains a docker image artifact", func() {
		BeforeEach(func() {
			deployManifestRequest.RequiredArtifacts = []clouddriver.Artifact{
				{
					Reference: "gcr.io/test-project/test-container-image:v1.0.0",
					Name:      "gcr.io/test-project/test-container-image",
					Type:      artifact.TypeDockerImage,
				},
			}
		})

		It("replaces the artifact reference", func() {
			u := fakeKubeClient.ApplyArgsForCall(0)
			p := kubernetes.NewPod(u.Object)
			containers := p.Object().Spec.Containers
			Expect(containers).To(HaveLen(1))
			Expect(containers[0].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
		})
	})

	When("the manifest uses source capacity", func() {
		BeforeEach(func() {
			deployManifestRequest = DeployManifestRequest{
				Manifests: []map[string]interface{}{
					{
						"kind":       "Deployment",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/use-source-capacity": "true",
							},
						},
						"spec": map[string]interface{}{
							"replicas": 1,
						},
					},
				},
			}
		})

		When("get current resource returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("GetReturns fake error"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("GetReturns fake error"))
			})
		})

		When("current resource is not found", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, k8serrors.NewNotFound(schema.GroupResource{Group: "", Resource: "fake resource"}, "fake resource not found"))
			})

			It("it does not error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				actual, _, _ := unstructured.NestedInt64(u.Object, "spec", "replicas")
				Expect(actual).To(Equal(int64(1)))
			})
		})

		When("current resource has different replicas value", func() {
			BeforeEach(func() {
				currentManifest := unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"spec": map[string]interface{}{
							"replicas": int64(2),
						},
					},
				}

				fakeKubeClient.GetReturns(&currentManifest, nil)
			})

			It("sets replicas", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				actual, _, _ := unstructured.NestedInt64(u.Object, "spec", "replicas")
				Expect(actual).To(Equal(int64(2)))
			})
		})
	})

	When("the manifest uses recreate strategy", func() {
		BeforeEach(func() {
			deployManifestRequest = DeployManifestRequest{
				Manifests: []map[string]interface{}{
					{
						"kind":       "Job",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/recreate": "true",
							},
							"name":      "test-name",
							"namespace": "test-namespace",
						},
					},
				},
			}
		})

		When("get current resource returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("GetReturns fake error"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("GetReturns fake error"))
			})
		})

		When("current resource does not exist", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, k8serrors.NewNotFound(schema.GroupResource{Group: "", Resource: "fake resource"}, "fake resource not found"))
			})

			It("skips delete", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(Equal(0))
			})
		})

		When("current resource exists", func() {
			BeforeEach(func() {
				currentManifest := unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Job",
						"name": "test-name",
					},
				}

				fakeKubeClient.GetReturns(&currentManifest, nil)
			})

			It("deletes resource before deploying", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("Job"))
				Expect(name).To(Equal("test-name"))
				Expect(namespace).To(Equal("test-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).To(BeNil())
				Expect(deleteOptions.PropagationPolicy).To(BeNil())
			})
		})
	})

	When("the manifest uses replace strategy", func() {
		BeforeEach(func() {
			deployManifestRequest = DeployManifestRequest{
				Manifests: []map[string]interface{}{
					{
						"kind":       "Job",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/replace": "true",
							},
							"name":      "test-name",
							"namespace": "test-namespace",
						},
					},
				},
			}
		})

		When("replace returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ReplaceReturns(kubernetes.Metadata{}, errors.New("ReplaceReturns fake error"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error replacing manifest (kind: Job, apiVersion: v1, name: test-name): ReplaceReturns fake error"))
			})
		})

		It("it succeeds, calling replace", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			Expect(fakeKubeClient.ReplaceCallCount()).To(Equal(1))
		})
	})

	Context("when the manifest uses Spinnaker managed traffic", func() {
		BeforeEach(func() {
			deployManifestRequest = DeployManifestRequest{
				Manifests: []map[string]interface{}{
					{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"name":      "test-name",
							"namespace": "test-namespace",
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
				},
				TrafficManagement: TrafficManagement{
					Enabled: true,
					Options: TrafficManagementOptions{
						EnableTraffic: true,
						Namespace:     "test-namespace",
						Services: []string{
							"service test-service",
							"service test-service2",
						},
					},
				},
			}
			fakeService := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Service",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name":      "test-service",
						"namespace": "test-namespace",
					},
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{
							"selectorKey1": "selectorValue1",
							"selectorKey2": "selectorValue2",
						},
					},
				},
			}
			fakeService2 := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Service",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name":      "test-service2",
						"namespace": "test-namespace",
					},
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{
							"selectorKey3": "selectorValue3",
							"selectorKey4": "selectorValue4",
						},
					},
				},
			}
			fakeKubeClient.GetReturnsOnCall(0, &fakeService, nil)
			fakeKubeClient.GetReturnsOnCall(1, &fakeService2, nil)
		})

		When("the load balancer annotation is already set", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"traffic.spinnaker.io/load-balancers": "[\"service test-service\"]",
								},
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						},
					},
					TrafficManagement: TrafficManagement{
						Enabled: true,
						Options: TrafficManagementOptions{
							EnableTraffic: true,
							Namespace:     "test-namespace",
							Services: []string{
								"service test-service",
								"service test-service2",
							},
						},
					},
				}
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("manifest already has traffic.spinnaker.io/load-balancers annotation set to [\"service test-service\"]. Failed attempting to set it to [service test-service, service test-service2]"))
			})
		})

		When("the load balancer is part of the current request's manifests", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-namespace",
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
						{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service",
								"namespace": "test-namespace",
							},
							"spec": map[string]interface{}{
								"selector": map[string]interface{}{
									"selectorKey1": "selectorValue1",
									"selectorKey2": "selectorValue2",
								},
							},
						},
					},
					TrafficManagement: TrafficManagement{
						Enabled: true,
						Options: TrafficManagementOptions{
							EnableTraffic: true,
							Namespace:     "test-namespace",
							Services: []string{
								"service test-service",
							},
						},
					},
				}
			})

			It("succeeds and does not call the cluster to get the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(BeZero())
			})
		})

		When("the load balancer is part of the current request's manifests and the namespace is overridden", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					NamespaceOverride: "test-namespace",
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-1",
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
						{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service",
								"namespace": "test-2",
							},
							"spec": map[string]interface{}{
								"selector": map[string]interface{}{
									"selectorKey1": "selectorValue1",
									"selectorKey2": "selectorValue2",
								},
							},
						},
					},
					TrafficManagement: TrafficManagement{
						Enabled: true,
						Options: TrafficManagementOptions{
							EnableTraffic: true,
							Namespace:     "test-namespace",
							Services: []string{
								"service test-service",
							},
						},
					},
				}
			})

			It("succeeds and does not call the cluster to get the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(BeZero())
			})
		})

		When("the client has requested to not forward requests to pods", func() {
			BeforeEach(func() {
				deployManifestRequest.TrafficManagement = TrafficManagement{
					Enabled: true,
					Options: TrafficManagementOptions{
						EnableTraffic: false,
						Namespace:     "test-namespace",
						Services: []string{
							"service test-service",
						},
					},
				}
			})

			It("annotates the manifest and does not call the cluster to get the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(BeZero())
				u := fakeKubeClient.ApplyArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerTrafficLoadBalancers]).To(Equal(`["service test-service"]`))
			})
		})

		When("the target manifest does not have any annotations", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
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
						{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service",
								"namespace": "test-namespace",
							},
							"spec": map[string]interface{}{
								"selector": map[string]interface{}{
									"selectorKey1": "selectorValue1",
									"selectorKey2": "selectorValue2",
								},
							},
						},
					},
					TrafficManagement: TrafficManagement{
						Enabled: true,
						Options: TrafficManagementOptions{
							EnableTraffic: true,
							Namespace:     "test-namespace",
							Services: []string{
								"service test-service",
							},
						},
					},
				}
			})

			It("sets the annotations", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerTrafficLoadBalancers]).To(Equal(`["service test-service"]`))
			})
		})

		When("it succeeds", func() {
			It("attaches the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				labels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				Expect(labels["labelKey1"]).To(Equal("labelValue1"))
				Expect(labels["labelKey2"]).To(Equal("labelValue2"))
				Expect(labels["selectorKey1"]).To(Equal("selectorValue1"))
				Expect(labels["selectorKey2"]).To(Equal("selectorValue2"))
				Expect(labels["selectorKey3"]).To(Equal("selectorValue3"))
				Expect(labels["selectorKey4"]).To(Equal("selectorValue4"))
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerTrafficLoadBalancers]).To(Equal(`["service test-service", "service test-service2"]`))
				kind, name, namespace := fakeKubeClient.GetArgsForCall(0)
				Expect(kind).To(Equal("service"))
				Expect(name).To(Equal("test-service"))
				Expect(namespace).To(Equal("test-namespace"))
				kind2, name2, namespace2 := fakeKubeClient.GetArgsForCall(1)
				Expect(kind2).To(Equal("service"))
				Expect(name2).To(Equal("test-service2"))
				Expect(namespace2).To(Equal("test-namespace"))
			})
		})
	})

	Context("when the manifest uses load balancer annotations", func() {
		BeforeEach(func() {
			deployManifestRequest = DeployManifestRequest{
				Manifests: []map[string]interface{}{
					{
						"kind":       "ReplicaSet",
						"apiVersion": "apps/v1",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "[\"service test-service\"]",
							},
							"name":      "test-name",
							"namespace": "test-namespace",
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
				},
			}
			fakeService := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Service",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name":      "test-service",
						"namespace": "test-namespace",
					},
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{
							"selectorKey1": "selectorValue1",
							"selectorKey2": "selectorValue2",
						},
					},
				},
			}
			fakeKubeClient.GetReturns(&fakeService, nil)
		})

		When("the load balancer is annotation is incorrectly formatted", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"traffic.spinnaker.io/load-balancers": "[\"test-ingress\", \"service test-service2\"]",
								},
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						},
					},
				}
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("Failed to attach load balancer 'test-ingress'. " +
					"Load balancers must be specified in the form '{kind} {name}', e.g. 'service my-service'."))
			})
		})

		When("the load balancer kind is not supported", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"traffic.spinnaker.io/load-balancers": "[\"ingress test-ingress\", \"service test-service2\"]",
								},
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						},
					},
				}
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("No support for load balancing via ingress exists in Spinnaker."))
			})
		})

		When("the load balancer is part of the current request's manifests", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"traffic.spinnaker.io/load-balancers": "[\"service test-service\"]",
								},
								"name":      "test-name",
								"namespace": "test-namespace",
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
						{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service",
								"namespace": "test-namespace",
							},
							"spec": map[string]interface{}{
								"selector": map[string]interface{}{
									"selectorKey1": "selectorValue1",
									"selectorKey2": "selectorValue2",
								},
							},
						},
					},
				}
			})

			It("succeeds and does not call the cluster to get the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(BeZero())
			})
		})

		When("the load balancer is part of the current request's manifests and the namespace is overridden", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					NamespaceOverride: "test-namespace",
					Manifests: []map[string]interface{}{
						{
							"kind":       "ReplicaSet",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"traffic.spinnaker.io/load-balancers": "[\"service test-service\"]",
								},
								"name":      "test-name",
								"namespace": "test-1",
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
						{
							"kind":       "Service",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-service",
								"namespace": "test-2",
							},
							"spec": map[string]interface{}{
								"selector": map[string]interface{}{
									"selectorKey1": "selectorValue1",
									"selectorKey2": "selectorValue2",
								},
							},
						},
					},
				}
			})

			It("succeeds and does not call the cluster to get the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(BeZero())
			})
		})

		When("getting the load balancer from the cluster returns a not found error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, k8serrors.NewNotFound(schema.GroupResource{Group: "", Resource: "fake resource"}, "fake resource not found"))
			})

			It("errors", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("Load balancer service test-service does not exist"))
			})
		})

		When("getting the load balancer from the cluster returns a generic error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("generic error"))
			})

			It("errors", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("error getting service test-service: generic error"))
			})
		})

		When("the service has no selectors", func() {
			BeforeEach(func() {
				fakeService := unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Service",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":      "test-service",
							"namespace": "test-namespace",
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{},
						},
					},
				}
				fakeKubeClient.GetReturns(&fakeService, nil)
			})

			It("errors", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("Service must have a non-empty selector in order to be attached to a workload"))
			})
		})

		When("the service selector and target labels are not disjoint", func() {
			BeforeEach(func() {
				fakeService := unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Service",
						"apiVersion": "v1",
						"metadata": map[string]interface{}{
							"name":      "test-service",
							"namespace": "test-namespace",
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"selectorKey1": "selectorValue1",
								"labelKey1":    "labelValue1",
							},
						},
					},
				}
				fakeKubeClient.GetReturns(&fakeService, nil)
			})

			It("errors", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("Service selector must have no label keys in common with target workload"))
			})
		})

		When("the kind is a pod", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "Pod",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"traffic.spinnaker.io/load-balancers": "[\"service test-service\"]",
								},
								"name":      "test-name",
								"namespace": "test-namespace",
								"labels": map[string]interface{}{
									"labelKey1": "labelValue1",
									"labelKey2": "labelValue2",
								},
							},
						},
					},
				}
			})

			It("attaches the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				labels := u.GetLabels()
				Expect(labels["labelKey1"]).To(Equal("labelValue1"))
				Expect(labels["labelKey2"]).To(Equal("labelValue2"))
				Expect(labels["selectorKey1"]).To(Equal("selectorValue1"))
				Expect(labels["selectorKey2"]).To(Equal("selectorValue2"))
			})
		})

		When("it succeeds", func() {
			It("attaches the load balancer", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				labels, _, _ := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
				Expect(labels["labelKey1"]).To(Equal("labelValue1"))
				Expect(labels["labelKey2"]).To(Equal("labelValue2"))
				Expect(labels["selectorKey1"]).To(Equal("selectorValue1"))
				Expect(labels["selectorKey2"]).To(Equal("selectorValue2"))
			})
		})
	})

	When("applying the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ApplyReturns(kubernetes.Metadata{}, errors.New("error applying manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error applying manifest (kind: Pod, apiVersion: v1, name: test-name-v000): error applying manifest"))
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

	Context("generating the cluster", func() {
		When("the kind is deployment", func() {
			kind := "deployment"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("sets the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(Equal(fmt.Sprintf("%s %s", kind, "test-name")))
			})
		})

		When("the kind is statefulSet", func() {
			kind := "statefulSet"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("sets the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(Equal(fmt.Sprintf("%s %s", kind, "test-name")))
			})
		})

		When("the kind is replicaSet", func() {
			kind := "replicaSet"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("sets the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(Equal(fmt.Sprintf("%s %s", kind, "test-name")))
			})
		})

		When("the kind is ingress", func() {
			kind := "ingress"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("sets the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(Equal(fmt.Sprintf("%s %s", kind, "test-name")))
			})
		})

		When("the kind is service", func() {
			kind := "service"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("sets the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(Equal(fmt.Sprintf("%s %s", kind, "test-name")))
			})
		})

		When("the kind is daemonSet", func() {
			kind := "daemonSet"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("sets the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(Equal(fmt.Sprintf("%s %s", kind, "test-name")))
			})
		})

		When("the kind is not a cluster type", func() {
			kind := "pod"

			BeforeEach(func() {
				fakeKubeClient.ApplyReturns(kubernetes.Metadata{Kind: kind}, nil)
			})

			It("does not set the cluster", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.Cluster).To(BeEmpty())
			})
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			u := fakeKubeClient.ApplyArgsForCall(0)
			Expect(string(u.GetNamespace())).To(Equal("default"))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
		})

		When("the kind is not supported", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "Namespace",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name": "sommeNamespace",
							},
						},
					},
				}
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'Namespace'"))
			})
		})

		When("the kind is supported", func() {
			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				Expect(string(u.GetNamespace())).To(Equal("provider-namespace"))
			})
		})
	})

	Context("annotating 'artifact.spinnaker.io/location'", func() {
		When("the kind is namespace-scoped and the namespace is not set", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "Pod",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name": "test-name",
							},
						},
					},
				}
			})

			It("annotates the object accordingly", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(Equal("default"))
			})
		})

		When("the kind is namespace scoped and the namespace is set", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "Pod",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						},
					},
				}
			})

			It("annotates the object accordingly", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(Equal("test-namespace"))
			})
		})

		When("the kind is not namespace-scoped", func() {
			BeforeEach(func() {
				deployManifestRequest = DeployManifestRequest{
					Manifests: []map[string]interface{}{
						{
							"kind":       "ClusterRole",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name": "test-name",
							},
						},
					},
				}
			})

			It("leaves the 'artifact.spinnaker.io/location' annotation empty", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				u := fakeKubeClient.ApplyArgsForCall(0)
				annotations := u.GetAnnotations()
				Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).ToNot(BeNil())
				Expect(annotations[kubernetes.AnnotationSpinnakerArtifactLocation]).To(BeEmpty())
			})
		})
	})
})
