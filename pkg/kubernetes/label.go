package kubernetes

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	LabelKubernetesSpinnakerApp = `app.kubernetes.io/spinnaker-app`
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	LabelKubernetesName      = `app.kubernetes.io/name`
	LabelKubernetesManagedBy = `app.kubernetes.io/managed-by`
)

func AddSpinnakerLabels(u *unstructured.Unstructured, application string) error {
	var err error

	gvk := u.GroupVersionKind()

	// Add reserved labels. Had some trouble with setting the kubernetes name as
	// this interferes with label selectors, so I changed that to be spinnaker-app.
	//
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-labels
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	// label(u, LabelKubernetesName, application)
	label(u, LabelKubernetesSpinnakerApp, application)
	label(u, LabelKubernetesManagedBy, spinnaker)

	if strings.EqualFold(gvk.Kind, "deployment") {
		d := NewDeployment(u.Object)

		// Add reserved labels.
		// d.LabelTemplate(LabelKubernetesName, application)
		d.LabelTemplate(LabelKubernetesSpinnakerApp, application)
		d.LabelTemplate(LabelKubernetesManagedBy, spinnaker)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := NewReplicaSet(u.Object)

		// Add reserved labels.
		// rs.LabelTemplate(LabelKubernetesName, application)
		rs.LabelTemplate(LabelKubernetesSpinnakerApp, application)
		rs.LabelTemplate(LabelKubernetesManagedBy, spinnaker)

		*u, err = rs.ToUnstructured()
		if err != nil {
			return err
		}
	}

	return nil
}

func label(o *unstructured.Unstructured, key, value string) {
	labels := o.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key] = value
	o.SetLabels(labels)
}
