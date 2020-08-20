package kubernetes

import (
	"encoding/base64"
	"encoding/json"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8s.io/client-go/rest"
)

type DeployManifestRequest struct {
	EnableTraffic     bool                     `json:"enableTraffic"`
	NamespaceOverride string                   `json:"namespaceOverride"`
	OptionalArtifacts []interface{}            `json:"optionalArtifacts"`
	CloudProvider     string                   `json:"cloudProvider"`
	Manifests         []map[string]interface{} `json:"manifests"`
	TrafficManagement struct {
		Options struct {
			EnableTraffic bool `json:"enableTraffic"`
		} `json:"options"`
		Enabled bool `json:"enabled"`
	} `json:"trafficManagement"`
	Moniker struct {
		App string `json:"app"`
	} `json:"moniker"`
	Source                   string        `json:"source"`
	Account                  string        `json:"account"`
	SkipExpressionEvaluation bool          `json:"skipExpressionEvaluation"`
	RequiredArtifacts        []interface{} `json:"requiredArtifacts"`
}

type OperationsResponse struct {
	ID          string `json:"id"`
	ResourceURI string `json:"resourceUri"`
}

type ManifestResponse struct {
	Account string `json:"account"`
	// Artifacts []struct {
	// 	CustomKind bool `json:"customKind"`
	// 	Metadata   struct {
	// 	} `json:"metadata"`
	// 	Name      string `json:"name"`
	// 	Reference string `json:"reference"`
	// 	Type      string `json:"type"`
	// } `json:"artifacts"`
	Events   []interface{}           `json:"events"`
	Location string                  `json:"location"`
	Manifest map[string]interface{}  `json:"manifest"`
	Metrics  []interface{}           `json:"metrics"`
	Moniker  ManifestResponseMoniker `json:"moniker"`
	Name     string                  `json:"name"`
	Status   manifest.Status         `json:"status"`
	Warnings []interface{}           `json:"warnings"`
}

type ManifestResponseMoniker struct {
	App     string `json:"app"`
	Cluster string `json:"cluster"`
}

func DeployManifests(c *gin.Context, taskID string, dm DeployManifestRequest) error {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)

	provider, err := sc.GetKubernetesProvider(dm.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: provider.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	kc.WithConfig(config)

	for _, manifest := range dm.Manifests {
		b, err := json.Marshal(manifest)
		if err != nil {
			return err
		}

		meta, err := kc.Apply(b, dm.Moniker.App)
		if err != nil {
			return err
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			APIGroup:     meta.Group,
			Name:         meta.Name,
			Namespace:    meta.Namespace,
			Resource:     meta.Resource,
			Version:      meta.Version,
			Kind:         meta.Kind,
			SpinnakerApp: dm.Moniker.App,
		}

		err = sc.CreateKubernetesResource(kr)
		if err != nil {
			return err
		}
	}

	return nil
}
