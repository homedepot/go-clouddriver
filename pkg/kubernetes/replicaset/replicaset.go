package replicaset

import (
	"encoding/json"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

func New(m map[string]interface{}) *v1.ReplicaSet {
	r := &v1.ReplicaSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &r)

	return r
}

func Status(m map[string]interface{}) manifest.Status {
	s := manifest.DefaultStatus

	r := New(m)

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
