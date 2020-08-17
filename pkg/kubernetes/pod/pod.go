package pod

import (
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	"github.com/mitchellh/mapstructure"
	v1 "k8s.io/api/core/v1"
)

func Status(m map[string]interface{}) manifest.Status {
	s := manifest.DefaultStatus

	p := &v1.Pod{}
	_ = mapstructure.Decode(m, &p)

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
