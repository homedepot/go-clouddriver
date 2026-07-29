package kubernetes

import "time"

type Resource struct {
	AccountName  string    `json:"accountName" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:1"`
	ID           string    `json:"id" gorm:"primary_key"`
	Timestamp    time.Time `json:"timestamp,omitempty" gorm:"type:timestamp;DEFAULT:current_timestamp"`
	TaskID       string    `json:"taskId" gorm:"index:task_id_idx"`
	TaskType     string    `json:"-"`
	APIGroup     string    `json:"apiGroup"`
	Name         string    `json:"name" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:3"`
	ArtifactName string    `json:"-"`
	Namespace    string    `json:"namespace"`
	Resource     string    `json:"resource"`
	Version      string    `json:"version"`
	Kind         string    `json:"kind" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:2;index:kind_idx"`
	SpinnakerApp string    `json:"spinnakerApp" gorm:"index:account_name_kind_name_spinnaker_app_idx,priority:4"`
	Cluster      string    `json:"-"`
}

// A covering index (idx_kubernetes_resources_kind_covering, on
// kind/account_name/name/spinnaker_app) was added and then removed during
// this same change. It was built to speed up ListKubernetesClustersByFields
// and ListKubernetesClustersByApplication after their UPPER(kind) predicate
// was made sargable (see internal/sql/client.go's clusterKinds comment), but
// validation against real production-shaped data showed it provided no
// measurable benefit for either query: the query planner never chose it
// naturally, and forcing it made the plan worse (added a temporary table
// for no reduction in rows scanned). The existing
// account_name_kind_name_spinnaker_app_idx (for ListKubernetesClustersByFields)
// and a separately-tracked, out-of-band index on (spinnaker_app, kind) (for
// ListKubernetesClustersByApplication) were already sufficient once the
// UPPER(kind) predicate was removed - that predicate change is what
// delivered the entire measured improvement, not this index. See the
// decision record for the full before/after evidence.

func (Resource) TableName() string {
	return "kubernetes_resources"
}
