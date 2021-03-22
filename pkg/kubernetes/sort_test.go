package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sort", func() {
	var (
		c               Controller
		err             error
		m               []map[string]interface{}
		sortedManifests []map[string]interface{}
	)

	BeforeEach(func() {
		c = NewController()
	})

	JustBeforeEach(func() {
		sortedManifests, err = c.SortManifests(m)
	})

	Describe("#SortManifests", func() {
		When("the manifests are unsorted", func() {
			BeforeEach(func() {
				m = []map[string]interface{}{
					{
						"kind": "ApiService",
					},
					{
						"kind": "ClusterRoleBinding",
					},
					{
						"kind": "ClusterRoleHandler",
					},
					{
						"kind": "ConfigMap",
					},
					{
						"kind": "ControllerRevision",
					},
					{
						"kind": "CronJob",
					},
					{
						"kind": "CustomResourceDefinition",
					},
					{
						"kind": "DaemonSet",
					},
					{
						"kind": "Deployment",
					},
					{
						"kind": "Event",
					},
					{
						"kind": "HorizontalPodAutoscaler",
					},
					{
						"kind": "Ingress",
					},
					{
						"kind": "Job",
					},
					{
						"kind": "LimitRange",
					},
					{
						"kind": "MutatingWebhookConfiguration",
					},
					{
						"kind": "Namespace",
					},
					{
						"kind": "NetworkPolicy",
					},
					{
						"kind": "PersistentVolumeClaim",
					},
					{
						"kind": "PersistentVolume",
					},
					{
						"kind": "PodDisruptionBudget",
					},
					{
						"kind": "Pod",
					},
					{
						"kind": "PodPreset",
					},
					{
						"kind": "PodSecurityPolicy",
					},
					{
						"kind": "ReplicaSet",
					},
					{
						"kind": "RoleBinding",
					},
					{
						"kind": "Role",
					},
					{
						"kind": "Secret",
					},
					{
						"kind": "ServiceAccount",
					},
					{
						"kind": "Service",
					},
					{
						"kind": "StatefulSet",
					},
					{
						"kind": "StorageClass",
					},
					{
						"kind": "UnregisteredClusterResource",
					},
					{
						"kind": "ValidatingWebhookConfiguration",
					},
					{
						"kind": "Unknown",
					},
				}
			})

			It("sorts the manifests", func() {
				Expect(err).To(BeNil())
				Expect(sortedManifests).To(HaveLen(len(m)))
				// priority 0
				Expect(sortedManifests[0]["kind"].(string)).To(Equal("Namespace"))
				// priority 20
				Expect(sortedManifests[1]["kind"].(string)).To(Equal("ClusterRoleHandler"))
				Expect(sortedManifests[2]["kind"].(string)).To(Equal("Role"))
				// priority 30
				Expect(sortedManifests[3]["kind"].(string)).To(Equal("ClusterRoleBinding"))
				Expect(sortedManifests[4]["kind"].(string)).To(Equal("CustomResourceDefinition"))
				Expect(sortedManifests[5]["kind"].(string)).To(Equal("RoleBinding"))
				// priority 40
				Expect(sortedManifests[6]["kind"].(string)).To(Equal("MutatingWebhookConfiguration"))
				Expect(sortedManifests[7]["kind"].(string)).To(Equal("PersistentVolume"))
				Expect(sortedManifests[8]["kind"].(string)).To(Equal("ServiceAccount"))
				Expect(sortedManifests[9]["kind"].(string)).To(Equal("StorageClass"))
				Expect(sortedManifests[10]["kind"].(string)).To(Equal("ValidatingWebhookConfiguration"))
				// priority 50
				Expect(sortedManifests[11]["kind"].(string)).To(Equal("ConfigMap"))
				Expect(sortedManifests[12]["kind"].(string)).To(Equal("PersistentVolumeClaim"))
				Expect(sortedManifests[13]["kind"].(string)).To(Equal("Secret"))
				// priority 70
				Expect(sortedManifests[14]["kind"].(string)).To(Equal("Ingress"))
				Expect(sortedManifests[15]["kind"].(string)).To(Equal("NetworkPolicy"))
				Expect(sortedManifests[16]["kind"].(string)).To(Equal("Service"))
				// priority 80
				Expect(sortedManifests[17]["kind"].(string)).To(Equal("ApiService"))
				// priority 90
				Expect(sortedManifests[18]["kind"].(string)).To(Equal("LimitRange"))
				Expect(sortedManifests[19]["kind"].(string)).To(Equal("PodDisruptionBudget"))
				Expect(sortedManifests[20]["kind"].(string)).To(Equal("PodPreset"))
				Expect(sortedManifests[21]["kind"].(string)).To(Equal("PodSecurityPolicy"))
				// priority 100
				Expect(sortedManifests[22]["kind"].(string)).To(Equal("CronJob"))
				Expect(sortedManifests[23]["kind"].(string)).To(Equal("DaemonSet"))
				Expect(sortedManifests[24]["kind"].(string)).To(Equal("Deployment"))
				Expect(sortedManifests[25]["kind"].(string)).To(Equal("Job"))
				Expect(sortedManifests[26]["kind"].(string)).To(Equal("Pod"))
				Expect(sortedManifests[27]["kind"].(string)).To(Equal("ReplicaSet"))
				Expect(sortedManifests[28]["kind"].(string)).To(Equal("StatefulSet"))
				// priority 110
				Expect(sortedManifests[29]["kind"].(string)).To(Equal("HorizontalPodAutoscaler"))
				// priority 1000
				Expect(sortedManifests[30]["kind"].(string)).To(Equal("ControllerRevision"))
				Expect(sortedManifests[31]["kind"].(string)).To(Equal("Event"))
				Expect(sortedManifests[32]["kind"].(string)).To(Equal("UnregisteredClusterResource"))
				Expect(sortedManifests[33]["kind"].(string)).To(Equal("Unknown"))
			})
		})
	})
})
