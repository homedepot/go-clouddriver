package kubernetes

import "time"

type Resource struct {
	AccountName  string    `json:"accountName" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:1;index:idx_kubernetes_resources_kind_covering,priority:2"`
	ID           string    `json:"id" gorm:"primary_key"`
	Timestamp    time.Time `json:"timestamp,omitempty" gorm:"type:timestamp;DEFAULT:current_timestamp"`
	TaskID       string    `json:"taskId" gorm:"index:task_id_idx"`
	TaskType     string    `json:"-"`
	APIGroup     string    `json:"apiGroup"`
	Name         string    `json:"name" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:3;index:idx_kubernetes_resources_kind_covering,priority:3"`
	ArtifactName string    `json:"-"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Version      string    `json:"version"`
	Kind         string    `json:"kind" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:2;index:kind_idx;index:idx_kubernetes_resources_kind_covering,priority:1"`
	SpinnakerApp string    `json:"spinnakerApp" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:4;index:idx_kubernetes_resources_kind_covering,priority:4"`
	Cluster      string    `json:"-"`
}

// Composite index declared above (currently created automatically by
// AutoMigrate in sql.Client.Connect on every startup):
//
//   idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app)
//     Covers ListKubernetesClustersByFields / ListKubernetesClustersByApplication,
//     which filter on kind and select/group on all four columns.
//
// This index is slated to move to manual DDL management (removing it from
// this struct's gorm tags entirely) to eliminate the multi-replica
// AutoMigrate race on a write-heavy production table - see
// docs/decisions/kind-covering-index-rollout-runbook.md. Do not remove the
// tag until that runbook's rollout steps have been executed against the
// database; until then, AutoMigrate is still the source of truth.
//
// A second covering index for the spinnaker_app-scoped queries
// (ListKubernetesClustersByApplication, ListKubernetesAccountsBySpinnakerApp)
// was deliberately deferred - add it only if slow-query logs show it's
// needed, to avoid paying write-amplification cost on this write-heavy
// table for an index that may not be justified.

func (Resource) TableName() string {
	return "kubernetes_resources"
}
