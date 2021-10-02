package kubernetes

import (
	"encoding/json"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

func NewReplicaSet(m map[string]interface{}) *ReplicaSet {
	r := &v1.ReplicaSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &r)

	return &ReplicaSet{rs: r}
}

type ReplicaSet struct {
	rs *v1.ReplicaSet
}

func (rs *ReplicaSet) Object() *v1.ReplicaSet {
	return rs.rs
}

func (rs *ReplicaSet) Status() manifest.Status {
	s := manifest.DefaultStatus
	r := rs.rs

	desired := int32(0)
	fullyLabeled := r.Status.FullyLabeledReplicas
	available := r.Status.AvailableReplicas
	ready := r.Status.ReadyReplicas

	if r.Spec.Replicas != nil {
		desired = *r.Spec.Replicas
	}

	if desired > fullyLabeled {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be fully-labeled"

		return s
	}

	if desired > ready {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be ready"

		return s
	}

	if desired > available {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be available"

		return s
	}

	if r.ObjectMeta.Generation != r.Status.ObservedGeneration {
		s.Stable.State = false
		s.Stable.Message = "Waiting for replicaset spec update to be observed"

		return s
	}

	return s
}
