package kubernetes

import (
	"encoding/json"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func NewDeployment(m map[string]interface{}) Deployment {
	d := &v1.Deployment{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &d)

	return &deployment{d: d}
}

type Deployment interface {
	ToUnstructured() (unstructured.Unstructured, error)
	AnnotateTemplate(string, string)
	GetSpec() v1.DeploymentSpec
	SetReplicas(*int32)
	LabelTemplate(string, string)
	Status() manifest.Status
	Object() *v1.Deployment
}

type deployment struct {
	d *v1.Deployment
}

func (d *deployment) ToUnstructured() (unstructured.Unstructured, error) {
	u := unstructured.Unstructured{}

	b, err := json.Marshal(d.d)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(b, &u.Object)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (d *deployment) Object() *v1.Deployment {
	return d.d
}

func (d *deployment) AnnotateTemplate(key, value string) {
	annotations := d.d.Spec.Template.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	d.d.Spec.Template.ObjectMeta.Annotations = annotations
}

func (d *deployment) GetSpec() v1.DeploymentSpec {
	return d.d.Spec
}

func (d *deployment) SetReplicas(replicas *int32) {
	d.d.Spec.Replicas = replicas
}

func (d *deployment) LabelTemplate(key, value string) {
	labels := d.d.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	d.d.Spec.Template.ObjectMeta.Labels = labels
}

func (d *deployment) Status() manifest.Status {
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
