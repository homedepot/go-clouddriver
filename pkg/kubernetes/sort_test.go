package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Sort", func() {
	var (
		m               []map[string]interface{}
		ul              []unstructured.Unstructured
		sortedManifests []unstructured.Unstructured
	)

	JustBeforeEach(func() {
		sortedManifests = SortManifests(ul)
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
				for _, v := range m {
					u, err := ToUnstructured(v)
					Expect(err).To(BeNil())
					ul = append(ul, u)
				}
			})

			It("sorts the manifests", func() {
				Expect(sortedManifests).To(HaveLen(len(m)))
				// priority 0
				Expect(sortedManifests[0].GetKind()).To(Equal("Namespace"))
				// priority 20
				Expect(sortedManifests[1].GetKind()).To(Equal("ClusterRoleHandler"))
				Expect(sortedManifests[2].GetKind()).To(Equal("Role"))
				// priority 30
				Expect(sortedManifests[3].GetKind()).To(Equal("ClusterRoleBinding"))
				Expect(sortedManifests[4].GetKind()).To(Equal("CustomResourceDefinition"))
				Expect(sortedManifests[5].GetKind()).To(Equal("RoleBinding"))
				// priority 40
				Expect(sortedManifests[6].GetKind()).To(Equal("MutatingWebhookConfiguration"))
				Expect(sortedManifests[7].GetKind()).To(Equal("PersistentVolume"))
				Expect(sortedManifests[8].GetKind()).To(Equal("ServiceAccount"))
				Expect(sortedManifests[9].GetKind()).To(Equal("StorageClass"))
				Expect(sortedManifests[10].GetKind()).To(Equal("ValidatingWebhookConfiguration"))
				// priority 50
				Expect(sortedManifests[11].GetKind()).To(Equal("ConfigMap"))
				Expect(sortedManifests[12].GetKind()).To(Equal("PersistentVolumeClaim"))
				Expect(sortedManifests[13].GetKind()).To(Equal("Secret"))
				// priority 70
				Expect(sortedManifests[14].GetKind()).To(Equal("Ingress"))
				Expect(sortedManifests[15].GetKind()).To(Equal("NetworkPolicy"))
				Expect(sortedManifests[16].GetKind()).To(Equal("Service"))
				// priority 80
				Expect(sortedManifests[17].GetKind()).To(Equal("ApiService"))
				// priority 90
				Expect(sortedManifests[18].GetKind()).To(Equal("LimitRange"))
				Expect(sortedManifests[19].GetKind()).To(Equal("PodDisruptionBudget"))
				Expect(sortedManifests[20].GetKind()).To(Equal("PodPreset"))
				Expect(sortedManifests[21].GetKind()).To(Equal("PodSecurityPolicy"))
				// priority 100
				Expect(sortedManifests[22].GetKind()).To(Equal("CronJob"))
				Expect(sortedManifests[23].GetKind()).To(Equal("DaemonSet"))
				Expect(sortedManifests[24].GetKind()).To(Equal("Deployment"))
				Expect(sortedManifests[25].GetKind()).To(Equal("Job"))
				Expect(sortedManifests[26].GetKind()).To(Equal("Pod"))
				Expect(sortedManifests[27].GetKind()).To(Equal("ReplicaSet"))
				Expect(sortedManifests[28].GetKind()).To(Equal("StatefulSet"))
				// priority 110
				Expect(sortedManifests[29].GetKind()).To(Equal("HorizontalPodAutoscaler"))
				// priority 1000
				Expect(sortedManifests[30].GetKind()).To(Equal("ControllerRevision"))
				Expect(sortedManifests[31].GetKind()).To(Equal("Event"))
				Expect(sortedManifests[32].GetKind()).To(Equal("UnregisteredClusterResource"))
				Expect(sortedManifests[33].GetKind()).To(Equal("Unknown"))
			})
		})
	})
})
