package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("Provider", func() {
	var (
		provider  Provider
		kind      string
		namespace = "test-namespace"
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
			clusterScopedKinds := []string{
				"apiService",
				"clusterRole",
				"clusterRoleBinding",
				"customResourceDefinition",
				"mutatingWebhookConfiguration",
				"namespace",
				"persistentVolume",
				"podSecurityPolicy",
				"storageClass",
				"validatingWebhookConfiguration",
			}

			BeforeEach(func() {
				provider.Namespaces = []string{namespace}
			})

			for _, k := range clusterScopedKinds {
				When("kind "+k+" is not allowed", func() {
					BeforeEach(func() {
						kind = k
					})

					It("errors", func() {
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: '" + kind + "'"))
					})
				})
			}

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
