package kubernetes_test

import (
	"errors"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Patch", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		action = actionHandler.NewPatchManifestAction(actionConfig)
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

	When("patching the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.PatchUsingStrategyReturns(kubernetes.Metadata{}, nil, errors.New("error patching manifest"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error patching manifest"))
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

	Context("when it succeeds", func() {
		Context("json patch type", func() {
			BeforeEach(func() {
				actionConfig.Operation.PatchManifest.Options.MergeStrategy = "json"
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("merge patch type", func() {
			BeforeEach(func() {
				actionConfig.Operation.PatchManifest.Options.MergeStrategy = "merge"
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})

		It("succeeds", func() {
			Expect(err).To(BeNil())
		})
	})
})
