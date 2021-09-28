package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	var (
		provider  Provider
		kind      string
		namespace string = "test-namespace"
	)

	Context("#ValidateKindStatus", func() {
		BeforeEach(func() {
			provider = Provider{}
			kind = "Deployment"
		})

		JustBeforeEach(func() {
			err = provider.ValidateKindStatus(kind)
		})

		When("Provider is namespace-scoped", func() {
			BeforeEach(func() {
				provider.Namespace = &namespace
			})

			When("kind is not allowed", func() {
				BeforeEach(func() {
					kind = "Namespace"
				})

				It("errors", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'Namespace'"))
				})
			})

			When("kind is allowed", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		When("Provider is cluster-scoped", func() {
			When("kind is not allowed", func() {
				BeforeEach(func() {
					kind = "Namespace"
				})

				It("succeeds", func() {
					Expect(err).To(BeNil())
				})
			})

			When("kind is allowed", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
