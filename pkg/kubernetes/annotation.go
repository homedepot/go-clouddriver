package kubernetes

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationSpinnakerArtifactLocation   = `artifact.spinnaker.io/location`
	AnnotationSpinnakerArtifactName       = `artifact.spinnaker.io/name`
	AnnotationSpinnakerArtifactType       = `artifact.spinnaker.io/type`
	AnnotationSpinnakerMonikerApplication = `moniker.spinnaker.io/application`
	AnnotationSpinnakerMonikerCluster     = `moniker.spinnaker.io/cluster`
	AnnotationSpinnakerStrategyVersioned  = `strategy.spinnaker.io/versioned`
)

// AddSpinnakerAnnotations adds Spinnaker-defined annotations to a given
// unstructured resource.
func AddSpinnakerAnnotations(u *unstructured.Unstructured, application string) error {
	var err error

	name := u.GetName()
	namespace := u.GetNamespace()
	gvk := u.GroupVersionKind()

	// possible bug ToLower
	t := fmt.Sprintf("kubernetes/%s", strcase.ToLowerCamel(gvk.Kind))
	cluster := fmt.Sprintf("%s %s", strcase.ToLowerCamel(gvk.Kind), name)

	// Add reserved annotations.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-annotations
	annotate(u, AnnotationSpinnakerArtifactLocation, namespace)
	annotate(u, AnnotationSpinnakerArtifactName, name)
	annotate(u, AnnotationSpinnakerArtifactType, t)
	annotate(u, AnnotationSpinnakerMonikerApplication, application)
	annotate(u, AnnotationSpinnakerMonikerCluster, cluster)

	if strings.EqualFold(gvk.Kind, "deployment") {
		d := NewDeployment(u.Object)

		// Add spinnaker annotations to the deployment pod template.
		d.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		d.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		d.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := NewReplicaSet(u.Object)

		// Add spinnaker annotations to the replicaset pod template.
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		*u, err = rs.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "daemonset") {
		ds := NewDaemonSet(u.Object)

		// Add spinnaker annotations to the daemonset pod template.
		ds.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		ds.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		ds.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		ds.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		ds.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		*u, err = ds.ToUnstructured()
		if err != nil {
			return err
		}
	}

	return nil
}

// AddSpinnakerVersionAnnotations adds the following annotations:
// `artifact.spinnaker.io/version`
// `moniker.spinnaker.io/sequence`
// to the manifest to identify the version number of that resource.
func AddSpinnakerVersionAnnotations(u *unstructured.Unstructured, version SpinnakerVersion) error {
	var err error

	annotate(u, AnnotationSpinnakerArtifactVersion, version.Long)
	annotate(u, AnnotationSpinnakerMonikerSequence, version.Short)

	gvk := u.GroupVersionKind()

	if strings.EqualFold(gvk.Kind, "deployment") {
		d := NewDeployment(u.Object)

		d.AnnotateTemplate(AnnotationSpinnakerArtifactVersion, version.Long)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := NewReplicaSet(u.Object)

		rs.AnnotateTemplate(AnnotationSpinnakerArtifactVersion, version.Long)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = rs.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "daemonset") {
		ds := NewReplicaSet(u.Object)

		ds.AnnotateTemplate(AnnotationSpinnakerArtifactVersion, version.Long)
		ds.AnnotateTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = ds.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "statefulset") {
		ss := NewStatefulSet(u.Object)

		ss.AnnotateTemplate(AnnotationSpinnakerArtifactVersion, version.Long)
		ss.AnnotateTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = ss.ToUnstructured()
		if err != nil {
			return err
		}
	}

	return nil
}

func annotate(o *unstructured.Unstructured, key, value string) {
	annotations := o.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	o.SetAnnotations(annotations)
}
