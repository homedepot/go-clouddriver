package kubernetes

type Operations []struct {
	DeployManifest *DeployManifestRequest `json:"deployManifest"`
	ScaleManifest  *ScaleManifestRequest  `json:"scaleManifest"`
}
