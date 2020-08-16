package clouddriver

type DeployManifests []DeployManifest

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
