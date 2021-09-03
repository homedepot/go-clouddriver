package kubernetes_test

import (
	"errors"
	"net/http"

	. "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
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
		Rollback(c, undoRolloutManifestRequest)
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

	When("setting the dynamic client returns an error", func() {
		BeforeEach(func() {
			fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("bad config"))
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
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "asdf",
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
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "100",
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
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "101",
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
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "99",
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
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "100",
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

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
		})
	})
})
