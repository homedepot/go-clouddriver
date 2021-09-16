package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
)

var _ = Describe("DaemonSet", func() {
	var (
		daemonSet *DaemonSet
	)

	BeforeEach(func() {
		ds := map[string]interface{}{}
		daemonSet = NewDaemonSet(ds)
	})

	Describe("#Status", func() {
		var s manifest.Status

		BeforeEach(func() {
			daemonSet.Object().Status.DesiredNumberScheduled = 4
		})

		JustBeforeEach(func() {
			s = daemonSet.Status()
		})

		When("there is no status", func() {
			BeforeEach(func() {
				ds := map[string]interface{}{}
				daemonSet = NewDaemonSet(ds)
			})

			It("returns status unstable and unavailable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("No status reported yet"))
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("No availability reported"))

			})
		})

		When("the update stategy is rolling update", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.Spec.UpdateStrategy = v1.DaemonSetUpdateStrategy{
					Type: v1.RollingUpdateDaemonSetStrategyType,
				}
			})

			It("returns status stable and available", func() {
				Expect(s.Stable.State).To(BeTrue())
				Expect(s.Available.State).To(BeTrue())
			})
		})

		When("the generations do not match", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.ObjectMeta.Generation = 100
				o.Status.ObservedGeneration = 99
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for status generation to match updated object generation"))
			})
		})

		When("scheduled replicas is less than desired", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.Status.CurrentNumberScheduled = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be scheduled"))
			})
		})

		When("updated replicas is less than desired", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.Status.CurrentNumberScheduled = int32(4)
				o.Status.UpdatedNumberScheduled = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all updated replicas to be scheduled"))
			})
		})

		When("there are less available replicas than desireds replicas", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.Status.CurrentNumberScheduled = int32(4)
				o.Status.UpdatedNumberScheduled = int32(4)
				o.Status.NumberAvailable = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be available"))
			})
		})

		When("there are less ready replicas than desireds replicas", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.Status.CurrentNumberScheduled = int32(4)
				o.Status.UpdatedNumberScheduled = int32(4)
				o.Status.NumberAvailable = int32(4)
				o.Status.NumberReady = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be ready"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				o := daemonSet.Object()
				o.Status.CurrentNumberScheduled = int32(4)
				o.Status.UpdatedNumberScheduled = int32(4)
				o.Status.NumberAvailable = int32(4)
				o.Status.NumberReady = int32(4)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeTrue())
			})
		})
	})
})
