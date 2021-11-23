package clouddriver

import "github.com/homedepot/go-clouddriver/internal/artifact"

type Artifact struct {
	CustomKind bool             `json:"customKind"`
	Location   string           `json:"location,omitempty"`
	Metadata   ArtifactMetadata `json:"metadata"`
	Name       string           `json:"name"`
	Reference  string           `json:"reference"`
	Type       artifact.Type    `json:"type"`
	Version    string           `json:"version,omitempty"`
}

type ArtifactMetadata struct {
	Account string `json:"account,omitempty"`
}
