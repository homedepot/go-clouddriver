package replicaset

import (
	"encoding/json"
	"log"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

func Status(m map[string]interface{}) manifest.Status {
	s := manifest.DefaultStatus

	r := &v1.ReplicaSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &r)

	desired := r.Spec.Replicas
	fullyLabeled := r.Status.FullyLabeledReplicas
	available := r.Status.AvailableReplicas
	ready := r.Status.ReadyReplicas

	if desired == nil {
		*desired = 0
	}

	if *desired > fullyLabeled {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be fully-labeled"

		return s
	}

	if *desired > ready {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be ready"

		return s
	}

	if *desired > available {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be available"

		return s
	}

	if r.ObjectMeta.Generation != r.Status.ObservedGeneration {
		log.Println("r.ObjectMeta.Generation", r.ObjectMeta.Generation)
		log.Println("r.Status.ObservedGeneration", r.Status.ObservedGeneration)

		s.Stable.State = false
		s.Stable.Message = "Waiting for replicaset spec update to be observed"

		return s
	}

	return s
}
