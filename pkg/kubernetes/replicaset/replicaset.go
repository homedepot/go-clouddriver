package replicaset

import (
	"encoding/json"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ReplicaSet struct {
	rs *v1.ReplicaSet
}

func New(m map[string]interface{}) ReplicaSet {
	r := &v1.ReplicaSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &r)

	return ReplicaSet{rs: r}
}

func (rs *ReplicaSet) ToUnstructured() (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}

	b, err := json.Marshal(rs.rs)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &u.Object)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (rs *ReplicaSet) AnnotateTemplate(key, value string) {
	annotations := rs.rs.Spec.Template.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	rs.rs.Spec.Template.ObjectMeta.Annotations = annotations
}

func (rs *ReplicaSet) LabelTemplate(key, value string) {
	labels := rs.rs.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	rs.rs.Spec.Template.ObjectMeta.Labels = labels
}

func (rs *ReplicaSet) GetSpec() v1.ReplicaSetSpec {
	return rs.rs.Spec
}

func (rs *ReplicaSet) GetStatus() v1.ReplicaSetStatus {
	return rs.rs.Status
}

func Status(m map[string]interface{}) manifest.Status {
	s := manifest.DefaultStatus

	rs := New(m)
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
