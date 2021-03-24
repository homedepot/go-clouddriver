package kubernetes

import (
	"encoding/json"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Pod interface {
	Status() manifest.Status
	GetObjectMeta() metav1.ObjectMeta
	GetPodStatus() v1.PodStatus
	GetLabels() map[string]string
	GetNamespace() string
	GetName() string
	GetUID() string
	Object() *v1.Pod
	GetSpec() v1.PodSpec
	ToUnstructured() (unstructured.Unstructured, error)
}

func NewPod(m map[string]interface{}) Pod {
	p := &v1.Pod{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &p)

	return &pod{p: p}
}

type pod struct {
	p *v1.Pod
}

func (p *pod) Object() *v1.Pod {
	return p.p
}

func (p *pod) GetObjectMeta() metav1.ObjectMeta {
	return p.p.ObjectMeta
}

func (p *pod) GetPodStatus() v1.PodStatus {
	return p.p.Status
}

func (p *pod) GetNamespace() string {
	return p.p.ObjectMeta.Namespace
}

func (p *pod) GetName() string {
	return p.p.ObjectMeta.Name
}

func (p *pod) GetUID() string {
	return string(p.p.ObjectMeta.UID)
}

func (p *pod) GetLabels() map[string]string {
	return p.p.ObjectMeta.Labels
}

func (p *pod) Status() manifest.Status {
	s := manifest.DefaultStatus

	if p.p.Status.Phase == v1.PodPending ||
		p.p.Status.Phase == v1.PodFailed ||
		p.p.Status.Phase == v1.PodUnknown {
		s.Stable.State = false
		s.Stable.Message = "Pod is " + string(p.p.Status.Phase)
		s.Available.State = false
		s.Available.Message = "Pod is " + string(p.p.Status.Phase)
	}

	return s
}

func (p *pod) GetSpec() v1.PodSpec {
	return p.p.Spec
}

func (p *pod) ToUnstructured() (unstructured.Unstructured, error) {
	u := unstructured.Unstructured{}

	b, err := json.Marshal(p.p)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(b, &u.Object)
	if err != nil {
		return u, err
	}

	return u, nil
}
