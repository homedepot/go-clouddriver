package kubernetes

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func NewRequirement(op string, key string, values []string) (*labels.Requirement, error) {
	switch strings.ToLower(op) {
	case "contains":
		return labels.NewRequirement(key, selection.In, values)
	case "not_contains":
		return labels.NewRequirement(key, selection.NotIn, values)
	case "equals":
		return labels.NewRequirement(key, selection.Equals, values)
	case "not_equals":
		return labels.NewRequirement(key, selection.NotEquals, values)
	case "exists":
		return labels.NewRequirement(key, selection.Exists, values)
	case "not_exists":
		return labels.NewRequirement(key, selection.DoesNotExist, values)
	default:
		return nil, fmt.Errorf("operator '%v' is not recognized", op)
	}
}

// DefaultLabelSelector returns the label selector
// `app.kubernetes.io/managed-by in (spinnaker,spinnaker-operator)`,
// which allows us to list all resources with a label selector
// managed by Spinnaker or Spinnaker Operator.
func DefaultLabelSelector() string {
	labelSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      LabelKubernetesManagedBy,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"spinnaker", "spinnaker-operator"},
			},
		},
	}
	// Since this is a defined label Selector we can ignore the error.
	ls, _ := metav1.LabelSelectorAsSelector(&labelSelector)

	return ls.String()
}
