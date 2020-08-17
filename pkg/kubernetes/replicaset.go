package kubernetes

import (
	"strings"

	"github.com/mitchellh/mapstructure"
)

type ReplicaSet struct {
	status Status
}

type ReplicaSetStatus struct {
	Phase string `json:"phase"`
}

func DecodePod(m map[string]interface{}) (*Pod, error) {
	p := &Pod{}
	return p, mapstructure.Decode(m, &p)
}

func (p *Pod) Status() ManifestStatus {
	s := DefaultStatus
	if strings.EqualFold(p.status.Phase, "pending") ||
		strings.EqualFold(p.status.Phase, "failed") ||
		strings.EqualFold(p.status.Phase, "unknown") {
		s.Stable.State = false
		s.Stable.Message = "Pod is " + strings.ToLower(p.status.Phase)
		s.Available.State = false
		s.Available.Message = "Pod is " + strings.ToLower(p.status.Phase)
	}

	return s
}
