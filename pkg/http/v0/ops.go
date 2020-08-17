package v0

import "github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"

type OpsRequest []struct {
	DeployManifest DeployManifest `json:"deployManifest"`
}

type DeployManifest struct {
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

type OpsResponse struct {
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
	Events   []interface{}          `json:"events"`
	Location string                 `json:"location"`
	Manifest map[string]interface{} `json:"manifest"`
	Metrics  []interface{}          `json:"metrics"`
	Moniker  struct {
		App     string `json:"app"`
		Cluster string `json:"cluster"`
	} `json:"moniker"`
	Name     string          `json:"name"`
	Status   manifest.Status `json:"status"`
	Warnings []interface{}   `json:"warnings"`
}
