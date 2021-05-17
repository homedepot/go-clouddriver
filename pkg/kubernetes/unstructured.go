package kubernetes

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
)

// ToUnstructured converts a map[string]interface{} to a kubernetes unstructured.Unstructured
// object. An unstructured's "Object" field is technically just a map[string]interface{},
// so we could do the following:
//
// return unstructured.Unstructured{ Object: manifest }
//
// and not have this function return an error, but we miss out on some validation that
// happens during `kubectl`, such as when you attempt to apply a bad manifest, like something
// without "kind" specified. Example:
//
// $ k apply -f bad.yaml --validate=false
// error: unable to decode "bad.yaml": Object 'Kind' is missing in '{}'
//
// If we decide to not cycle through the encoding/decoding process, these errors
// will be deferred to when the manifest gets applied, and will fail with some error like
//
// error applying manifest (kind: , apiVersion: v1, name: bad-gke): no matches for kind "" in version "apps/v1"
//
// For now, we are not deferring the error.
func ToUnstructured(manifest map[string]interface{}) (unstructured.Unstructured, error) {
	b, err := json.Marshal(manifest)
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(b, nil, nil)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	// Convert the runtime.Object to unstructured.Unstructured.
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	return unstructured.Unstructured{
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
