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
)

type SpinnakerVersion struct {
	Long  string
	Short string
}

//Create a similar function for labels
func (c *controller) AddSpinnakerVersionAnnotations(u *unstructured.Unstructured, application string, version SpinnakerVersion) error {
	annotate(u, AnnotationSpinnakerArtifactVersion, version.Long)
	annotate(u, AnnotationSpinnakerMonikerSequence, version.Short)
	return nil
}

func (c *controller) AddSpinnakerVersionLabels(u *unstructured.Unstructured, application string, version SpinnakerVersion) error {
	label(u, LabelSpinnakerMonikerSequence, version.Short)
	// u.Spec.Template.ObjectMeta

	// .spec.template.metadata.annotations:
	//   artifact.spinnaker.io/version: vNNN
	//   moniker.spinnaker.io/sequence: "N"
	// .spec.template.metadata.labels:
	//   moniker.spinnaker.io/sequence: "N"

	return nil
}

func (c *controller) GetCurrentVersion(ul *unstructured.UnstructuredList, kind, name string) string {
	cluster := kind + " " + name
	currentVersion := "0"
	if len(ul.Items) == 0 {
		return currentVersion
	}
	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	manifestFilter := NewManifestFilter(ul.Items)
	results := manifestFilter.FilterOnCluster(cluster)
	if len(results) == 0 {
		return currentVersion
	}

	//filter out empty moniker.spinnaker.io/sequence labels
	results = manifestFilter.FilterWhereLabelDoesNotExist(LabelSpinnakerSequence)
	if len(results) == 0 {
		return currentVersion
	}

	// For now, we sort on creation timestamp to grab the manifest.
	sort.Slice(results, func(i, j int) bool {
		return results[i].GetCreationTimestamp().String() > results[j].GetCreationTimestamp().String()
	})
	currentVersion = results[0].GetResourceVersion()

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
	latestVersionShortFormat := strconv.Itoa(latestVersionInt)
	latestVersionLongFormat := ""
	if latestVersionInt < 999 {
		latestVersionLongFormat = fmt.Sprintf("v%03d", latestVersionInt)
	} else {
		latestVersionLongFormat = fmt.Sprintf("v%d", latestVersionInt)
	}

	return SpinnakerVersion{
		Short: latestVersionShortFormat,
		Long:  latestVersionLongFormat,
	}
}
