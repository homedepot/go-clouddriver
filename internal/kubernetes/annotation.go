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
	AnnotationSpinnakerMonikerDetail      = `moniker.spinnaker.io/detail`
	AnnotationSpinnakerMonikerStack       = `moniker.spinnaker.io/stack`
	AnnotationSpinnakerStrategyVersioned  = `strategy.spinnaker.io/versioned`
	AnnotationSpinnakerServerSideApply    = `strategy.spinnaker.io/server-side-apply`
)

// AddSpinnakerAnnotations adds Spinnaker-defined annotations to a given
// unstructured resource.
func AddSpinnakerAnnotations(u *unstructured.Unstructured, application string) error {
	name := u.GetName()
	t := fmt.Sprintf("kubernetes/%s", strcase.ToLowerCamel(u.GetKind()))
	cluster := fmt.Sprintf("%s %s", strcase.ToLowerCamel(u.GetKind()), name)
	namespace := u.GetNamespace()

	// Add reserved annotations.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-annotations
	annotate(u, AnnotationSpinnakerArtifactLocation, namespace)
	annotate(u, AnnotationSpinnakerArtifactName, name)
	annotate(u, AnnotationSpinnakerArtifactType, t)
	annotate(u, AnnotationSpinnakerMonikerApplication, application)
	annotate(u, AnnotationSpinnakerMonikerCluster, cluster)

	if strings.EqualFold(u.GetKind(), "deployment") ||
		strings.EqualFold(u.GetKind(), "replicaset") ||
		strings.EqualFold(u.GetKind(), "daemonset") {
		err := AnnotateTemplate(u, AnnotationSpinnakerArtifactLocation, namespace)
		if err != nil {
			return err
		}

		err = AnnotateTemplate(u, AnnotationSpinnakerArtifactName, name)
		if err != nil {
			return err
		}

		err = AnnotateTemplate(u, AnnotationSpinnakerArtifactType, t)
		if err != nil {
			return err
		}

		err = AnnotateTemplate(u, AnnotationSpinnakerMonikerApplication, application)
		if err != nil {
			return err
		}

		err = AnnotateTemplate(u, AnnotationSpinnakerMonikerCluster, cluster)
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
	annotate(u, AnnotationSpinnakerArtifactVersion, version.Long)
	annotate(u, AnnotationSpinnakerMonikerSequence, version.Short)

	if strings.EqualFold(u.GetKind(), "deployment") ||
		strings.EqualFold(u.GetKind(), "replicaset") ||
		strings.EqualFold(u.GetKind(), "statefulset") ||
		strings.EqualFold(u.GetKind(), "daemonset") {
		err := AnnotateTemplate(u, AnnotationSpinnakerArtifactVersion, version.Long)
		if err != nil {
			return err
		}

		err = AnnotateTemplate(u, AnnotationSpinnakerMonikerSequence, version.Short)
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

// AnnotateTemplate annotates the nested string map located at
// .spec.template.metadata.annotations.
func AnnotateTemplate(u *unstructured.Unstructured, key, value string) error {
	templateAnnotations, found, err := unstructured.NestedStringMap(u.Object,
		"spec", "template", "metadata", "annotations")
	if err != nil {
		return err
	}

	if !found {
		templateAnnotations = map[string]string{}
	}

	templateAnnotations[key] = value

	err = unstructured.SetNestedStringMap(u.Object, templateAnnotations, "spec", "template", "metadata", "annotations")
	if err != nil {
		return err
	}

	return nil
}

// SpinnakerMonikerApplication returns the value the annotation
// `moniker.spinnaker.io/application` of the given Kubernetes
// unstructured resource, or an empty string if not present.
func SpinnakerMonikerApplication(u unstructured.Unstructured) string {
	annotations := u.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[AnnotationSpinnakerMonikerApplication]; ok {
			return value
		}
	}

	return ""
}

// AnnotationMatches returns true if the value of annotation
// matches the given value.
//
// If value is "", then match when annotation is missing or also "".
// If value is "*", then always match.
// If value is "something", then match only when annotation value is "something".
//
// See https://github.com/spinnaker/clouddriver/blob/96755fec0c04b6e361efb6d1c19a7afc3926e302/clouddriver-core/src/main/java/com/netflix/spinnaker/clouddriver/core/ProjectClustersService.java#L197
func AnnotationMatches(u unstructured.Unstructured, annotationKey, value string) bool {
	annotationValue := ""

	annotations := u.GetAnnotations()
	if annotations != nil {
		annotationValue = annotations[annotationKey]
	}

	return value == "*" || strings.EqualFold(value, annotationValue)
}
