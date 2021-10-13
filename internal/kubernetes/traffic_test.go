package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Traffic", func() {
	var (
		err          error
		fakeResource unstructured.Unstructured
		lbs          []string
	)

	Context("#LoadBalancers", func() {
		JustBeforeEach(func() {
			lbs, err = LoadBalancers(fakeResource)
		})

		When("annotation is missing", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "ReplicaSet",
					},
				}
			})

			It("returns no load balancers", func() {
				Expect(err).To(BeNil())
				Expect(lbs).To(HaveLen(0))
			})
		})

		When("the annotation is not an array", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "ReplicaSet",
						"metadata": map[string]interface{}{
							"name":      "test-name",
							"namespace": "test-namespace",
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "string",
							},
						},
					},
				}
			})

			It("errors", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error unmarshaling annotation 'traffic.spinnaker.io/load-balancers' for resource (kind: ReplicaSet, name: test-name, namespace: test-namespace) into string slice: invalid character 's' looking for beginning of value"))
			})
		})

		When("the annotation is formatted incorrectly", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "ReplicaSet",
						"metadata": map[string]interface{}{
							"name":      "test-name",
							"namespace": "test-namespace",
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "[\"test-lb-service\"]",
							},
						},
					},
				}
			})

			It("errors", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to attach load balancer 'test-lb-service'. Load balancers must be specified in the form '{kind} {name}', e.g. 'service my-service'."))
			})
		})

		When("the load balancer kind is not supported by spinnaker", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "ReplicaSet",
						"metadata": map[string]interface{}{
							"name":      "test-name",
							"namespace": "test-namespace",
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "[\"service test-lb-service\",\"ingress test-lb-service\"]",
							},
						},
					},
				}
			})

			It("errors", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("No support for load balancing via ingress exists in Spinnaker."))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				fakeResource = unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind": "ReplicaSet",
						"metadata": map[string]interface{}{
							"name":      "test-name",
							"namespace": "test-namespace",
							"annotations": map[string]interface{}{
								"traffic.spinnaker.io/load-balancers": "[\"service test-lb-service\",\"service test-lb-service2\"]",
							},
						},
					},
				}
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(lbs).To(HaveLen(2))
				Expect(lbs[0]).To(Equal("service test-lb-service"))
				Expect(lbs[1]).To(Equal("service test-lb-service2"))
			})
		})
	})
})
