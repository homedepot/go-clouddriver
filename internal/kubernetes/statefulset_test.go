package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
)

var _ = Describe("Statefulset", func() {
	var (
		ss *StatefulSet
	)

	BeforeEach(func() {
		s := map[string]interface{}{}
		ss = NewStatefulSet(s)
	})

	Describe("#Object", func() {
		var s *v1.StatefulSet

		BeforeEach(func() {
			s = ss.Object()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(s).ToNot(BeNil())
			})
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		BeforeEach(func() {
			replicas := int32(4)
			o := ss.Object()
			o.Spec.Replicas = &replicas
			o.Status.Replicas = replicas
			o.Status.ReadyReplicas = replicas
			o.Status.UpdatedReplicas = replicas
			o.Status.CurrentReplicas = replicas
		})

		JustBeforeEach(func() {
			s = ss.Status()
		})

		When("the update strategy type is OnDelete", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Spec.UpdateStrategy.Type = "OnDelete"
			})

			It("returns the expected status", func() {
				Expect(s).To(Equal(manifest.DefaultStatus))
			})
		})

		When("there is no reported status", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Status = v1.StatefulSetStatus{}
			})

			It("returns the expected status", func() {
				Expect(s).To(Equal(manifest.NoneReported))
			})
		})

		When("the generation does not match", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.ObjectMeta.Generation = int64(99)
				o.Status.ObservedGeneration = int64(100)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for status generation to match updated object generation"))
			})
		})

		When("there are more desired replicas than existing replicas", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Status.Replicas = 3
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for at least the desired replica count to be met"))
			})
		})

		When("there are more desired replicas than ready replicas", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Status.ReadyReplicas = int32(3)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all updated replicas to be ready"))
			})
		})

		Context("when the update type is a rolling update", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Spec.UpdateStrategy.Type = "RollingUpdate"
				rollingUpdate := v1.RollingUpdateStatefulSetStrategy{}
				o.Spec.UpdateStrategy.RollingUpdate = &rollingUpdate
			})

			When("the partitioned rollout has not finished", func() {
				BeforeEach(func() {
					partition := int32(2)
					o := ss.Object()
					o.Status.UpdatedReplicas = 1
					o.Status.CurrentReplicas = 3
					o.Spec.UpdateStrategy.RollingUpdate.Partition = &partition
				})

				It("returns the expected status", func() {
					Expect(s.Stable.State).To(BeFalse())
					Expect(s.Stable.Message).To(Equal("Waiting for partitioned rollout to finish"))
				})
			})

			When("the partitioned rollout has finished", func() {
				BeforeEach(func() {
					partition := int32(2)
					o := ss.Object()
					o.Status.UpdatedReplicas = 2
					o.Spec.UpdateStrategy.RollingUpdate.Partition = &partition
				})

				It("returns the expected status", func() {
					Expect(s.Stable.State).To(BeTrue())
					Expect(s.Stable.Message).To(Equal("Partitioned roll out complete"))
				})
			})
		})

		When("the desired replicas is more than the current replicas", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Status.CurrentReplicas = 2
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all updated replicas to be scheduled"))
			})
		})

		When("the current revision is not equal to the update revision", func() {
			BeforeEach(func() {
				o := ss.Object()
				o.Status.UpdateRevision = "100"
				o.Status.CurrentRevision = "99"
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for the updated revision to match the current revision"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(s).To(Equal(manifest.DefaultStatus))
			})
		})
	})
})
