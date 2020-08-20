package kubernetes

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

func label(o *unstructured.Unstructured, key, value string) {
	labels := o.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	o.SetLabels(labels)
}
