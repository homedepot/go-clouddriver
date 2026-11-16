package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime/schema"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("selectKindForGroup", func() {
	var (
		candidates = []schema.GroupVersionKind{
			{Group: "gateway.networking.k8s.io", Version: "v1", Kind: "Gateway"},
			{Group: "networking.istio.io", Version: "v1beta1", Kind: "Gateway"},
			{Group: "", Version: "v1", Kind: "Pod"},
		}
	)

	It("picks the candidate in the requested group, not the first one", func() {
		gvk, ok := SelectKindForGroup(candidates, "networking.istio.io")
		Expect(ok).To(BeTrue())
		Expect(gvk.Version).To(Equal("v1beta1"))
	})

	It("matches the group case-insensitively", func() {
		gvk, ok := SelectKindForGroup(candidates, "Gateway.Networking.K8s.IO")
		Expect(ok).To(BeTrue())
		Expect(gvk.Group).To(Equal("gateway.networking.k8s.io"))
	})

	It("matches the core group with the empty string", func() {
		gvk, ok := SelectKindForGroup(candidates, "")
		Expect(ok).To(BeTrue())
		Expect(gvk.Kind).To(Equal("Pod"))
	})

	It("reports no match for a group that is not a candidate", func() {
		_, ok := SelectKindForGroup(candidates, "getambassador.io")
		Expect(ok).To(BeFalse())
	})
})
