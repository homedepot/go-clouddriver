package kubernetes

import (
	"encoding/json"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/autoscaling/v1"
	"fmt"
)

type HorizontalPodAutoscaler interface {
    Status() manifest.Status
	Object() *v1.HorizontalPodAutoscaler	
}

func NewHorizontalPodAutoscaler(m map[string]interface{}) HorizontalPodAutoscaler {
	hpa := &v1.HorizontalPodAutoscaler{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &hpa)

	return &horizontalPodAutoscaler{hpa: hpa}
}

type horizontalPodAutoscaler struct{
	hpa *v1.HorizontalPodAutoscaler
}

func (hpa *horizontalPodAutoscaler) Object() *v1.HorizontalPodAutoscaler {
	return hpa.hpa
}

func (hpa *horizontalPodAutoscaler) Status() manifest.Status {
	s := manifest.DefaultStatus

	hpaStatus := hpa.hpa.Status
	if(hpaStatus.DesiredReplicas > hpaStatus.CurrentReplicas){
		s.Stable.State = false
		s.Stable.Message = fmt.Sprintf("Waiting for HPA to complete a scale up, current: %d desired: %d",hpaStatus.CurrentReplicas, hpaStatus.DesiredReplicas)
		s.Available.State = false
		s.Available.Message = fmt.Sprintf("Waiting for HPA to complete a scale up, current: %d desired: %d",hpaStatus.CurrentReplicas, hpaStatus.DesiredReplicas)
	}
	if(hpaStatus.DesiredReplicas < hpaStatus.CurrentReplicas) {
		s.Stable.State = false
		s.Stable.Message = fmt.Sprintf("Waiting for HPA to complete a scale down, current: %d desired: %d",hpaStatus.CurrentReplicas, hpaStatus.DesiredReplicas)
		s.Available.State = false
		s.Available.Message = fmt.Sprintf("Waiting for HPA to complete a scale down, current: %d desired: %d",hpaStatus.CurrentReplicas, hpaStatus.DesiredReplicas)
	}

	return s
}