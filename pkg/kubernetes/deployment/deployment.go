package deployment

import (
	"encoding/json"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

func New(m map[string]interface{}) *v1.Deployment {
	p := &v1.Deployment{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &p)

	return p
}

func Status(m map[string]interface{}) manifest.Status {
	s := manifest.DefaultStatus
	d := New(m)

	if d.ObjectMeta.Generation != d.Status.ObservedGeneration {
		s.Stable.State = false
		s.Stable.Message = "Waiting for status generation to match updated object generation"

		return s
	}

	conditions := d.Status.Conditions
	for _, condition := range conditions {
		if strings.EqualFold(condition.Reason, "deploymentpaused") {
			s.Paused.State = true
		}

		if strings.EqualFold(string(condition.Type), "available") &&
			strings.EqualFold(string(condition.Status), "false") {
			s.Available.State = false
			s.Available.Message = condition.Reason
			s.Stable.State = false
			s.Stable.Message = condition.Reason
		}

		if strings.EqualFold(string(condition.Type), "progressing") &&
			strings.EqualFold(condition.Reason, "progressdeadlineexceeded") {
			s.Failed.State = true
		}
	}

	desiredReplicas := int32(0)

	if d.Spec.Replicas != nil {
		desiredReplicas = *d.Spec.Replicas
	}

	{
		updatedReplicas := d.Status.UpdatedReplicas
		if updatedReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be updated"

			return s
		}

		statusReplicas := d.Status.Replicas
		if statusReplicas > updatedReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for old replicas to finish termination"

			return s
		}
	}

	{
		availableReplicas := d.Status.AvailableReplicas
		if availableReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be available"

			return s
		}
	}

	{
		readyReplicas := d.Status.ReadyReplicas
		if readyReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be ready"

			return s
		}
	}

	return s
}
