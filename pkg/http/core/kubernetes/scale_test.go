package kubernetes_test

import (
	"errors"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scale", func() {
	BeforeEach(func() {
		setup()
		action = actionHandler.NewScaleManifestAction(actionConfig)
	})

	JustBeforeEach(func() {
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

	When("setting the dynamic client returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.SetDynamicClientForConfigReturns(errors.New("error setting client"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error setting client"))
		})
	})

	When("getting the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturns(nil, errors.New("error getting manifest"))
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error getting manifest"))
		})
	})

	When("converting the replicas returns an error", func() {
		BeforeEach(func() {
			actionConfig.Operation.ScaleManifest.Replicas = "asdf"
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("strconv.Atoi: parsing \"asdf\": invalid syntax"))
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

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(err).To(BeNil())
		})
	})
})
