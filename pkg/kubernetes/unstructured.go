package kubernetes

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
)

func (c *controller) ToUnstructured(manifest map[string]interface{}) (*unstructured.Unstructured, error) {
	b, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(b, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convert the runtime.Object to unstructured.Unstructured.
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: m,
	}, nil
}

func SetDefaultNamespaceIfScopedAndNoneSet(u *unstructured.Unstructured, helper *resource.Helper) {
	namespace := u.GetNamespace()
	if helper.NamespaceScoped && namespace == "" {
		namespace = "default"
		u.SetNamespace(namespace)
	}
}

func SetNamespaceIfScoped(namespace string, u *unstructured.Unstructured, helper *resource.Helper) {
	if helper.NamespaceScoped {
		u.SetNamespace(namespace)
	}
}
