package kubernetes

import (
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationSpinnakerMaxVersionHistory = "strategy.spinnaker.io/max-version-history"
	AnnotationSpinnakerRecreate          = "strategy.spinnaker.io/recreate"
	AnnotationSpinnakerReplaced          = "strategy.spinnaker.io/replace"
	AnnotationSpinnakerUseSourceCapacity = "strategy.spinnaker.io/use-source-capacity"
)

// MaxVersionHistory returns true if the value of the annotation
// `strategy.spinnaker.io/max-version-history` of the given Kubernetes
// unstructured resource, or 0 if annotation is not present.
//
// See https://spinnaker.io/docs/reference/providers/kubernetes-v2/#strategy for more info.
func MaxVersionHistory(u unstructured.Unstructured) (maxVersionHistory int, err error) {
	maxVersionHistory = 0

	annotations := u.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[AnnotationSpinnakerMaxVersionHistory]; ok {
			maxVersionHistory, err = strconv.Atoi(value)
			if err != nil {
				return
			}
		}
	}

	return
}

// Recreate returns true if the given Kubernetes unstructured resource
// has the annotation `strategy.spinnaker.io/recreate` set to "true".
//
// See https://spinnaker.io/docs/reference/providers/kubernetes-v2/#strategy for more info.
func Recreate(u unstructured.Unstructured) bool {
	annotations := u.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[AnnotationSpinnakerRecreate]; ok {
			return value == "true"
		}
	}

	return false
}

// Replace returns true if the given Kubernetes unstructured resource
// has the annotation `strategy.spinnaker.io/replace` set to "true".
//
// See https://spinnaker.io/docs/reference/providers/kubernetes-v2/#strategy for more info.
func Replace(u unstructured.Unstructured) bool {
	annotations := u.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[AnnotationSpinnakerReplaced]; ok {
			return value == "true"
		}
	}

	return false
}

// UseSourceCapacity returns true is a given Kubernetes unstructured resource
// has the annotation `strategy.spinnaker.io/use-source-capacity` set to "true"
// and it is of kind Deployment, ReplicaSet, or StatefulSet.
//
// See https://spinnaker.io/docs/reference/providers/kubernetes-v2/#strategy for more info.
func UseSourceCapacity(u unstructured.Unstructured) bool {
	switch strings.ToLower(u.GetKind()) {
	case "deployment", "replicaset", "statefulset":
		annotations := u.GetAnnotations()
		if annotations != nil {
			if value, ok := annotations[AnnotationSpinnakerUseSourceCapacity]; ok {
				return value == "true"
			}
		}
	}

	return false
}
