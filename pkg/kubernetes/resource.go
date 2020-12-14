package kubernetes

import "time"

type Resource struct {
	AccountName  string    `json:"accountName"`
	ID           string    `json:"id" gorm:"primary_key"`
	Timestamp    time.Time `json:"timestamp,omitempty" gorm:"type:timestamp;DEFAULT:current_timestamp"`
	TaskID       string    `json:"taskId"`
	TaskType     string    `json:"-"`
	APIGroup     string    `json:"apiGroup"`
	Name         string    `json:"name"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Version      string    `json:"version"`
	Kind         string    `json:"kind"`
	SpinnakerApp string    `json:"spinnakerApp"`
	Cluster      string    `json:"-"`
}

func (Resource) TableName() string {
	return "kubernetes_resources"
}
