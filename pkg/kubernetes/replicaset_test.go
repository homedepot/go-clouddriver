package kubernetes_test

import (
	. "github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Replicaset", func() {
	var (
		rs  ReplicaSet
		err error
	)

	BeforeEach(func() {
		r := map[string]interface{}{}
		rs = NewReplicaSet(r)
	})

	Describe("#ToUnstructured", func() {
		BeforeEach(func() {
			_, err = rs.ToUnstructured()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#GetReplicaSetSpec", func() {
		BeforeEach(func() {
			_ = rs.GetReplicaSetSpec()
		})

		It("succeeds", func() {
		})
	})

	Describe("#GetReplicaSetStatus", func() {
		BeforeEach(func() {
			_ = rs.GetReplicaSetStatus()
		})

		It("succeeds", func() {
		})
	})

	Describe("#ListImages", func() {
		images := []string{}

		BeforeEach(func() {
			o := rs.Object()
			containers := []corev1.Container{
				{
					Image: "test-image1",
				},
				{
					Image: "test-image2",
				},
			}
			o.Spec.Template.Spec.Containers = containers
			images = rs.ListImages()
		})

		It("succeeds", func() {
			Expect(images).To(HaveLen(2))
		})
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

	Describe("#AnnotateTemplate", func() {
		BeforeEach(func() {
			rs.AnnotateTemplate("test", "value")
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				o := rs.Object()
				annotations := o.Spec.Template.ObjectMeta.Annotations
				Expect(annotations["test"]).To(Equal("value"))
			})
		})
	})

	Describe("#LabelTemplate", func() {
		BeforeEach(func() {
			rs.LabelTemplate("test", "value")
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				o := rs.Object()
				labels := o.Spec.Template.ObjectMeta.Labels
				Expect(labels["test"]).To(Equal("value"))
			})
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		BeforeEach(func() {
			replicas := int32(4)
			rs.SetReplicas(&replicas)
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
