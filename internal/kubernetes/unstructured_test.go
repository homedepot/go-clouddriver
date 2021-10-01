package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/resource"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("Unstructured", func() {
	var (
		m   map[string]interface{}
		err error
	)

	Describe("#ToUnstructured", func() {
		var u unstructured.Unstructured

		BeforeEach(func() {
			m = map[string]interface{}{
				"kind":       "Namespace",
				"apiVersion": "v1",
			}
		})

		JustBeforeEach(func() {
			u, err = ToUnstructured(m)
		})

		When("object kind is missing", func() {
			BeforeEach(func() {
				m = map[string]interface{}{}
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Object 'Kind' is missing in '{}'"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(u).ToNot(BeNil())
			})
		})
	})

	Describe("#SetDefaultNamespaceIfScopedAndNoneSet", func() {
		var (
			u      unstructured.Unstructured
			helper *resource.Helper
		)

		BeforeEach(func() {
			m = map[string]interface{}{
				"kind":       "Pod",
				"apiVersion": "v1",
				"metadata":   map[string]interface{}{},
			}
			helper = &resource.Helper{
				NamespaceScoped: false,
			}
			u, err = ToUnstructured(m)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			SetDefaultNamespaceIfScopedAndNoneSet(&u, helper)
		})

		When("it is scoped", func() {
			BeforeEach(func() {
				helper.NamespaceScoped = true
			})

			It("sets the default namespace", func() {
				n := u.GetNamespace()
				Expect(n).To(Equal("default"))
			})
		})

		When("it is not scoped", func() {
			BeforeEach(func() {
				helper.NamespaceScoped = false
			})

			It("does not set the namespace", func() {
				n := u.GetNamespace()
				Expect(n).To(BeEmpty())
			})
		})
	})
})
