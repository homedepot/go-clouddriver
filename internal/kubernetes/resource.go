package kubernetes

import "time"

type Resource struct {
	AccountName  string    `json:"accountName" gorm:"index:kind_account_name_kind_name_spinnaker_app_idx,sort:asc;index:account_name_idx,sort:asc"`
	ID           string    `json:"id" gorm:"primary_key"`
	Timestamp    time.Time `json:"timestamp,omitempty" gorm:"type:timestamp;DEFAULT:current_timestamp"`
	TaskID       string    `json:"taskId" gorm:"index:task_id_idx,sort:asc"`
	TaskType     string    `json:"-"`
	APIGroup     string    `json:"apiGroup"`
	Name         string    `json:"name" gorm:"index:kind_account_name_kind_name_spinnaker_app_idx,sort:asc"`
	ArtifactName string    `json:"-"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Version      string    `json:"version"`
	Kind         string    `json:"kind" gorm:"index:kind_account_name_kind_name_spinnaker_app_idx,sort:asc;index:kind_idx,sort:asc"`
	SpinnakerApp string    `json:"spinnakerApp" gorm:"index:kind_account_name_kind_name_spinnaker_app_idx,sort:asc"`
	Cluster      string    `json:"-"`
}

func (Resource) TableName() string {
	return "kubernetes_resources"
}
