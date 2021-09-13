package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	LabelKubernetesName           = `app.kubernetes.io/name`
	LabelKubernetesManagedBy      = `app.kubernetes.io/managed-by`
	LabelSpinnakerMonikerSequence = `moniker.spinnaker.io/sequence`
)

// AddSpinnakerLabels labels a given unstructured Kubernetes resource
// with Spinnaker defined labels.
func AddSpinnakerLabels(u *unstructured.Unstructured, application string) error {
	// Add reserved labels. Do not overwrite the "qpp.kubernetes.io/name" label
	// as this could affect label selectors.
	//
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-labels
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	label(u, LabelKubernetesManagedBy, spinnaker)
	labelIfNotExists(u, LabelKubernetesName, application)

	if strings.EqualFold(u.GetKind(), "deployment") ||
		strings.EqualFold(u.GetKind(), "replicaset") ||
		strings.EqualFold(u.GetKind(), "daemonset") {
		err := labelTemplate(u, LabelKubernetesManagedBy, spinnaker)
		if err != nil {
			return err
		}

		err = labelTemplateIfNotExists(u, LabelKubernetesName, application)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddSpinnakerVersionLabels adds the `moniker.spinnaker.io/sequence` label
// to the manifest to identify the version number of that resource.
func AddSpinnakerVersionLabels(u *unstructured.Unstructured, version SpinnakerVersion) error {
	label(u, LabelSpinnakerMonikerSequence, version.Short)

	if strings.EqualFold(u.GetKind(), "deployment") ||
		strings.EqualFold(u.GetKind(), "replicaset") ||
		strings.EqualFold(u.GetKind(), "statefulset") ||
		strings.EqualFold(u.GetKind(), "daemonset") {
		err := labelTemplate(u, LabelSpinnakerMonikerSequence, version.Short)
		if err != nil {
			return err
		}
	}

	return nil
}

func label(u *unstructured.Unstructured, key, value string) {
	labels := u.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	u.SetLabels(labels)
}

func labelTemplate(u *unstructured.Unstructured, key, value string) error {
	templateLabels, found, err := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
	if err != nil {
		return err
	}

	if !found {
		templateLabels = map[string]string{}
	}

	templateLabels[key] = value

	err = unstructured.SetNestedStringMap(u.Object, templateLabels, "spec", "template", "metadata", "labels")
	if err != nil {
		return err
	}

	return nil
}

func labelIfNotExists(u *unstructured.Unstructured, key, value string) {
	labels := u.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	if _, ok := labels[key]; !ok {
		labels[key] = value
	}

	u.SetLabels(labels)
}

func labelTemplateIfNotExists(u *unstructured.Unstructured, key, value string) error {
	templateLabels, found, err := unstructured.NestedStringMap(u.Object, "spec", "template", "metadata", "labels")
	if err != nil {
		return err
	}

	if !found {
		templateLabels = map[string]string{}
	}

	if templateLabels[key] != "" {
		return nil
	}

	templateLabels[key] = value

	err = unstructured.SetNestedStringMap(u.Object, templateLabels, "spec", "template", "metadata", "labels")
	if err != nil {
		return err
	}

	return nil
}
