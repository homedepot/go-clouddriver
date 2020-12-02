package kubernetes_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	. "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
)

var _ = Describe("Delete", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		Delete(c, deleteManifestRequest)
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

	When("the mode is label", func() {
		BeforeEach(func() {
			deleteManifestRequest.Mode = "label"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusNotImplemented))
			Expect(c.Errors.Last().Error()).To(Equal("requested to delete manifest deployment test-deployment using mode label which is not implemented"))
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
			Expect(*deleteOptions.PropagationPolicy).To(Equal(v1.DeletePropagationForeground))
		})
	})
})
