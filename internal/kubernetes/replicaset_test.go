package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
)

var _ = Describe("Replicaset", func() {
	var (
		rs *ReplicaSet
	)

	BeforeEach(func() {
		r := map[string]interface{}{}
		rs = NewReplicaSet(r)
	})

	Describe("#Object", func() {
		var r *v1.ReplicaSet

		BeforeEach(func() {
			r = rs.Object()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(r).ToNot(BeNil())
			})
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		BeforeEach(func() {
			replicas := int32(4)
			rs.Object().Spec.Replicas = &replicas
		})

		JustBeforeEach(func() {
			s = rs.Status()
		})

		When("there are more desired replicas than fully labeled replicas", func() {
			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be fully-labeled"))
			})
		})

		When("there are more desired replicas than ready replicas", func() {
			BeforeEach(func() {
				o := rs.Object()
				o.Status.FullyLabeledReplicas = int32(4)
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be ready"))
			})
		})

		When("there are more desired replicas than available", func() {
			BeforeEach(func() {
				o := rs.Object()
				o.Status.FullyLabeledReplicas = int32(4)
				o.Status.ReadyReplicas = int32(4)
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be available"))
			})
		})

		When("the generations do not match", func() {
			BeforeEach(func() {
				o := rs.Object()
				o.Status.FullyLabeledReplicas = int32(4)
				o.Status.ReadyReplicas = int32(4)
				o.Status.AvailableReplicas = int32(4)
				o.ObjectMeta.Generation = 100
				o.Status.ObservedGeneration = 99
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for replicaset spec update to be observed"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				o := rs.Object()
				o.Status.FullyLabeledReplicas = int32(4)
				o.Status.ReadyReplicas = int32(4)
				o.Status.AvailableReplicas = int32(4)
				o.ObjectMeta.Generation = 100
				o.Status.ObservedGeneration = 100
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeTrue())
			})
		})
	})
})
