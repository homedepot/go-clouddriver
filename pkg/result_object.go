package clouddriver

type ResultObject struct {
	BoundArtifacts                    []interface{}            `json:"boundArtifacts"`
	CreatedArtifacts                  []CreatedArtifact        `json:"createdArtifacts"`
	ManifestNamesByNamespace          map[string][]string      `json:"manifestNamesByNamespace"`
	ManifestNamesByNamespaceToRefresh map[string][]string      `json:"manifestNamesByNamespaceToRefresh"`
	Manifests                         []map[string]interface{} `json:"manifests"`
}
