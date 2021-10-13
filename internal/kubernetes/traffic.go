package kubernetes

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationSpinnakerTrafficLoadBalancers = "traffic.spinnaker.io/load-balancers"
)

// LoadBalancers returns a slice of load balancers from the annotation
// `traffic.spinnaker.io/load-balancers`. It errors if this annotation
// is not a string slice, if it is formatted incorrectly, or if
// the kind is not supported.
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
					fmt.Errorf("error unmarshaling annotation 'traffic.spinnaker.io/load-balancers' for resource (kind: %s, name: %s, namespace: %s) into string slice: %v",
						u.GetKind(),
						u.GetName(),
						u.GetNamespace(),
						err)
			}

			for _, lb := range lbs {
				a := strings.Split(lb, " ")
				if len(a) != 2 {
					return nil,
						fmt.Errorf("Failed to attach load balancer '%s'. Load balancers must be specified in the form '{kind} {name}', e.g. 'service my-service'.", lb)
				}

				kind := a[0]
				if !strings.EqualFold(kind, "service") {
					return nil,
						fmt.Errorf("No support for load balancing via %s exists in Spinnaker.", kind)
				}
			}
		}
	}

	return lbs, nil
}
