package kubernetes

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	AnnotationSpinnakerArtifactVersion = `artifact.spinnaker.io/version`
	AnnotationSpinnakerMonikerSequence = `moniker.spinnaker.io/sequence`
	// Maximum latest version before cycling back to V000.
	maxLatestVersion = 999
)

type SpinnakerVersion struct {
	Long  string
	Short string
}

// GetCurrentVersion returns the latest "Spinnaker version" from an unstructured
// list of Kubernetes resources.
func GetCurrentVersion(ul *unstructured.UnstructuredList, kind, name string) string {
	currentVersion := "-1"
	cluster := ""

	if ul == nil || len(ul.Items) == 0 {
		return currentVersion
	}

	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	manifestFilter := NewManifestFilter(ul.Items)
	cluster = kind + " " + name

	results := manifestFilter.FilterOnClusterAnnotation(cluster)
	if len(results) == 0 {
		return currentVersion
	}

	// Filter out empty moniker.spinnaker.io/sequence labels
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

// IsVersioned returns true is a given Kubernetes unstructured resource
// is "versioned". A resource is version if its annotation
// `strategy.spinnaker.io/versioned` is set to "true" or if it is of kind
// Pod, ReplicaSet, ConfigMap, or Secret.
//
// See https://spinnaker.io/reference/providers/kubernetes-v2/#workloads for more info.
func IsVersioned(u unstructured.Unstructured) bool {
	annotations := u.GetAnnotations()
	if annotations != nil {
		if _, ok := annotations[AnnotationSpinnakerStrategyVersioned]; ok {
			return annotations[AnnotationSpinnakerStrategyVersioned] == "true"
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

func IncrementVersion(currentVersion string) SpinnakerVersion {
	currentVersionInt, _ := strconv.Atoi(currentVersion)
	latestVersionInt := currentVersionInt + 1

	if latestVersionInt > maxLatestVersion {
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

func VersionVolumes(u *unstructured.Unstructured, artifacts map[string]clouddriver.TaskCreatedArtifact) error {
	var (
		err     error
		volumes []v1.Volume
	)

	if len(artifacts) > 0 {
		switch strings.ToLower(u.GetKind()) {
		case "deployment":
			d := NewDeployment(u.Object)
			volumes = d.GetSpec().Template.Spec.Volumes

			overwriteVolumesNames(volumes, artifacts)

			*u, err = d.ToUnstructured()
			if err != nil {
				return err
			}

		case "pod":
			p := NewPod(u.Object)
			volumes = p.GetSpec().Volumes

			overwriteVolumesNames(volumes, artifacts)

			*u, err = p.ToUnstructured()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func overwriteVolumesNames(volumes []v1.Volume, artifacts map[string]clouddriver.TaskCreatedArtifact) {
	for _, volume := range volumes {
		if volume.VolumeSource.ConfigMap != nil {
			if val, ok := artifacts[volume.ConfigMap.Name]; ok && strings.EqualFold(val.Type, "kubernetes/configMap") {
				volume.ConfigMap.Name = val.Reference
			}
		}

		if volume.VolumeSource.Secret != nil {
			if val, ok := artifacts[volume.Secret.SecretName]; ok && strings.EqualFold(val.Type, "kubernetes/secret") {
				volume.Secret.SecretName = val.Reference
			}
		}
	}
}
