package kubernetes_test

import (
	. "github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
)

var _ = Describe("Deployment", func() {
	var (
		deployment Deployment
		err        error
	)

	BeforeEach(func() {
		d := map[string]interface{}{}
		deployment = NewDeployment(d)
	})

	Describe("#ToUnstructured", func() {
		BeforeEach(func() {
			_, err = deployment.ToUnstructured()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#Object", func() {
		var d *v1.Deployment

		BeforeEach(func() {
			d = deployment.Object()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(d).ToNot(BeNil())
			})
		})
	})

	Describe("#AnnotateTemplate", func() {
		BeforeEach(func() {
			deployment.AnnotateTemplate("test", "value")
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				o := deployment.Object()
				annotations := o.Spec.Template.ObjectMeta.Annotations
				Expect(annotations["test"]).To(Equal("value"))
			})
		})
	})

	Describe("#GetSpec", func() {
		BeforeEach(func() {
			_ = deployment.GetSpec()
		})

		It("succeeds", func() {
		})
	})

	Describe("#LabelTemplate", func() {
		BeforeEach(func() {
			deployment.LabelTemplate("test", "value")
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				o := deployment.Object()
				labels := o.Spec.Template.ObjectMeta.Labels
				Expect(labels["test"]).To(Equal("value"))
			})
		})
	})

	Describe("#LabelTemplateIfNotExists", func() {
		JustBeforeEach(func() {
			deployment.LabelTemplateIfNotExists("test", "value")
		})

		When("the label exists", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Spec.Template.ObjectMeta.Labels = map[string]string{
					"test": "taken",
				}
			})

			It("does not label the template", func() {
				o := deployment.Object()
				labels := o.Spec.Template.ObjectMeta.Labels
				Expect(labels["test"]).To(Equal("taken"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				o := deployment.Object()
				labels := o.Spec.Template.ObjectMeta.Labels
				Expect(labels["test"]).To(Equal("value"))
			})
		})
	})

	Describe("#SetReplicas", func() {
		BeforeEach(func() {
			replicas := int32(4)
			deployment.SetReplicas(&replicas)
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				o := deployment.Object()
				replicas := o.Spec.Replicas
				Expect(*replicas).To(Equal(int32(4)))
			})
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		BeforeEach(func() {
			replicas := int32(4)
			deployment.SetReplicas(&replicas)
		})

		JustBeforeEach(func() {
			s = deployment.Status()
		})

		When("the generations do not match", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.ObjectMeta.Generation = 100
				o.Status.ObservedGeneration = 99
			})

			It("returns status unstable", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for status generation to match updated object generation"))
			})
		})

		When("the deployment is paused", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.Conditions = []v1.DeploymentCondition{
					{
						Reason: "deploymentPaused",
					},
				}
			})

			It("returns paused state", func() {
				Expect(s.Paused.State).To(BeTrue())
			})
		})

		When("the deployment has a condition type of available and status false", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.Conditions = []v1.DeploymentCondition{
					{
						Type:   "available",
						Status: "false",
						Reason: "test reason",
					},
				}
			})

			It("returns the expected status", func() {
				Expect(s.Available.State).To(BeFalse())
				Expect(s.Available.Message).To(Equal("test reason"))
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be updated"))
			})
		})

		When("the deployment has condition of type progressing with deadline exceeded", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.Conditions = []v1.DeploymentCondition{
					{
						Type:   "progressing",
						Reason: "progressdeadlineexceeded",
					},
				}
			})

			It("returns the expected status", func() {
				Expect(s.Failed.State).To(BeTrue())
			})
		})

		When("updated replicas is less than desired", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.UpdatedReplicas = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be updated"))
			})
		})

		When("there are more status replicas than updated replicas", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.UpdatedReplicas = int32(4)
				o.Status.Replicas = int32(6)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for old replicas to finish termination"))
			})
		})

		When("there are less available replicas than desired replicas", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.UpdatedReplicas = int32(4)
				o.Status.AvailableReplicas = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be available"))
			})
		})

		When("there are less ready replicas than desired replicas", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.UpdatedReplicas = int32(4)
				o.Status.AvailableReplicas = int32(4)
				o.Status.ReadyReplicas = int32(2)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeFalse())
				Expect(s.Stable.Message).To(Equal("Waiting for all replicas to be ready"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				o := deployment.Object()
				o.Status.UpdatedReplicas = int32(4)
				o.Status.AvailableReplicas = int32(4)
				o.Status.ReadyReplicas = int32(4)
			})

			It("returns the expected status", func() {
				Expect(s.Stable.State).To(BeTrue())
			})
		})
	})
})
