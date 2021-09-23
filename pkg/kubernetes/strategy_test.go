package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Strategy", func() {

	var (
		err               error
		fakeResource      unstructured.Unstructured
		maxVersionHistory int
		recreate          bool
		replace           bool
		useSourceCapacity bool
	)

	Context("#GetMaxVersionHistory", func() {
		JustBeforeEach(func() {
			maxVersionHistory, err = GetMaxVersionHistory(fakeResource)
		})

		When("annotation is missing", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
					},
				}
			})

			It("returns 0", func() {
				Expect(maxVersionHistory).To(Equal(0))
				Expect(err).To(BeNil())
			})
		})

		When("annotation is set to a non-integer", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/max-version-history": "one",
							},
						},
					},
				}
			})

			It("errors", func() {
				Expect(err).ToNot(BeNil())
				Expect(maxVersionHistory).To(Equal(0))
			})
		})

		When("annotation is set to an integer", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/max-version-history": "5",
							},
						},
					},
				}
			})

			It("succeeds", func() {
				Expect(maxVersionHistory).To(Equal(5))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("#Recreate", func() {
		JustBeforeEach(func() {
			recreate = Recreate(fakeResource)
		})

		When("annotation is missing", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
					},
				}
			})

			It("returns false", func() {
				Expect(recreate).To(Equal(false))
			})
		})

		When("annotation is set to false", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/recreate": "false",
							},
						},
					},
				}
			})

			It("returns false", func() {
				Expect(recreate).To(Equal(false))
			})
		})

		When("annotation is set to true", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/recreate": "true",
							},
						},
					},
				}
			})

			It("returns true", func() {
				Expect(recreate).To(Equal(true))
			})
		})
	})

	Context("#Replace", func() {
		JustBeforeEach(func() {
			replace = Replace(fakeResource)
		})

		When("annotation is missing", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
					},
				}
			})

			It("returns false", func() {
				Expect(replace).To(Equal(false))
			})
		})

		When("annotation is set to false", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/replace": "false",
							},
						},
					},
				}
			})

			It("returns false", func() {
				Expect(replace).To(Equal(false))
			})
		})

		When("annotation is set to true", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/replace": "true",
							},
						},
					},
				}
			})

			It("returns true", func() {
				Expect(replace).To(Equal(true))
			})
		})
	})

	Context("#UseSourceCapacity", func() {
		JustBeforeEach(func() {
			useSourceCapacity = UseSourceCapacity(fakeResource)
		})

		When("annotation is missing", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
					},
				}
			})

			It("returns false", func() {
				Expect(useSourceCapacity).To(Equal(false))
			})
		})

		When("annotation is false", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "Deployment",
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"strategy.spinnaker.io/use-source-capacity": "false",
							},
						},
					},
				}
			})

			It("returns false", func() {
				Expect(useSourceCapacity).To(Equal(false))
			})
		})

		When("annotation is true", func() {
			When("kind is pod", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "Pod",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/use-source-capacity": "true",
								},
							},
						},
					}
				})

				It("returns false", func() {
					Expect(UseSourceCapacity(fakeResource)).To(Equal(false))
				})
			})

			When("kind is deployment", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "Deployment",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/use-source-capacity": "true",
								},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(UseSourceCapacity(fakeResource)).To(Equal(true))
				})
			})

			When("kind is replicaSet", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "ReplicaSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/use-source-capacity": "true",
								},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(UseSourceCapacity(fakeResource)).To(Equal(true))
				})
			})

			When("kind is statefuleSet", func() {
				BeforeEach(func() {
					fakeResource = unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind": "StatefulSet",
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{
									"strategy.spinnaker.io/use-source-capacity": "true",
								},
							},
						},
					}
				})

				It("returns true", func() {
					Expect(UseSourceCapacity(fakeResource)).To(Equal(true))
				})
			})
		})
	})
})
