package kubernetes_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	. "github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

var _ = Describe("Delete", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		kubernetesController.Delete(c, deleteManifestRequest)
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

	When("the mode is static", func() {
		When("getting the gvr returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{}, errors.New("error getting gvr"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error getting gvr"))
			})
		})

		When("deleting the resource returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceReturns(errors.New("error deleting resource"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error deleting resource"))
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

		When("the cascading option is true", func() {
			BeforeEach(func() {
				t := true
				deleteManifestRequest.Options.Cascading = &t
			})

			It("leaves the delete propagation to foreground", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("deployment"))
				Expect(name).To(Equal("test-deployment"))
				Expect(namespace).To(Equal("test-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
				Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
				Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
				Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationForeground))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
			})
		})

		When("orphan dependants is set to false", func() {
			BeforeEach(func() {
				f := false
				deleteManifestRequest.Options.OrphanDependants = &f
				deleteManifestRequest.Options.Cascading = nil
			})

			It("leaves the delete propagation to foreground", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("deployment"))
				Expect(name).To(Equal("test-deployment"))
				Expect(namespace).To(Equal("test-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
				Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
				Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
				Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationForeground))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
			})
		})

		When("orphan dependents is true", func() {
			BeforeEach(func() {
				t := true
				deleteManifestRequest.Options.OrphanDependants = &t
				deleteManifestRequest.Options.Cascading = nil
			})

			It("sets the delete propagation to orphan", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("deployment"))
				Expect(name).To(Equal("test-deployment"))
				Expect(namespace).To(Equal("test-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
				Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
				Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
				Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationOrphan))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
			})
		})

		When("cascading is false", func() {
			BeforeEach(func() {
				f := false
				deleteManifestRequest.Options.OrphanDependants = nil
				deleteManifestRequest.Options.Cascading = &f
			})

			It("sets the delete propagation to orphan", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("deployment"))
				Expect(name).To(Equal("test-deployment"))
				Expect(namespace).To(Equal("test-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
				Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
				Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
				Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationOrphan))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
			})
		})

		When("no propagation policy is set", func() {
			BeforeEach(func() {
				deleteManifestRequest.Options.OrphanDependants = nil
				deleteManifestRequest.Options.Cascading = nil
			})

			It("leaves the delete propagation to foreground", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("deployment"))
				Expect(name).To(Equal("test-deployment"))
				Expect(namespace).To(Equal("test-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
				Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
				Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
				Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationForeground))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
			})
		})
	})

	When("the mode is label", func() {
		BeforeEach(func() {
			deleteManifestRequest.Mode = "label"
		})

		When("getting there are no label selectors", func() {
			BeforeEach(func() {
				deleteManifestRequest.LabelSelectors = DeleteManifestRequestLabelSelectors{}
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("requested to delete manifests by label, but no label selectors were provided"))
			})
		})

		When("there are no kinds selected", func() {
			BeforeEach(func() {
				deleteManifestRequest.Kinds = []string{}
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("requested to delete manifests by label, but no kinds were selected"))
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

		When("listing by gvr returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(nil, errors.New("error listing by gvr"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error listing by gvr"))
			})
		})

		When("deleting the resource returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceReturns(errors.New("error deleting resource"))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
				Expect(c.Errors.Last().Error()).To(Equal("error deleting resource"))
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

		When("there ar no resources to delete", func() {
			BeforeEach(func() {
				fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{}, nil)
			})

			It("gracefully succeeds w/o deleting", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceCallCount()).To(Equal(0))
				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeNoOp))
			})
		})

		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			_, listOptions := fakeKubeClient.ListByGVRArgsForCall(0)
			Expect(listOptions.LabelSelector).To(Equal("key1=key1-value1,key2,key3 notin (key3-value1,key3-value2)"))
			Expect(listOptions.FieldSelector).To(Equal("metadata.namespace=test-namespace"))
			kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
			Expect(kind).To(Equal("deployment"))
			Expect(name).To(Equal("test-name"))
			Expect(namespace).To(Equal("test-namespace"))
			Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
			Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
			Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
			Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationOrphan))
			kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
			Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
		})
	})

	When("the mode is invalid", func() {
		BeforeEach(func() {
			deleteManifestRequest.Mode = "invalid"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusNotImplemented))
			Expect(c.Errors.Last().Error()).To(Equal("requested to delete manifest deployment test-deployment using mode invalid which is not implemented"))
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
			Expect(kind).To(Equal("deployment"))
			Expect(name).To(Equal("test-deployment"))
			Expect(namespace).To(Equal("test-namespace"))
			Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
			Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
			Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
			Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationOrphan))
			kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
			Expect(kr.TaskType).To(Equal(clouddriver.TaskTypeDelete))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
		})

		When("the kind is not supported", func() {
			BeforeEach(func() {
				deleteManifestRequest.ManifestName = "namespace someNamespace"
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'namespace'"))
			})
		})

		When("the kind is supported", func() {
			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				kind, name, namespace, deleteOptions := fakeKubeClient.DeleteResourceByKindAndNameAndNamespaceArgsForCall(0)
				Expect(kind).To(Equal("deployment"))
				Expect(name).To(Equal("test-deployment"))
				Expect(namespace).To(Equal("provider-namespace"))
				Expect(deleteOptions.GracePeriodSeconds).ToNot(BeNil())
				Expect(*deleteOptions.GracePeriodSeconds).To(Equal(int64(10)))
				Expect(deleteOptions.PropagationPolicy).ToNot(BeNil())
				Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationOrphan))
			})
		})
	})

})
