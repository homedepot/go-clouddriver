package kubernetes

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

type StatefulSet interface {
	Object() *v1.StatefulSet
	SetReplicas(*int32)
	Status() manifest.Status
}

type statefulSet struct {
	ss *v1.StatefulSet
}

func NewStatefulSet(m map[string]interface{}) StatefulSet {
	s := &v1.StatefulSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &s)

	return &statefulSet{ss: s}
}

func (ss *statefulSet) Object() *v1.StatefulSet {
	return ss.ss
}

func (ss *statefulSet) SetReplicas(replicas *int32) {
	ss.ss.Spec.Replicas = replicas
}

func (ss *statefulSet) Status() manifest.Status {
	s := manifest.DefaultStatus
	x := ss.ss

	if strings.EqualFold(string(x.Spec.UpdateStrategy.Type), "ondelete") {
		return s
	}

	if reflect.DeepEqual(x.Status, v1.StatefulSetStatus{}) {
		s = manifest.NoneReported
		return s
	}

	if x.ObjectMeta.Generation != x.Status.ObservedGeneration {
		s.Stable.State = false
		s.Stable.Message = "Waiting for status generation to match updated object generation"
		return s
	}

	desired := int32(0)
	if x.Spec.Replicas != nil {
		desired = *x.Spec.Replicas
	}

	existing := x.Status.Replicas
	if desired > existing {
		s.Stable.State = false
		s.Stable.Message = "Waiting for at least the desired replica count to be met"

		return s
	}

	ready := x.Status.ReadyReplicas
	if desired > ready {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all updated replicas to be ready"

		return s
	}

	updType := string(x.Spec.UpdateStrategy.Type)
	rollUpd := x.Spec.UpdateStrategy.RollingUpdate
	updated := x.Status.UpdatedReplicas

	if strings.EqualFold(updType, "rollingupdate") && rollUpd != nil {
		partition := rollUpd.Partition
		if partition != nil && (updated < (existing - *partition)) {
			s.Stable.State = false
			s.Stable.Message = "Waiting for partitioned rollout to finish"
			return s
		}
		s.Stable.State = true
		s.Stable.Message = "Partitioned roll out complete"
		return s
	}

	current := x.Status.CurrentReplicas
	if desired > current {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all updated replicas to be scheduled"

		return s
	}

	updateRev := x.Status.UpdateRevision
	currentRev := x.Status.CurrentRevision

	if currentRev != updateRev {
		s.Stable.State = false
		s.Stable.Message = "Waiting for the updated revision to match the current revision"
		return s
	}

	return s
}
