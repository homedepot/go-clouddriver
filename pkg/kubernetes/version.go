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
	LongVersion  string
	ShortVersion string
}

//Create a similar function for labels
func (c *controller) AddSpinnakerVersionAnnotations(u *unstructured.Unstructured, application string, version SpinnakerVersion) error {
	// var err error

	annotate(u, AnnotationSpinnakerArtifactVersion, version.LongVersion)
	annotate(u, AnnotationSpinnakerMonikerSequence, version.ShortVersion)

	//ToDo
	// Add spinnaker versioning labels and annotations.
	// 		.metadata.annotations:
	//   artifact.spinnaker.io/version: vNNN
	//   moniker.spinnaker.io/sequence: "N"
	// .metadata.labels:
	//   moniker.spinnaker.io/sequence: "N"
	// .spec.template.metadata.annotations:
	//   artifact.spinnaker.io/version: vNNN
	//   moniker.spinnaker.io/sequence: "N"
	// .spec.template.metadata.labels:
	//   moniker.spinnaker.io/sequence: "N"

	return nil
}

func (c *controller) GetCurrentVersion(ul *unstructured.UnstructuredList, kind, name string) (string, error) {

	cluster := kind + " " + name
	currentVersion := "0"
	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	results := FilterOnCluster(ul.Items, cluster)
	if len(results) == 0 {
		return currentVersion, nil
	}

	//filter out empty moniker.spinnaker.io/sequence labels
	results = FilterWhereLabelDoesNotExist(results, LabelSpinnakerSequence)
	if len(results) == 0 {
		return currentVersion, nil
	}

	// For now, we sort on creation timestamp to grab the manifest.
	sort.Slice(results, func(i, j int) bool {
		return results[i].GetCreationTimestamp().String() > results[j].GetCreationTimestamp().String()
	})
	currentVersion = results[0].GetResourceVersion()

	return currentVersion, nil
}

func (c *controller) IsVersioned(u *unstructured.Unstructured) bool {
	annotations := u.GetAnnotations()
	if annotations != nil {
		if _, ok := annotations[AnnotationSpinnakerVersioned]; ok {
			return true
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
	latestVersionLongFormat := fmt.Sprintf("v%03d", latestVersionInt)
	return SpinnakerVersion{
		ShortVersion: latestVersionShortFormat,
		LongVersion:  latestVersionLongFormat,
	}
}
