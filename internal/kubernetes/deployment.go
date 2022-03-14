package kubernetes

import (
	"encoding/json"
	"strings"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

func NewDeployment(m map[string]interface{}) *Deployment {
	d := &v1.Deployment{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &d)

	return &Deployment{d: d}
}

type Deployment struct {
	d *v1.Deployment
}

func (d *Deployment) Object() *v1.Deployment {
	return d.d
}

func (d *Deployment) Status() manifest.Status {
	s := manifest.DefaultStatus

	if d.d.ObjectMeta.Generation != d.d.Status.ObservedGeneration {
		s.Stable.State = false
		s.Stable.Message = "Waiting for status generation to match updated object generation"

		return s
	}

	conditions := d.d.Status.Conditions
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
			s.Failed.Message = condition.Message
		}
	}

	desiredReplicas := int32(0)

	if d.d.Spec.Replicas != nil {
		desiredReplicas = *d.d.Spec.Replicas
	}

	{
		updatedReplicas := d.d.Status.UpdatedReplicas
		if updatedReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be updated"

			return s
		}

		statusReplicas := d.d.Status.Replicas
		if statusReplicas > updatedReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for old replicas to finish termination"

			return s
		}
	}

	{
		availableReplicas := d.d.Status.AvailableReplicas
		if availableReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be available"

			return s
		}
	}

	{
		readyReplicas := d.d.Status.ReadyReplicas
		if readyReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be ready"

			return s
		}
	}

	return s
}
