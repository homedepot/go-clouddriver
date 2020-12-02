package kubernetes

import (
	"encoding/json"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ReplicaSet interface {
	ToUnstructured() (unstructured.Unstructured, error)
	AnnotateTemplate(string, string)
	GetReplicaSetSpec() v1.ReplicaSetSpec
	GetReplicaSetStatus() v1.ReplicaSetStatus
	LabelTemplate(string, string)
	LabelTemplateIfNotExists(string, string)
	Status() manifest.Status
	SetReplicas(*int32)
	ListImages() []string
	Object() *v1.ReplicaSet
}

func NewReplicaSet(m map[string]interface{}) ReplicaSet {
	r := &v1.ReplicaSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &r)

	return &replicaSet{rs: r}
}

type replicaSet struct {
	rs *v1.ReplicaSet
}

func (rs *replicaSet) Object() *v1.ReplicaSet {
	return rs.rs
}

func (rs *replicaSet) ToUnstructured() (unstructured.Unstructured, error) {
	u := unstructured.Unstructured{}

	b, err := json.Marshal(rs.rs)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(b, &u.Object)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (rs *replicaSet) AnnotateTemplate(key, value string) {
	annotations := rs.rs.Spec.Template.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	rs.rs.Spec.Template.ObjectMeta.Annotations = annotations
}

func (rs *replicaSet) LabelTemplate(key, value string) {
	labels := rs.rs.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	rs.rs.Spec.Template.ObjectMeta.Labels = labels
}

func (rs *replicaSet) LabelTemplateIfNotExists(key, value string) {
	labels := rs.rs.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	if _, ok := labels[key]; !ok {
		labels[key] = value
	}

	rs.rs.Spec.Template.ObjectMeta.Labels = labels
}

func (rs *replicaSet) SetReplicas(replicas *int32) {
	rs.rs.Spec.Replicas = replicas
}

func (rs *replicaSet) GetReplicaSetSpec() v1.ReplicaSetSpec {
	return rs.rs.Spec
}

func (rs *replicaSet) GetReplicaSetStatus() v1.ReplicaSetStatus {
	return rs.rs.Status
}

func (rs *replicaSet) ListImages() []string {
	images := []string{}
	for _, container := range rs.rs.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}

	return images
}

func (rs *replicaSet) Status() manifest.Status {
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
