package kubernetes

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	listTimeout = int64(30)
)

const (
	AnnotationSpinnakerArtifactVersion = `artifact.spinnaker.io/version`
	AnnotationSpinnakerMonikerSequence = `moniker.spinnaker.io/sequence`
	Annotation
)

type SpinnakerVersion struct {
	Long  string
	Short string
}

func (c *controller) GetCurrentVersion(ul *unstructured.UnstructuredList, kind, name string) string {
	currentVersion := "-1"
	cluster := ""

	if len(ul.Items) == 0 {
		return currentVersion
	}

	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	manifestFilter := NewManifestFilter(ul.Items)

	re := regexp.MustCompile(`(.*)-v(\d){3}$`)
	subm := re.FindSubmatch([]byte(name))

	if subm == nil {
		cluster = kind + " " + name
	} else {
		cluster = string(subm[1])
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

func (c *controller) VersionVolumes(u *unstructured.Unstructured, namespace, application string, kubeClient Client) error {
	var (
		err     error
		volumes []v1.Volume
	)

	switch strings.ToLower(u.GetKind()) {
	case "deployment":
		d := NewDeployment(u.Object)
		volumes = d.GetSpec().Template.Spec.Volumes

		err = c.OverwriteVolumeNames(volumes, namespace, application, kubeClient)
		if err != nil {
			return err
		}

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}

	case "pod":
		p := NewPod(u.Object)
		volumes = p.GetSpec().Volumes

		err = c.OverwriteVolumeNames(volumes, namespace, application, kubeClient)
		if err != nil {
			return err
		}

		*u, err = p.ToUnstructured()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *controller) OverwriteVolumeNames(volumes []v1.Volume, namespace, application string, kubeClient Client) error {
	for _, volume := range volumes {
		if volume.VolumeSource.ConfigMap != nil {
			currentVersion, err := c.GetVolumeVersion(volume.ConfigMap.Name, "configMap", namespace, application, kubeClient)
			if err != nil {
				return err
			}

			volume.ConfigMap.Name = currentVersion
		}

		if volume.VolumeSource.Secret != nil {
			currentVersion, err := c.GetVolumeVersion(volume.Secret.SecretName, "secret", namespace, application, kubeClient)
			if err != nil {
				return err
			}

			volume.Secret.SecretName = currentVersion
		}
	}

	return nil
}

func (c *controller) GetVolumeVersion(name, kind, namespace, application string, kubeClient Client) (string, error) {
	//get resourcces
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			LabelKubernetesName:      application,
			LabelKubernetesManagedBy: Spinnaker,
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      LabelSpinnakerMonikerSequence,
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}

	lo := metav1.ListOptions{
		LabelSelector:  labels.Set(labelSelector.MatchLabels).String(),
		TimeoutSeconds: &listTimeout,
	}

	results, err := kubeClient.ListResourcesByKindAndNamespace(kind, namespace, lo)
	if err != nil {
		return "", err
	}

	cluster := kind + " " + name

	manifestFilter := NewManifestFilter(results.Items)
	resources := manifestFilter.FilterOnClusterAnnotation(cluster)
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].GetCreationTimestamp().String() > resources[j].GetCreationTimestamp().String()
	})

	return resources[0].GetName(), nil
}
