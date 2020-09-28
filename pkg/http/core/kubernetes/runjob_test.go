package kubernetes_test

import (
	"errors"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("RunJob", func() {
	BeforeEach(func() {
		setup()
		fakeUnstructured := unstructured.Unstructured{
			Object: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
						kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
						"deployment.kubernetes.io/revision":        "100",
					},
					"generateName": "test-",
				},
			},
		}
		fakeKubeController.ToUnstructuredReturns(&fakeUnstructured, nil)
	})

	JustBeforeEach(func() {
		action = actionHandler.NewRunJobAction(actionConfig)
		err = action.Run()
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error getting provider"))
		})
	})

	When("there is an error decoding the CA data for the kubernetes provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{CAData: "{}{}{}{}"}, nil)
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("illegal base64 data at input byte 0"))
		})
	})

	When("getting the gcloud access token returns an error", func() {
		BeforeEach(func() {
			fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error getting token"))
		})
	})

	When("creating the kube client returns an error", func() {
		BeforeEach(func() {
			fakeKubeController.NewClientReturns(nil, errors.New("bad config"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("bad config"))
		})
	})

	When("getting the unstructured manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeController.ToUnstructuredReturns(nil, errors.New("error converting to unstructured"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error converting to unstructured"))
		})
	})

	When("adding the spinnaker annotations returns an error", func() {
		BeforeEach(func() {
			fakeKubeController.AddSpinnakerAnnotationsReturns(errors.New("error adding annotations"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error adding annotations"))
		})
	})

	When("adding the spinnaker labels returns an error", func() {
		BeforeEach(func() {
			fakeKubeController.AddSpinnakerLabelsReturns(errors.New("error adding labels"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error adding labels"))
		})
	})

	When("applying the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ApplyReturns(kubernetes.Metadata{}, errors.New("error applying manifest"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error applying manifest"))
		})
	})

	When("creating the resource returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error creating resource"))
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(err).To(BeNil())
		})

		It("generates the name correctly", func() {
			Expect(err).To(BeNil())
			u := fakeKubeClient.ApplyArgsForCall(0)
			name := u.GetName()
			Expect(name).To(HavePrefix("test-"))
			Expect(name).To(HaveLen(10))
		})
	})
})
