package kubernetes_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Rollback", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		kubernetesController.Rollback(c, undoRolloutManifestRequest)
	})

	When("the application is not set", func() {
		BeforeEach(func() {
			c.Request.Header.Del("X-Spinnaker-Application")
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("no application provided"))
		})
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("internal: error getting kubernetes provider spin-cluster-account: error getting provider"))
		})
	})

	When("getting the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturns(nil, errors.New("error getting manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error getting manifest"))
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

	When("listing replicasets returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ListByGVRReturns(nil, errors.New("error listing replicasets"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error listing replicasets"))
		})
	})

	When("the replicaset cannot be found", func() {
		BeforeEach(func() {
			undoRolloutManifestRequest.ManifestName = "deployment wrong"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusNotFound))
			Expect(c.Errors.Last().Error()).To(Equal("revision not found"))
		})
	})

	When("applying the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ApplyReturns(kubernetes.Metadata{}, errors.New("error applying manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error applying manifest"))
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

	Context("when the mode is static", func() {
		BeforeEach(func() {
			undoRolloutManifestRequest.Mode = "static"
			undoRolloutManifestRequest.NumRevisionsBack = 1
		})

		When("num revisions back is less than 1", func() {
			BeforeEach(func() {
				undoRolloutManifestRequest.NumRevisionsBack = 0
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("number of revisions back was less than 1"))
			})
		})

		When("num revisions back is out of range", func() {
			BeforeEach(func() {
				undoRolloutManifestRequest.NumRevisionsBack = 100
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("number of revisions back was out of range"))
			})
		})

		When("there is an error converting the sequence to an int", func() {
			fakeUnstructured1 := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "asdf",
						},
						"name": "test-name",
					},
				},
			}
			fakeUnstructured2 := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "100",
						},
						"name": "test-name",
					},
				},
			}
			fakeUnstructured3 := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "101",
						},
						"name": "test-name",
					},
				},
			}
			fakeUnstructuredList := &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					fakeUnstructured1,
					fakeUnstructured2,
					fakeUnstructured3,
				},
			}

			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(fakeUnstructuredList, nil)
			})

			It("continues", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})

		When("it succeeds", func() {
			fakeUnstructured1 := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "99",
						},
						"name": "test-name",
					},
				},
			}
			fakeUnstructured2 := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "100",
						},
						"name": "test-name",
					},
				},
			}
			fakeUnstructuredList := &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					fakeUnstructured1,
					fakeUnstructured2,
				},
			}

			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(fakeUnstructuredList, nil)
			})

			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})
	})

	When("the annotation moniker is incorrect", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
			fakeUnstructured := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "100",
						},
						"name":      "test-name",
						"namespace": "provider-namespace",
					},
				},
			}
			fakeUnstructuredList := &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					fakeUnstructured,
				},
			}
			fakeKubeClient.ListByGVRReturns(fakeUnstructuredList, nil)
		})

		It("returns status 404 not found", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusNotFound))
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			u := fakeKubeClient.ApplyArgsForCall(0)
			b, _ := json.Marshal(&u)
			Expect(u.GetNamespace()).To(Equal(""))
			Expect(string(b)).To(Equal("{\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/name\":\"test-deployment\",\"artifact.spinnaker.io/type\":\"kubernetes/deployment\",\"moniker.spinnaker.io/application\":\"test-app\"},\"creationTimestamp\":null},\"spec\":{\"selector\":null,\"strategy\":{},\"template\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"containers\":null}}},\"status\":{}}"))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
			fakeUnstructured := unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "test-kind",
					"apiVersion": "test-api-version",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-app",
							kubernetes.AnnotationSpinnakerArtifactName:       "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType:       "kubernetes/deployment",
							"deployment.kubernetes.io/revision":              "100",
						},
						"name":      "test-name",
						"namespace": "provider-namespace",
					},
				},
			}
			fakeUnstructuredList := &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					fakeUnstructured,
				},
			}
			fakeKubeClient.ListByGVRReturns(fakeUnstructuredList, nil)
		})

		When("the kind is not supported", func() {
			BeforeEach(func() {
				undoRolloutManifestRequest.ManifestName = "namespace fake-namespace"
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'namespace'"))
			})
		})

		When("the kind is supported", func() {
			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				_, _, namespace := fakeKubeClient.GetArgsForCall(0)
				Expect(namespace).To(Equal("provider-namespace"))
				u := fakeKubeClient.ApplyArgsForCall(0)
				b, _ := json.Marshal(&u)
				Expect(string(b)).To(Equal("{\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/name\":\"test-deployment\",\"artifact.spinnaker.io/type\":\"kubernetes/deployment\",\"moniker.spinnaker.io/application\":\"test-app\"},\"creationTimestamp\":null},\"spec\":{\"selector\":null,\"strategy\":{},\"template\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"containers\":null}}},\"status\":{}}"))
			})
		})
	})
})
