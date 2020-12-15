package kubernetes

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	AnnotationSpinnakerArtifactVersion = `artifact.spinnaker.io/version`
	AnnotationSpinnakerMonikerSequence = `moniker.spinnaker.io/sequence`
	Annotation
)

type SpinnakerVersion struct {
	Long  string
	Short string
}

func (c *controller) AddSpinnakerVersionAnnotations(u *unstructured.Unstructured, version SpinnakerVersion) error {
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

func (c *controller) AddSpinnakerVersionLabels(u *unstructured.Unstructured, version SpinnakerVersion) error {
	var err error

	label(u, LabelSpinnakerMonikerSequence, version.Short)

	gvk := u.GroupVersionKind()

	if strings.EqualFold(gvk.Kind, "deployment") {
		d := NewDeployment(u.Object)

		d.LabelTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := NewReplicaSet(u.Object)

		rs.LabelTemplate(AnnotationSpinnakerMonikerSequence, version.Short)

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

func (c *controller) GetCurrentVersion(ul *unstructured.UnstructuredList, kind, name string) string {
	currentVersion := "-1"
	cluster := ""

	if len(ul.Items) == 0 {
		return currentVersion
	}

	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	manifestFilter := NewManifestFilter(ul.Items)
	lastIndex := strings.LastIndex(name, "-v")

	if lastIndex != -1 {
		cluster = kind + " " + name[:lastIndex]
	} else {
		cluster = kind + " " + name
	}

	results := manifestFilter.FilterOnClusterAnnotation(cluster)
	if len(results) == 0 {
		return currentVersion
	}

	//filter out empty moniker.spinnaker.io/sequence labels
	manifestFilter2 := NewManifestFilter(results)
	results = manifestFilter2.FilterOnLabel(LabelSpinnakerMonikerSequence)

	if len(results) == 0 {
		return currentVersion
	}

	// For now, we sort on creation timestamp to grab the manifest.
	sort.Slice(results, func(i, j int) bool {
		return results[i].GetCreationTimestamp().String() > results[j].GetCreationTimestamp().String()
	})

	annotations := results[0].GetAnnotations()
	currentVersion = annotations[AnnotationSpinnakerMonikerSequence]

	return currentVersion
}

func (c *controller) IsVersioned(u *unstructured.Unstructured) bool {
	annotations := u.GetAnnotations()
	if annotations != nil {
		if _, ok := annotations[AnnotationSpinnakerStrategyVersioned]; ok {
			if annotations[AnnotationSpinnakerStrategyVersioned] == "true" {
				return true
			} else {
				return false
			}
		}
	}

	kind := strings.ToLower(u.GetKind())
	if strings.EqualFold(kind, "pod") ||
		strings.EqualFold(kind, "replicaSet") ||
		strings.EqualFold(kind, "ConfigMap") ||
		strings.EqualFold(kind, "Secret") {
		return true
	}

	return false
}

func (c *controller) IncrementVersion(currentVersion string) SpinnakerVersion {
	currentVersionInt, _ := strconv.Atoi(currentVersion)
	latestVersionInt := currentVersionInt + 1

	if latestVersionInt > 999 {
		latestVersionInt = 0
	}

	latestVersionShortFormat := strconv.Itoa(latestVersionInt)
	latestVersionLongFormat := ""
	latestVersionLongFormat = fmt.Sprintf("v%03d", latestVersionInt)

	return SpinnakerVersion{
		Short: latestVersionShortFormat,
		Long:  latestVersionLongFormat,
	}
}
