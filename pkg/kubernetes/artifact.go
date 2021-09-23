package kubernetes

import (
	"strings"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// BindDockerImageArtifacts finds containers and init containers in a given unstructured
// Kubernetes object and replaces the `image` value of the container
// with the artifact reference if the artifact name matches the
// container's `image` value.
//
// Pods define containers at the JSON path ".spec.containers" and
// init containers at ".spec.initContainers".
// All other kinds define containers at the JSON path ".spec.template.spec.containers"
// and init containers at ".spec.template.spec.initContainers".
//
// Java source code here:
// https://github.com/spinnaker/clouddriver/blob/58ab154b0ec0d62772201b5b319af349498a4e3f/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/artifact/Replacer.java#L166
func BindDockerImageArtifacts(u *unstructured.Unstructured,
	artifacts map[string]clouddriver.TaskCreatedArtifact) error {
	// Bind artifact to the "containers" field.
	err := bindDockerImageArtifacts(u, artifacts, "containers")
	if err != nil {
		return err
	}
	// Bind artifact to the "initContainers" field.
	err = bindDockerImageArtifacts(u, artifacts, "initContainers")
	if err != nil {
		return err
	}

	return nil
}

func bindDockerImageArtifacts(u *unstructured.Unstructured,
	artifacts map[string]clouddriver.TaskCreatedArtifact, field string) error {
	var (
		containers []interface{}
		found      bool
		err        error
	)

	if len(artifacts) == 0 {
		return nil
	}

	if strings.EqualFold(u.GetKind(), "pod") {
		containers, found, err = unstructured.NestedSlice(u.Object, "spec", field)
	} else {
		containers, found, err = unstructured.NestedSlice(u.Object, "spec", "template", "spec", field)
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
		err = unstructured.SetNestedSlice(u.Object, containers, "spec", field)
	} else {
		err = unstructured.SetNestedSlice(u.Object, containers, "spec", "template", "spec", field)
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
