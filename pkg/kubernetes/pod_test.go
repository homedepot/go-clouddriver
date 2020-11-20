package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"

	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
)

var _ = Describe("Pod", func() {
	var (
		pod Pod
	)

	BeforeEach(func() {
		p := map[string]interface{}{}
		pod = NewPod(p)
	})

	Describe("#GetObjectMeta", func() {
		BeforeEach(func() {
			_ = pod.GetObjectMeta()
		})

		It("succeeds", func() {
		})
	})

	Describe("#GetPodStatus", func() {
		BeforeEach(func() {
			_ = pod.GetPodStatus()
		})

		It("succeeds", func() {
		})
	})

	Describe("#GetNamespace", func() {
		BeforeEach(func() {
			_ = pod.GetNamespace()
		})

		It("succeeds", func() {
		})
	})

	Describe("#GetName", func() {
		BeforeEach(func() {
			_ = pod.GetName()
		})

		It("succeeds", func() {
		})
	})

	Describe("GetUID", func() {
		BeforeEach(func() {
			_ = pod.GetUID()
		})

		It("succeeds", func() {
		})
	})

	Describe("GetLabels", func() {
		BeforeEach(func() {
			_ = pod.GetLabels()
		})

		It("succeeds", func() {
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		JustBeforeEach(func() {
			s = pod.Status()
		})

		When("pod phase is pending", func() {
			BeforeEach(func() {
				o := pod.Object()
				o.Status.Phase = v1.PodPending
			})

			It("returns expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Pod is Pending"))
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("Pod is Pending"))
			})
		})

		When("pod phase is failed", func() {
			BeforeEach(func() {
				o := pod.Object()
				o.Status.Phase = v1.PodFailed
			})

			It("returns expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Pod is Failed"))
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("Pod is Failed"))
			})
		})

		When("pod phase is unknown", func() {
			BeforeEach(func() {
				o := pod.Object()
				o.Status.Phase = v1.PodUnknown
			})

			It("returns expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Pod is Unknown"))
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("Pod is Unknown"))
			})
		})
	})
})
