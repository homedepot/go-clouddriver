package kubernetes

import (
	"fmt"
	"strings"
)

var (
	clusterScopedKinds = []string{
		"apiService",
		"clusterRole",
		"clusterRoleBinding",
		"customResourceDefinition",
		"mutatingWebhookConfiguration",
		"namespace",
		"persistentVolume",
		"podSecurityPolicy",
		"storageClass",
		"validatingWebhookConfiguration",
	}
)

type Provider struct {
	Name          string              `json:"name" gorm:"primary_key"`
	Host          string              `json:"host"`
	CAData        string              `json:"caData" gorm:"type:text"`
	BearerToken   string              `json:"bearerToken,omitempty" gorm:"size:2048"`
	TokenProvider string              `json:"tokenProvider,omitempty" gorm:"size:32;not null;default:'google'"`
	Namespace     *string             `json:"namespace,omitempty" gorm:"size:253"`
	Namespaces    []string            `json:"namespaces,omitempty" gorm:"-"`
	Permissions   ProviderPermissions `json:"permissions" gorm:"-"`
	// Providers can hold instances of clients.
	Client    Client    `json:"-" gorm:"-"`
	Clientset Clientset `json:"-" gorm:"-"`
}

type ProviderPermissions struct {
	Read  []string `json:"read" gorm:"-"`
	Write []string `json:"write" gorm:"-"`
}

func (Provider) TableName() string {
	return "kubernetes_providers"
}

type ProviderNamespaces struct {
	//ID          string `json:"-" gorm:"primary_key"`
	AccountName string `json:"accountName"`
	Namespace   string `json:"namespace,omitempty"`
}

func (ProviderNamespaces) TableName() string {
	return "kubernetes_providers_namespaces"
}

// ValidateKindStatus verifies that this provider can access the given kind.
// This begins to support `omitKinds`, but only in the context of namespace-scoped
// providers.
//
// When a provider is limited to namespace, then it cannot access these kinds:
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
func (p *Provider) ValidateKindStatus(kind string) error {
	if p.Namespace == nil && len(p.Namespaces) == 0 {
		return nil
	}

	for _, value := range clusterScopedKinds {
		if strings.EqualFold(value, kind) {
			return fmt.Errorf("namespace-scoped account not allowed to access cluster-scoped kind: '%s'", kind)
		}
	}

	return nil
}

// WithClient sets the kubernetes client for this provider.
func (p *Provider) WithClient(client Client) {
	p.Client = client
}

// WithClientset sets the kubernetes clientset for this provider.
func (p *Provider) WithClientset(clientset Clientset) {
	p.Clientset = clientset
}
