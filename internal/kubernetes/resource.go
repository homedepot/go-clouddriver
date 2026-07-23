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

// Composite index declared above (created automatically by AutoMigrate in
// sql.Client.Connect on every startup - this is intentional and should stay
// this way, so any new environment/install of this project gets this index
// with no manual step required):
//
//   idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app)
//     Covers ListKubernetesClustersByFields / ListKubernetesClustersByApplication,
//     which filter on kind and select/group on all four columns.
//
// If you operate a deployment of this project with more than one replica
// AND a kubernetes_resources table large enough that a live ALTER TABLE is
// an operational concern for your environment, consider manually
// pre-creating this index via a directly-run ALTER TABLE ... ADD INDEX ...,
// ALGORITHM=INPLACE, LOCK=NONE statement before rolling out a change that
// introduces it, so every replica's AutoMigrate call finds it already
// present (via HasIndex) and skips CreateIndex - this avoids multiple
// replicas racing to create the same index concurrently on a rolling
// deploy. This is a deployment-time operational choice, not a code change:
// the gorm tag stays as-is either way.
//
// A second covering index for the spinnaker_app-scoped queries
// (ListKubernetesClustersByApplication, ListKubernetesAccountsBySpinnakerApp)
// was deliberately deferred - add it only if slow-query logs show it's
// needed, to avoid paying write-amplification cost on this write-heavy
// table for an index that may not be justified.

func (Resource) TableName() string {
	return "kubernetes_resources"
}
