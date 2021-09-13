package kubernetes

import (
	"fmt"
	"strings"

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
