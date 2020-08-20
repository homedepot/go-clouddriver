package kubernetes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func annotate(o *unstructured.Unstructured, key, value string) {
	annotations := o.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	o.SetAnnotations(annotations)
}
