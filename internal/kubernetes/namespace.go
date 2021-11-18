package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SetNamespaceOnManifest updates the namespace on the given manifest.
//
// If no namespace is set on the manifest and no namespace override is passed
// in then we set the namespace to 'default'.
//
// If namespaceOverride is empty it will NOT override the namespace set
// on the manifest.
//
// We only override the namespace if the manifest is NOT cluster scoped
// (i.e. a ClusterRole) and namespaceOverride is NOT an empty string.
func SetNamespaceOnManifest(u *unstructured.Unstructured, namespaceOverride string) {
	if namespaceOverride == "" {
		setDefaultNamespaceIfScopedAndNoneSet(u)
	} else {
		setNamespaceIfScoped(namespaceOverride, u)
	}
}

func setDefaultNamespaceIfScopedAndNoneSet(u *unstructured.Unstructured) {
	namespace := u.GetNamespace()
	if isNamespaceScoped(u.GetKind()) && namespace == "" {
		namespace = "default"
		u.SetNamespace(namespace)
	}
}

func setNamespaceIfScoped(namespace string, u *unstructured.Unstructured) {
	if isNamespaceScoped(u.GetKind()) {
		u.SetNamespace(namespace)
	}
}

// isNamespaceScoped returns true if the kind is namespace-scoped.
//
// Cluster-scoped kinds are:
//   - apiService
//   - clusterRole
//   - clusterRoleBinding
//   - customResourceDefinition
//   - mutatingWebhookConfiguration
//   - namespace
//   - persistentVolume
//   - podSecurityPolicy
//   - storageClass
//   - validatingWebhookConfiguration
//
// See https://github.com/spinnaker/clouddriver/blob/58ab154b0ec0d62772201b5b319af349498a4e3f/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/description/manifest/KubernetesKindProperties.java#L31
// for clouddriver OSS namespace-scoped kinds.
func isNamespaceScoped(kind string) bool {
	namespaceScoped := true

	for _, value := range clusterScopedKinds {
		if strings.EqualFold(value, kind) {
			namespaceScoped = false

			break
		}
	}

	return namespaceScoped
}
