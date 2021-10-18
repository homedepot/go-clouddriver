package kubernetes

import (
	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

type OperationsResponse struct {
	ID          string `json:"id"`
	ResourceURI string `json:"resourceUri"`
}

type Operations []Operation

type Operation struct {
	CleanupArtifacts       *CleanupArtifactsRequest       `json:"cleanupArtifacts"`
	DeleteManifest         *DeleteManifestRequest         `json:"deleteManifest"`
	DeployManifest         *DeployManifestRequest         `json:"deployManifest"`
	DisableManifest        *DisableManifestRequest        `json:"disableManifest"`
	PatchManifest          *PatchManifestRequest          `json:"patchManifest"`
	RollingRestartManifest *RollingRestartManifestRequest `json:"rollingRestartManifest"`
	RunJob                 *RunJobRequest                 `json:"runJob"`
	ScaleManifest          *ScaleManifestRequest          `json:"scaleManifest"`
	UndoRolloutManifest    *UndoRolloutManifestRequest    `json:"undoRolloutManifest"`
}

type DeployManifestRequest struct {
	EnableTraffic     bool                     `json:"enableTraffic"`
	NamespaceOverride string                   `json:"namespaceOverride"`
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
	Source                   string                 `json:"source"`
	Account                  string                 `json:"account"`
	SkipExpressionEvaluation bool                   `json:"skipExpressionEvaluation"`
	RequiredArtifacts        []clouddriver.Artifact `json:"requiredArtifacts"`
	OptionalArtifacts        []clouddriver.Artifact `json:"optionalArtifacts"`
}

type DisableManifestRequest struct {
	App           string `json:"app"`
	CloudProvider string `json:"cloudProvider"`
	ManifestName  string `json:"manifestName"`
	Location      string `json:"location"`
	Account       string `json:"account"`
}

type PatchManifestRequest struct {
	App      string `json:"app"`
	Cluster  string `json:"cluster"`
	Criteria string `json:"criteria"`
	// Kind          string                         `json:"kind"`
	ManifestName  string                         `json:"manifestName"`
	Source        string                         `json:"source"`
	Mode          string                         `json:"mode"`
	PatchBody     map[string]interface{}         `json:"patchBody"`
	CloudProvider string                         `json:"cloudProvider"`
	AllArtifacts  []PatchManifestRequestArtifact `json:"allArtifacts"`
	Options       PatchManifestRequestOptions    `json:"options"`
	// Manifests         []map[string]interface{}       `json:"manifests"`
	Location string `json:"location"`
	Account  string `json:"account"`
	// RequiredArtifacts []interface{}                  `json:"requiredArtifacts"`
}

type PatchManifestRequestArtifact struct {
	CustomKind bool   `json:"customKind"`
	Reference  string `json:"reference"`
	Metadata   struct {
		Account string `json:"account"`
	} `json:"metadata"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Type     string `json:"type"`
	Version  string `json:"version"`
}

// Merge strategy can be "strategic", "json", or "merge".
type PatchManifestRequestOptions struct {
	MergeStrategy string `json:"mergeStrategy"`
	Record        bool   `json:"record"`
}

// why are artifacts commented out here? possibly causing the problem of artifacts not getting bound correctly
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

type ManifestCoordinatesResponse struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
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
	App            string                              `json:"app"`
	Mode           string                              `json:"mode"`
	LabelSelectors DeleteManifestRequestLabelSelectors `json:"labelSelectors"`
	ManifestName   string                              `json:"manifestName"`
	CloudProvider  string                              `json:"cloudProvider"`
	Options        DeleteManifestRequestOptions        `json:"options"`
	Kinds          []string                            `json:"kinds"`
	Location       string                              `json:"location"`
	User           string                              `json:"user"`
	Account        string                              `json:"account"`
}

type DeleteManifestRequestLabelSelectors struct {
	Selectors []DeleteManifestRequestLabelSelector `json:"selectors"`
}

type DeleteManifestRequestLabelSelector struct {
	Kind   string   `json:"kind"`
	Values []string `json:"values"`
	Key    string   `json:"key"`
}

type DeleteManifestRequestOptions struct {
	Cascading          bool   `json:"cascading"`
	OrphanDependants   *bool  `json:"orphanDependants"`
	GracePeriodSeconds *int64 `json:"gracePeriodSeconds"`
}

type UndoRolloutManifestRequest struct {
	Mode             string `json:"mode"`
	ManifestName     string `json:"manifestName"`
	CloudProvider    string `json:"cloudProvider"`
	Location         string `json:"location"`
	NumRevisionsBack int    `json:"numRevisionsBack"`
	User             string `json:"user"`
	Account          string `json:"account"`
	Revision         string `json:"revision"`
}

type RollingRestartManifestRequest struct {
	CloudProvider string `json:"cloudProvider"`
	ManifestName  string `json:"manifestName"`
	Location      string `json:"location"`
	User          string `json:"user"`
	Account       string `json:"account"`
}

type RunJobRequest struct {
	Account       string                 `json:"account"`
	Alias         string                 `json:"alias"`
	Application   string                 `json:"application"`
	CloudProvider string                 `json:"cloudProvider"`
	Manifest      map[string]interface{} `json:"manifest"`
	// OptionalArtifacts []struct {
	// 	Type       string `json:"type"`
	// 	CustomKind bool   `json:"customKind"`
	// 	Name       string `json:"name"`
	// 	Version    string `json:"version"`
	// 	Location   string `json:"location"`
	// 	Reference  string `json:"reference"`
	// 	Metadata   struct {
	// 		Account string `json:"account"`
	// 	} `json:"metadata"`
	// } `json:"optionalArtifacts"`
	// PreconfiguredJobParameters []struct {
	// 	Mapping     string `json:"mapping"`
	// 	Name        string `json:"name"`
	// 	Description string `json:"description"`
	// 	Label       string `json:"label"`
	// 	Type        string `json:"type"`
	// 	Order       int    `json:"order"`
	// } `json:"preconfiguredJobParameters"`
	// WaitForCompletion bool   `json:"waitForCompletion"`
	// Source            string `json:"source"`
	// Parameters        struct {
	// 	IMAGEPATHS     string `json:"IMAGE_PATHS"`
	// 	SOURCEREGISTRY string `json:"SOURCE_REGISTRY"`
	// 	TARGETREGISTRY string `json:"TARGET_REGISTRY"`
	// } `json:"parameters"`
	// RequiredArtifacts []interface{} `json:"requiredArtifacts"`
}
