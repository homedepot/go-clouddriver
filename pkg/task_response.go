package clouddriver

type TaskResponse struct {
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
	ResultObjects []ResultObject `json:"resultObjects"`
	Status        TaskStatus     `json:"status"`
}

type TaskStatus struct {
	Complete  bool   `json:"complete"`
	Completed bool   `json:"completed"`
	Failed    bool   `json:"failed"`
	Phase     string `json:"phase"`
	Retryable bool   `json:"retryable"`
	Status    string `json:"status"`
}
