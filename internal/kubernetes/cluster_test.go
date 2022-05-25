package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster", func() {
	Describe("#Cluster", func() {
		var (
			kind    string
			cluster string
		)

		JustBeforeEach(func() {
			cluster = Cluster(kind, "test-name")
		})

		When("the kind is daemonSet", func() {
			BeforeEach(func() {
				kind = "DaemonSet"
			})

			It("sets the cluster", func() {
				Expect(cluster).To(Equal("daemonSet test-name"))
			})
		})

		When("the kind is deployment", func() {
			BeforeEach(func() {
				kind = "Deployment"
			})

			It("sets the cluster", func() {
				Expect(cluster).To(Equal("deployment test-name"))
			})
		})

		When("the kind is ingress", func() {
			BeforeEach(func() {
				kind = "Ingress"
			})

			It("sets the cluster", func() {
				Expect(cluster).To(Equal("ingress test-name"))
			})
		})

		When("the kind is replicaSet", func() {
			BeforeEach(func() {
				kind = "ReplicaSet"
			})

			It("sets the cluster", func() {
				Expect(cluster).To(Equal("replicaSet test-name"))
			})
		})

		When("the kind is service", func() {
			BeforeEach(func() {
				kind = "Service"
			})

			It("sets the cluster", func() {
				Expect(cluster).To(Equal("service test-name"))
			})
		})

		When("the kind is statefulSet", func() {
			BeforeEach(func() {
				kind = "StatefulSet"
			})

			It("sets the cluster", func() {
				Expect(cluster).To(Equal("statefulSet test-name"))
			})
		})

		When("the kind is not a cluster type", func() {
			BeforeEach(func() {
				kind = "Pod"
			})

			It("does not set the cluster", func() {
				Expect(cluster).To(BeEmpty())
			})
		})
	})
})
