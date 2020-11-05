package kubernetes

type ListElement struct {
	APIVersion string                   `json:"apiVersion"`
	Kind       string                   `json:"kind"`
	Items      []map[string]interface{} `json:"items"`
}
