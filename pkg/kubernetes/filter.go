package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func FilterOnCluster(items []unstructured.Unstructured, cluster string) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}
	for _, item := range items {
		annotations := item.GetAnnotations()
		if annotations != nil {
			if strings.EqualFold(annotations[AnnotationSpinnakerMonikerCluster], cluster) {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}

func FilterWhereLabelDoesNotExist(items []unstructured.Unstructured, label string) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}
	for _, item := range items {
		labels := item.GetLabels()
		if labels != nil {
			if _, ok := labels[label]; ok {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}
