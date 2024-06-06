package kubernetes

import (
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// Priorities of Kubernetes resources are defined in the source code here:
	// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/KubernetesHandler.java#L129
	lowestPriority                       = 1000
	workloadAttachmentPriority           = 110
	workloadControllerPriority           = 100
	workloadPriority                     = 100
	workloadModifierPriority             = 90
	pdbPriority                          = 90
	apiServicePriority                   = 80
	networkResourcePriority              = 70
	mountableDataPriority                = 50
	mountableDataBackingResourcePriority = 40
	serviceAccountPriority               = 40
	storageClassPriority                 = 40
	admissionPriority                    = 40
	resourceDefinitionPriority           = 30
	roleBindingPriority                  = 30
	rolePriority                         = 20
	namespacePriority                    = 0
)

var (
	// Define the priorities in a case insensitive map.
	// A given kind's priority can be found in each of the Kubernetes<KIND>Handler.java files here:
	// https://github.com/spinnaker/clouddriver/tree/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler
	priorities = map[string]int{
		strings.ToLower("ApiService"):         apiServicePriority,
		strings.ToLower("ClusterRoleBinding"): roleBindingPriority,
		strings.ToLower("ClusterRoleHandler"): rolePriority,
		strings.ToLower("ConfigMap"):          mountableDataPriority,
		// Controller revisions cannot be deployed.
		// See https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/KubernetesControllerRevisionHandler.java#L33
		strings.ToLower("ControllerRevision"):       lowestPriority,
		strings.ToLower("CronJob"):                  workloadControllerPriority,
		strings.ToLower("CustomResourceDefinition"): resourceDefinitionPriority,
		strings.ToLower("DaemonSet"):                workloadControllerPriority,
		strings.ToLower("Deployment"):               workloadControllerPriority,
		// Events cannot be deployed.
		// See https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/KubernetesEventHandler.java#L48
		strings.ToLower("Event"):                          lowestPriority,
		strings.ToLower("HorizontalPodAutoscaler"):        workloadAttachmentPriority,
		strings.ToLower("Ingress"):                        networkResourcePriority,
		strings.ToLower("Job"):                            workloadControllerPriority,
		strings.ToLower("LimitRange"):                     workloadModifierPriority,
		strings.ToLower("MutatingWebhookConfiguration"):   admissionPriority,
		strings.ToLower("Namespace"):                      namespacePriority,
		strings.ToLower("NetworkPolicy"):                  networkResourcePriority,
		strings.ToLower("PersistentVolumeClaim"):          mountableDataPriority,
		strings.ToLower("PersistentVolume"):               mountableDataBackingResourcePriority,
		strings.ToLower("PodDisruptionBudget"):            pdbPriority,
		strings.ToLower("Pod"):                            workloadPriority,
		strings.ToLower("PodPreset"):                      workloadModifierPriority,
		strings.ToLower("PodSecurityPolicy"):              workloadModifierPriority,
		strings.ToLower("ReplicaSet"):                     workloadControllerPriority,
		strings.ToLower("RoleBinding"):                    roleBindingPriority,
		strings.ToLower("Role"):                           rolePriority,
		strings.ToLower("Secret"):                         mountableDataPriority,
		strings.ToLower("ServiceAccount"):                 serviceAccountPriority,
		strings.ToLower("Service"):                        networkResourcePriority,
		strings.ToLower("StatefulSet"):                    workloadControllerPriority,
		strings.ToLower("StorageClass"):                   storageClassPriority,
		strings.ToLower("UnregisteredClusterResource"):    lowestPriority,
		strings.ToLower("ValidatingWebhookConfiguration"): admissionPriority,
	}
)

// SortManifests takes in a list of manifests and sorts them by the priority of their kind.
// The kind's priorities are defined above in the var 'priorities'. Lower numbered priorities
// should be deployed first.
func SortManifests(manifests []unstructured.Unstructured) []unstructured.Unstructured {
	// Map of priorities to lists of manifests.
	manifestMap := map[int][]unstructured.Unstructured{
		0:    {},
		20:   {},
		30:   {},
		40:   {},
		50:   {},
		70:   {},
		80:   {},
		90:   {},
		100:  {},
		110:  {},
		1000: {},
	}

	for _, manifest := range manifests {
		if _, ok := priorities[strings.ToLower(manifest.GetKind())]; ok {
			priority := priorities[strings.ToLower(manifest.GetKind())]
			s := manifestMap[priority]
			s = append(s, manifest)
			manifestMap[priority] = s
		} else {
			s := manifestMap[lowestPriority]
			s = append(s, manifest)
			manifestMap[lowestPriority] = s
		}
	}

	// Store the keys in slice in sorted asc order.
	keys := make([]int, len(manifestMap))
	i := 0

	for k := range manifestMap {
		keys[i] = k
		i++
	}

	sort.Ints(keys)

	sortedManifests := []unstructured.Unstructured{}
	for _, key := range keys {
		sortedManifests = append(sortedManifests, manifestMap[key]...)
	}

	return sortedManifests
}
