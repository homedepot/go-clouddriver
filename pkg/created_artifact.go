package clouddriver

type CreatedArtifact struct {
	CustomKind bool   `json:"customKind"`
	Location   string `json:"location"`
	Metadata   struct {
		Account string `json:"account"`
	} `json:"metadata"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Type      string `json:"type"`
	Version   string `json:"version"`
}
