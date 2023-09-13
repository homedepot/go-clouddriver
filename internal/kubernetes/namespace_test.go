package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("Namespace", func() {
	var (
		err       error
		m         map[string]interface{}
		manifest  unstructured.Unstructured
		namespace string
	)

	Context("#SetNamspaceOnManifest", func() {
		BeforeEach(func() {
			m = map[string]interface{}{
				"kind":       "DaemonSet",
				"apiVersion": "apps/v1",
				"metadata": map[string]interface{}{
					"name":      "test-name",
					"namespace": "test-namespace",
				},
			}
		})

		JustBeforeEach(func() {
			manifest, err = ToUnstructured(m)
			Expect(err).To(BeNil())
			SetNamespaceOnManifest(&manifest, namespace)
		})

		When("there is no namespace override", func() {
			BeforeEach(func() {
				namespace = ""
			})

			Context("kubernetes kind is namespace scoped", func() {
				When("manifest namespace is not set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name": "test-name",
							},
						}
					})

					It("sets namespace to 'default'", func() {
						Expect(manifest.GetNamespace()).To(Equal("default"))
					})
				})

				When("manifest namespace is set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						}
					})

					It("does not override namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal("test-namespace"))
					})
				})
			})

			Context("kubernetes kind is cluster scoped", func() {
				When("manifest namespace is not set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Namespace",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name": "test-name",
							},
						}
					})

					It("does not override namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal(""))
					})
				})

				When("manifest namespace is set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Namespace",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						}
					})

					It("does not override namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal("test-namespace"))
					})
				})
			})
		})

		When("there is a namespace override", func() {
			BeforeEach(func() {
				namespace = "test-namespace-override"
			})

			Context("kubernetes kind is namespace scoped", func() {
				When("manifest namespace is not set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name": "test-name",
							},
						}
					})

					It("overrides namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal("test-namespace-override"))
					})
				})

				When("manifest namespace is set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						}
					})

					It("overrides namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal("test-namespace-override"))
					})
				})
			})

			Context("kubernetes kind is cluster scoped", func() {
				When("manifest namespace is not set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Namespace",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name": "test-name",
							},
						}
					})

					It("does not override namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal(""))
					})
				})

				When("manifest namespace is set", func() {
					BeforeEach(func() {
						m = map[string]interface{}{
							"kind":       "Namespace",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name":      "test-name",
								"namespace": "test-namespace",
							},
						}
					})

					It("does not override namespace", func() {
						Expect(manifest.GetNamespace()).To(Equal("test-namespace"))
					})
				})
			})
		})
	})
})
