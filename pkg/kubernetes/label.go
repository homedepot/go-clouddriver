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

func (c *controller) AddSpinnakerLabels(u *unstructured.Unstructured, application string) error {
	var err error

	gvk := u.GroupVersionKind()

	// Add reserved labels. Do not overwrite the "qpp.kubernetes.io/name" label
	// as this could affect label selectors.
	//
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-labels
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	label(u, LabelKubernetesManagedBy, spinnaker)
	labelIfNotExists(u, LabelKubernetesName, application)

	if strings.EqualFold(gvk.Kind, "deployment") {
		d := NewDeployment(u.Object)

		// Add reserved labels.
		d.LabelTemplate(LabelKubernetesManagedBy, spinnaker)
		d.LabelTemplateIfNotExists(LabelKubernetesName, application)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := NewReplicaSet(u.Object)

		// Add reserved labels.
		rs.LabelTemplate(LabelKubernetesManagedBy, spinnaker)
		rs.LabelTemplateIfNotExists(LabelKubernetesName, application)

		*u, err = rs.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "daemonset") {
		ds := NewDaemonSet(u.Object)

		// Add reserved labels.
		ds.LabelTemplate(LabelKubernetesManagedBy, spinnaker)
		ds.LabelTemplateIfNotExists(LabelKubernetesName, application)

		*u, err = ds.ToUnstructured()
		if err != nil {
			return err
		}
	}

	return nil
}

// AddSpinnakerVersionLabels adds the `moniker.spinnaker.io/sequence` label
// to the manifest to identify the version number of that resource.
func (c *controller) AddSpinnakerVersionLabels(u *unstructured.Unstructured, version SpinnakerVersion) error {
	var err error

	label(u, LabelSpinnakerMonikerSequence, version.Short)

	gvk := u.GroupVersionKind()

	if strings.EqualFold(gvk.Kind, "deployment") {
		d := NewDeployment(u.Object)

		d.LabelTemplate(LabelSpinnakerMonikerSequence, version.Short)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := NewReplicaSet(u.Object)

		rs.LabelTemplate(LabelSpinnakerMonikerSequence, version.Short)

		*u, err = rs.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "demonset") {
		ds := NewReplicaSet(u.Object)

		ds.LabelTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = ds.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "statefulset") {
		ss := NewStatefulSet(u.Object)

		ss.LabelTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = ss.ToUnstructured()
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
