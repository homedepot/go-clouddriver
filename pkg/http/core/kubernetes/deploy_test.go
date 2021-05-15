package kubernetes_test

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Deploy", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		Deploy(c, deployManifestRequest)
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
				Expect(fakeKubeClient.ApplyWithNamespaceOverrideCallCount()).To(Equal(2))
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
			u, _ := fakeKubeClient.ApplyWithNamespaceOverrideArgsForCall(0)
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

	When("applying the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{}, errors.New("error applying manifest"))
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
				fakeKubeClient.ApplyWithNamespaceOverrideReturns(kubernetes.Metadata{Kind: kind}, nil)
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
		})
	})
})
