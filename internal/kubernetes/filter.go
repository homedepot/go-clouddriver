package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// FilterOnAnnotations takes a slice of unstructured and returns
// the filtered slice based on a given annotation key and value.
func FilterOnAnnotation(items []unstructured.Unstructured,
	annotationKey, annotationValue string) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}

	for _, item := range items {
		annotations := item.GetAnnotations()
		if annotations != nil &&
			strings.EqualFold(annotations[annotationKey], annotationValue) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// FilterOnLabelExists takes a slice of unstructured and returns
// a filtered slice where a given label exists in each
// of the unstructured objects.
func FilterOnLabelExists(items []unstructured.Unstructured,
	label string) []unstructured.Unstructured {
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
