package kubernetes

import (
	"encoding/json"

	v1 "k8s.io/api/core/v1"
)

type Service interface {
	Selector() map[string]string
}

type service struct {
	s *v1.Service
}

func NewService(m map[string]interface{}) Service {
	s := &v1.Service{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &s)

	return &service{s: s}
}

func (s *service) Selector() map[string]string {
	return s.s.Spec.Selector
}
