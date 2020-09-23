package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
)

var _ = Describe("HorizontalPodAutoscaler", func() {
	var (
		hpa HorizontalPodAutoscaler
	)

	BeforeEach(func(){
		hpa = NewHorizontalPodAutoscaler(map[string]interface{}{})
	})

	Describe("#Status", func() {
		var s manifest.Status

		JustBeforeEach(func() {
			s = hpa.Status()
		})

		When("DesiredReplicas > CurrentReplicas", func() {
			BeforeEach(func() {
				o := hpa.Object()
				o.Status.DesiredReplicas = 5
				o.Status.CurrentReplicas = 3
			})

			It("returns expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for HPA to complete a scale up, current: 3 desired: 5"))
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("Waiting for HPA to complete a scale up, current: 3 desired: 5"))
			})
		})

		When("DesiredReplicas < CurrentReplicas", func() {
			BeforeEach(func() {
				o := hpa.Object()
				o.Status.DesiredReplicas = 3
				o.Status.CurrentReplicas = 5
			})

			It("returns expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for HPA to complete a scale down, current: 5 desired: 3"))
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("Waiting for HPA to complete a scale down, current: 5 desired: 3"))
			})
		})
	})
})
