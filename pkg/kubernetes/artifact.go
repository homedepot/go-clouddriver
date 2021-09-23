package kubernetes

import (
	"strings"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ReplaceDockerImageArtifacts finds containers in a given unstructured
// Kubernetes object and replaces the `image` value of the container
// with the artifact reference if the artiface name matches the
// container's `image` value.
//
// Pods define containers at the JSON path ".spec.containers",
// all other kinds defines containers at the JSON path ".spec.template.spec.containers".
func ReplaceDockerImageArtifacts(u *unstructured.Unstructured,
	artifacts map[string]clouddriver.TaskCreatedArtifact) error {
	var (
		containers []interface{}
		found      bool
		err        error
	)

	if len(artifacts) == 0 {
		return nil
	}

	if strings.EqualFold(u.GetKind(), "pod") {
		containers, found, err = unstructured.NestedSlice(u.Object, "spec", "containers")
	} else {
		containers, found, err = unstructured.NestedSlice(u.Object, "spec", "template", "spec", "containers")
	}
	// An error is thrown if the nested slice is found but is not of type
	// []interface{}.
	if err != nil {
		return err
	}
	// If the nested slice is not found, return nil as this Kubernetes
	// kind might not have containers.
	if !found {
		return nil
	}

	overwriteContainerImages(containers, artifacts)

	if strings.EqualFold(u.GetKind(), "pod") {
		err = unstructured.SetNestedSlice(u.Object, containers, "spec", "containers")
	} else {
		err = unstructured.SetNestedSlice(u.Object, containers, "spec", "template", "spec", "containers")
	}

	return err
}

func overwriteContainerImages(containers []interface{}, artifacts map[string]clouddriver.TaskCreatedArtifact) {
	for _, container := range containers {
		c, ok := container.(map[string]interface{})
		if !ok {
			continue
		}

		image, found, err := unstructured.NestedString(c, "image")
		if err == nil && found {
			if artifact, ok := artifacts[image]; ok && strings.EqualFold(artifact.Type, "docker/image") {
				_ = unstructured.SetNestedField(c, artifact.Reference, "image")
			}
		}
	}
}
