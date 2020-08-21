package kubernetes

import "github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"

type OperationsResponse struct {
	ID          string `json:"id"`
	ResourceURI string `json:"resourceUri"`
}

type Operations []struct {
	DeployManifest         *DeployManifestRequest         `json:"deployManifest"`
	ScaleManifest          *ScaleManifestRequest          `json:"scaleManifest"`
	CleanupArtifacts       *CleanupArtifactsRequest       `json:"cleanupArtifacts"`
	DeleteManifest         *DeleteManifestRequest         `json:"deleteManifest"`
	UndoRolloutManifest    *UndoRolloutManifestRequest    `json:"undoRolloutManifest"`
	RollingRestartManifest *RollingRestartManifestRequest `json:"rollingRestartManifest"`
}

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

type ScaleManifestRequest struct {
	Replicas      string `json:"replicas"`
	ManifestName  string `json:"manifestName"`
	CloudProvider string `json:"cloudProvider"`
	Location      string `json:"location"`
	User          string `json:"user"`
	Account       string `json:"account"`
}

type CleanupArtifactsRequest struct {
	Manifests []map[string]interface{} `json:"manifests"`
	Account   string                   `json:"account"`
}

type DeleteManifestRequest struct {
	ManifestName  string `json:"manifestName"`
	CloudProvider string `json:"cloudProvider"`
	Options       struct {
		OrphanDependants   bool `json:"orphanDependants"`
		GracePeriodSeconds int  `json:"gracePeriodSeconds"`
	} `json:"options"`
	Location string `json:"location"`
	User     string `json:"user"`
	Account  string `json:"account"`
}

type UndoRolloutManifestRequest struct {
	ManifestName  string `json:"manifestName"`
	CloudProvider string `json:"cloudProvider"`
	Location      string `json:"location"`
	User          string `json:"user"`
	Account       string `json:"account"`
	Revision      string `json:"revision"`
}

type RollingRestartManifestRequest struct {
	CloudProvider string `json:"cloudProvider"`
	ManifestName  string `json:"manifestName"`
	Location      string `json:"location"`
	User          string `json:"user"`
	Account       string `json:"account"`
}
