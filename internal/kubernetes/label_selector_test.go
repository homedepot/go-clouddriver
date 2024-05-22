package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
)

var _ = Describe("Label Selector", func() {
	var (
		err         error
		op          string
		key         string
		values      []string
		requirement *labels.Requirement
	)

	Context("#NewRequirement", func() {
		JustBeforeEach(func() {
			requirement, err = NewRequirement(op, key, values)
		})

		When("the op is invalid", func() {
			Context("the op is empty", func() {
				BeforeEach(func() {
					op = ""
					key = "key1"
					values = []string{"value1"}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(err.Error()).To(Equal("operator '' is not recognized"))
					Expect(requirement).To(BeNil())
				})
			})

			Context("the op is unknown", func() {
				BeforeEach(func() {
					op = "ANY"
					key = "key1"
					values = []string{"value1"}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(err.Error()).To(Equal("operator 'ANY' is not recognized"))
					Expect(requirement).To(BeNil())
				})
			})
		})

		When("the op is EQUALS", func() {
			BeforeEach(func() {
				op = "EQUALS"
				key = "key1"
				values = []string{"value1"}
			})

			Context("the key is empty", func() {
				BeforeEach(func() {
					key = ""
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal(""))
					Expect(requirement.Operator()).Should(Equal(selection.Equals))
					Expect(requirement.String()).To(Equal("=value1"))
				})
			})

			Context("values is nil", func() {
				BeforeEach(func() {
					values = nil
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.Equals))
					Expect(requirement.String()).To(Equal("key1="))
				})
			})

			Context("values is empty", func() {
				BeforeEach(func() {
					values = []string{}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.Equals))
					Expect(requirement.String()).To(Equal("key1="))
				})
			})

			Context("there is more than one value", func() {
				BeforeEach(func() {
					values = []string{"value1", "value2"}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.Equals))
					Expect(requirement.String()).To(Equal("key1=value1,value2"))
				})
			})

			Context("key and values are set", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.String()).To(Equal("key1=value1"))
				})
			})
		})

		When("the op is NOT_EQUALS", func() {
			BeforeEach(func() {
				op = "NOT_EQUALS"
				key = "key1"
				values = []string{"value1"}
			})

			Context("the key is empty", func() {
				BeforeEach(func() {
					key = ""
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal(""))
					Expect(requirement.Operator()).Should(Equal(selection.NotEquals))
					Expect(requirement.String()).To(Equal("!=value1"))
				})
			})

			Context("values is nil", func() {
				BeforeEach(func() {
					values = nil
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotEquals))
					Expect(requirement.String()).To(Equal("key1!="))
				})
			})

			Context("values is empty", func() {
				BeforeEach(func() {
					values = []string{}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotEquals))
					Expect(requirement.String()).To(Equal("key1!="))
				})
			})

			Context("there is more than one value", func() {
				BeforeEach(func() {
					values = []string{"value1", "value2"}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotEquals))
					Expect(requirement.String()).To(Equal("key1!=value1,value2"))
				})
			})

			Context("key and values are set", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotEquals))
					Expect(requirement.String()).To(Equal("key1!=value1"))
				})
			})
		})

		When("the op is EXISTS", func() {
			BeforeEach(func() {
				op = "EXISTS"
				key = "key1"
				values = []string{}
			})

			Context("the key is empty", func() {
				BeforeEach(func() {
					key = ""
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal(""))
					Expect(requirement.Operator()).Should(Equal(selection.Exists))
					Expect(requirement.String()).To(Equal(""))
				})
			})

			Context("values is set", func() {
				BeforeEach(func() {
					values = []string{"value1"}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.Exists))
					Expect(requirement.String()).To(Equal("key1"))
				})
			})

			Context("just key is set", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.Exists))
					Expect(requirement.String()).To(Equal("key1"))
				})
			})
		})

		When("the op is NOT_EXISTS", func() {
			BeforeEach(func() {
				op = "NOT_EXISTS"
				key = "key1"
				values = []string{}
			})

			Context("the key is empty", func() {
				BeforeEach(func() {
					key = ""
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal(""))
					Expect(requirement.Operator()).Should(Equal(selection.DoesNotExist))
					Expect(requirement.String()).To(Equal("!"))
				})
			})

			Context("values is set", func() {
				BeforeEach(func() {
					values = []string{"value1"}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.DoesNotExist))
					Expect(requirement.String()).To(Equal("!key1"))
				})
			})

			Context("just key is set", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.DoesNotExist))
					Expect(requirement.String()).To(Equal("!key1"))
				})
			})
		})

		When("the op is CONTAINS", func() {
			BeforeEach(func() {
				op = "CONTAINS"
				key = "key1"
				values = []string{"value1", "value2"}
			})

			Context("the key is empty", func() {
				BeforeEach(func() {
					key = ""
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal(""))
					Expect(requirement.Operator()).Should(Equal(selection.In))
					Expect(requirement.String()).To(Equal(" in (value1,value2)"))
				})
			})

			Context("values is nil", func() {
				BeforeEach(func() {
					values = nil
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.In))
					Expect(requirement.String()).To(Equal("key1 in ()"))
				})
			})

			Context("values is empty", func() {
				BeforeEach(func() {
					values = []string{}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.In))
					Expect(requirement.String()).To(Equal("key1 in ()"))
				})
			})

			Context("there is one value", func() {
				BeforeEach(func() {
					values = []string{"value1"}
				})

				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.In))
					Expect(requirement.String()).To(Equal("key1 in (value1)"))
				})
			})

			Context("there is more than one value", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.In))
					Expect(requirement.String()).To(Equal("key1 in (value1,value2)"))
				})
			})
		})

		When("the op is NOT_CONTAINS", func() {
			BeforeEach(func() {
				op = "NOT_CONTAINS"
				key = "key1"
				values = []string{"value1", "value2"}
			})

			Context("the key is empty", func() {
				BeforeEach(func() {
					key = ""
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal(""))
					Expect(requirement.Operator()).Should(Equal(selection.NotIn))
					Expect(requirement.String()).To(Equal(" notin (value1,value2)"))
				})
			})

			Context("values is nil", func() {
				BeforeEach(func() {
					values = nil
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotIn))
					Expect(requirement.String()).To(Equal("key1 notin ()"))
				})
			})

			Context("values is empty", func() {
				BeforeEach(func() {
					values = []string{}
				})

				It("errors", func() {
					Expect(err).To(Not(BeNil()))
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotIn))
					Expect(requirement.String()).To(Equal("key1 notin ()"))
				})
			})

			Context("there is one value", func() {
				BeforeEach(func() {
					values = []string{"value1"}
				})

				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotIn))
					Expect(requirement.String()).To(Equal("key1 notin (value1)"))
				})
			})

			Context("there is more than one value", func() {
				It("succeeds", func() {
					Expect(err).To(BeNil())
					Expect(requirement.Key()).Should(Equal("key1"))
					Expect(requirement.Operator()).Should(Equal(selection.NotIn))
					Expect(requirement.String()).To(Equal("key1 notin (value1,value2)"))
				})
			})
		})
	})

	Context("#DefaultLabelSelector", func() {
		var selector string

		JustBeforeEach(func() {
			selector = DefaultLabelSelector()
		})

		It("succeeds", func() {
			Expect(selector).To(Equal("app.kubernetes.io/managed-by in (spinnaker,spinnaker-operator)"))
		})
	})
})
