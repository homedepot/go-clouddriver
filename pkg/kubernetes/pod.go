package kubernetes

import (
	"encoding/json"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/core/v1"
)

func NewPod(m map[string]interface{}) *Pod {
	p := &v1.Pod{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &p)

	return &Pod{p: p}
}

type Pod struct {
	p *v1.Pod
}

func (p *Pod) Object() *v1.Pod {
	return p.p
}

func (p *Pod) Status() manifest.Status {
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
