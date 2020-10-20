package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//go:generate counterfeiter . ManifestFilter
type ManifestFilter interface {
	FilterOnCluster([]unstructured.Unstructured, string) []unstructured.Unstructured
	FilterOnLabel([]unstructured.Unstructured, string) []unstructured.Unstructured
}

type manifestFilter struct {
	items []unstructured.Unstructured
}

func NewManifestFilter(items []unstructured.Unstructured) manifestFilter {
	return manifestFilter{items: items}
}

func (f *manifestFilter) FilterOnCluster(cluster string) []unstructured.Unstructured {
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

func (f *manifestFilter) FilterOnLabel(label string) []unstructured.Unstructured {
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
