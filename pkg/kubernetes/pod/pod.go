package pod

import (
	"encoding/json"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/core/v1"
)

func New(m map[string]interface{}) *v1.Pod {
	p := &v1.Pod{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &p)

	return p
}

func Status(m map[string]interface{}) manifest.Status {
	s := manifest.DefaultStatus

	p := New(m)

	if p.Status.Phase == v1.PodPending ||
		p.Status.Phase == v1.PodFailed ||
		p.Status.Phase == v1.PodUnknown {
		s.Stable.State = false
		s.Stable.Message = "Pod is " + string(p.Status.Phase)
		s.Available.State = false
		s.Available.Message = "Pod is " + string(p.Status.Phase)
	}

	return s
}
