package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
)

var _ = Describe("Job", func() {
	var (
		job *Job
	)

	BeforeEach(func() {
		j := map[string]interface{}{}
		job = NewJob(j)
	})

	Describe("#State", func() {
		var s string

		JustBeforeEach(func() {
			s = job.State()
		})

		When("the job has not completed", func() {
			BeforeEach(func() {
				o := job.Object()
				o.Status.CompletionTime = nil
			})

			It("returns expected state", func() {
				Expect(s).To(Equal("Running"))
			})
		})

		When("the job has failed", func() {
			BeforeEach(func() {
				completions := int32(1)
				o := job.Object()
				o.Status.CompletionTime = &metav1.Time{}
				o.Spec.Completions = &completions
				o.Status.Conditions = []v1.JobCondition{
					{
						Type: "Failed",
					},
				}
			})

			It("returns expected state", func() {
				Expect(s).To(Equal("Failed"))
			})
		})

		When("the job is partially successful", func() {
			BeforeEach(func() {
				completions := int32(1)
				o := job.Object()
				o.Status.CompletionTime = &metav1.Time{}
				o.Spec.Completions = &completions
			})

			It("returns expected state", func() {
				Expect(s).To(Equal("Running"))
			})
		})

		When("the job succeeded", func() {
			BeforeEach(func() {
				completions := int32(1)
				o := job.Object()
				o.Status.CompletionTime = &metav1.Time{}
				o.Status.Succeeded = int32(1)
				o.Spec.Completions = &completions
			})

			It("returns expected state", func() {
				Expect(s).To(Equal("Succeeded"))
			})
		})
	})

	Describe("#Status", func() {
		var s manifest.Status

		BeforeEach(func() {
			completions := int32(1)
			o := job.Object()
			o.Status.Succeeded = 1
			o.Spec.Completions = &completions
		})

		JustBeforeEach(func() {
			s = job.Status()
		})

		Context("succeeded pods is less than completions", func() {
			BeforeEach(func() {
				o := job.Object()
				completions := int32(2)
				o.Status.Succeeded = 1
				o.Spec.Completions = &completions
			})

			When("there is a failed condition", func() {
				BeforeEach(func() {
					o := job.Object()
					o.Status.Conditions = []v1.JobCondition{
						{
							Type:    v1.JobFailed,
							Message: "Some failure message",
						},
					}
				})

				It("returns status failed", func() {
					Expect(s.Failed.State).To(BeTrue())
					Expect(s.Failed.Message).To(Equal("Some failure message"))
				})
			})

			When("the job is not finished", func() {
				It("returns the expected status", func() {
					Expect(s.Stable.State).To(BeFalse())
					Expect(s.Stable.Message).To(Equal("Waiting for jobs to finish"))
				})
			})
		})

		When("it succeeeds", func() {
			It("succeeds", func() {
				Expect(s.Stable.State).To(BeTrue())
				Expect(s.Available.State).To(BeTrue())
			})
		})
	})
})
