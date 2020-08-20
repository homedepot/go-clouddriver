package deployment

import (
	"encoding/json"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ToUnstructured(d *v1.Deployment) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &u.Object)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func AnnotateTemplate(d *v1.Deployment, key, value string) {
	annotations := d.Spec.Template.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	d.Spec.Template.ObjectMeta.Annotations = annotations
}

func LabelTemplate(d *v1.Deployment, key, value string) {
	labels := d.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	d.Spec.Template.ObjectMeta.Labels = labels
}

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
