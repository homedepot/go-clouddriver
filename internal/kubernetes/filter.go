package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func NewManifestFilter(items []unstructured.Unstructured) *ManifestFilter {
	return &ManifestFilter{items: items}
}

type ManifestFilter struct {
	items []unstructured.Unstructured
}

func (f *ManifestFilter) FilterOnClusterAnnotation(cluster string) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}

	for _, item := range f.items {
		annotations := item.GetAnnotations()
		if annotations != nil {
			if strings.EqualFold(annotations[AnnotationSpinnakerMonikerCluster], cluster) {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}

func (f *ManifestFilter) FilterOnLabel(label string) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}

	for _, item := range f.items {
		labels := item.GetLabels()
		if labels != nil {
			if _, ok := labels[label]; ok {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}
