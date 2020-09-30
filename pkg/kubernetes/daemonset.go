package kubernetes

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func NewDaemonSet(m map[string]interface{}) DaemonSet {
	ds := &v1.DaemonSet{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &ds)

	return &daemonSet{ds: ds}
}

type DaemonSet interface {
	ToUnstructured() (unstructured.Unstructured, error)
	AnnotateTemplate(string, string)
	LabelTemplate(string, string)
	LabelTemplateIfNotExists(string, string)
	Status() manifest.Status
	Object() *v1.DaemonSet
}

type daemonSet struct {
	ds *v1.DaemonSet
}

func (ds *daemonSet) ToUnstructured() (unstructured.Unstructured, error) {
	u := unstructured.Unstructured{}

	b, err := json.Marshal(ds.ds)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(b, &u.Object)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (ds *daemonSet) Object() *v1.DaemonSet {
	return ds.ds
}

func (ds *daemonSet) AnnotateTemplate(key, value string) {
	annotations := ds.ds.Spec.Template.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	ds.ds.Spec.Template.ObjectMeta.Annotations = annotations
}

func (ds *daemonSet) LabelTemplate(key, value string) {
	labels := ds.ds.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	ds.ds.Spec.Template.ObjectMeta.Labels = labels
}

func (ds *daemonSet) LabelTemplateIfNotExists(key, value string) {
	labels := ds.ds.Spec.Template.ObjectMeta.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	if _, ok := labels[key]; !ok {
		labels[key] = value
	}
	ds.ds.Spec.Template.ObjectMeta.Labels = labels
}

func (ds *daemonSet) Status() manifest.Status {
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

	desiredReplicas := *&ds.ds.Status.DesiredNumberScheduled

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
