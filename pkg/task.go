package clouddriver

import "github.com/gin-gonic/gin"

const TaskIDKey = `TaskID`

func NewDefaultTask(id string) Task {
	return Task{
		ID:            id,
		ResultObjects: []TaskResultObject{},
		Status: TaskStatus{
			Complete:  true,
			Completed: true,
			Failed:    false,
			Phase:     "ORCHESTRATION",
			Retryable: false,
			Status:    "Orchestration completed.",
		},
	}
}

func TaskIDFromContext(c *gin.Context) string {
	return c.MustGet(TaskIDKey).(string)
}

type Task struct {
	ID string `json:"id"`
	// SagaIds []interface{} `json:"sagaIds"`
	// History []struct {
	// 	Phase  string `json:"phase"`
	// 	Status string `json:"status"`
	// } `json:"history"`
	// OwnerIDClouddriverSQL   string `json:"ownerId$clouddriver_sql"`
	// RequestIDClouddriverSQL string `json:"requestId$clouddriver_sql"`
	// Retryable                 bool  `json:"retryable"`
	// StartTimeMsClouddriverSQL int64 `json:"startTimeMs$clouddriver_sql"`
	ResultObjects []TaskResultObject `json:"resultObjects"`
	Status        TaskStatus         `json:"status"`
}

type TaskStatus struct {
	Complete  bool   `json:"complete"`
	Completed bool   `json:"completed"`
	Failed    bool   `json:"failed"`
	Phase     string `json:"phase"`
	Retryable bool   `json:"retryable"`
	Status    string `json:"status"`
}

type TaskResultObject struct {
	BoundArtifacts                    []interface{}            `json:"boundArtifacts"`
	CreatedArtifacts                  []TaskCreatedArtifact    `json:"createdArtifacts"`
	DeployedNamesByLocation           map[string][]string      `json:"deployedNamesByLocation"`
	ManifestNamesByNamespace          map[string][]string      `json:"manifestNamesByNamespace"`
	ManifestNamesByNamespaceToRefresh map[string][]string      `json:"manifestNamesByNamespaceToRefresh"`
	Manifests                         []map[string]interface{} `json:"manifests"`
}

type TaskCreatedArtifact struct {
	CustomKind bool                        `json:"customKind"`
	Location   string                      `json:"location"`
	Metadata   TaskCreatedArtifactMetadata `json:"metadata"`
	Name       string                      `json:"name"`
	Reference  string                      `json:"reference"`
	Type       string                      `json:"type"`
	Version    string                      `json:"version"`
}

type TaskCreatedArtifactMetadata struct {
	Account string `json:"account"`
}
