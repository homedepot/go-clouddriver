package kubernetes

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
)

func NewDaemonSet(m map[string]interface{}) *DaemonSet {
	ds := &v1.DaemonSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &ds)

	return &DaemonSet{ds: ds}
}

type DaemonSet struct {
	ds *v1.DaemonSet
}

func (ds *DaemonSet) Object() *v1.DaemonSet {
	return ds.ds
}

func (ds *DaemonSet) Status() manifest.Status {
	s := manifest.DefaultStatus

	if reflect.DeepEqual(ds.ds.Status, v1.DaemonSetStatus{}) {
		s = manifest.NoneReported

		return s
	}

	if strings.EqualFold(string(ds.ds.Spec.UpdateStrategy.Type), "rollingupdate") {
		return s
	}

	if ds.ds.ObjectMeta.Generation != ds.ds.Status.ObservedGeneration {
		s.Stable.State = false
		s.Stable.Message = "Waiting for status generation to match updated object generation"

		return s
	}

	desiredReplicas := ds.ds.Status.DesiredNumberScheduled

	{
		scheduledReplicas := ds.ds.Status.CurrentNumberScheduled
		if scheduledReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be scheduled"

			return s
		}
	}

	{
		updatedReplicas := ds.ds.Status.UpdatedNumberScheduled
		if updatedReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all updated replicas to be scheduled"

			return s
		}
	}

	{
		availableReplicas := ds.ds.Status.NumberAvailable
		if availableReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be available"

			return s
		}
	}

	{
		readyReplicas := ds.ds.Status.NumberReady
		if readyReplicas < desiredReplicas {
			s.Stable.State = false
			s.Stable.Message = "Waiting for all replicas to be ready"

			return s
		}
	}

	return s
}
