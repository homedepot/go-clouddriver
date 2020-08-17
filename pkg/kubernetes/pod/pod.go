package pod

import (
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/mitchellh/mapstructure"
)

type pod struct {
	status status
}

type status struct {
	Phase string `json:"phase"`
}

func Status(m map[string]interface{}) kubernetes.ManifestStatus {
	s := kubernetes.DefaultStatus

	p := &pod{}
	if err := mapstructure.Decode(m, &p); err != nil {
		return kubernetes.NoneReported
	}

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
