package kubernetes

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationSpinnakerTrafficLoadBalancers = "traffic.spinnaker.io/load-balancers"
)

// LoadBalancers returns a slice of load balancers from the annotation
// `traffic.spinnaker.io/load-balancers`. It errors if this annotation
// is not a string slice format like '["service my-service", "service my-service2"]'.
//
// See https://spinnaker.io/docs/reference/providers/kubernetes-v2/#traffic for more info.
func LoadBalancers(u unstructured.Unstructured) ([]string, error) {
	var lbs []string

	annotations := u.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[AnnotationSpinnakerTrafficLoadBalancers]; ok {
			err := json.Unmarshal([]byte(value), &lbs)
			if err != nil {
				return nil,
					fmt.Errorf("error unmarshaling annotation 'traffic.spinnaker.io/load-balancers' "+
						"for resource (kind: %s, name: %s, namespace: %s) into string slice: %v",
						u.GetKind(),
						u.GetName(),
						u.GetNamespace(),
						err)
			}
		}
	}

	return lbs, nil
}
