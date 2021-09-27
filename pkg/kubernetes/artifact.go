package kubernetes

import (
	"strings"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/artifact"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	// Artifact type kubernetes/configMap JSON path locations.
	jsonPathConfigMapVolume                    = `.spec.template.spec.volumes.*.configMap.name`
	jsonPathConfigMapPodVolume                 = `.spec.volumes.*.configMap.name`
	jsonPathConfigMapProjectedVolume           = `.spec.template.spec.volumes.*.projected.sources.*.configMap.name`
	jsonPathConfigMapPodProjectedVolume        = `.spec.volumes.*.projected.sources.*.configMap.name`
	jsonPathConfigMapKeyValueContainers        = `.spec.template.spec.containers.*.env.*.valueFrom.configMapKeyRef.name`
	jsonPathConfigMapKeyValueContainersPod     = `.spec.containers.*.env.*.valueFrom.configMapKeyRef.name`
	jsonPathConfigMapKeyValueInitContainers    = `.spec.template.spec.initContainers.*.env.*.valueFrom.configMapKeyRef.name`
	jsonPathConfigMapKeyValueInitContainersPod = `.spec.initContainers.*.env.*.valueFrom.configMapKeyRef.name`
	jsonPathConfigMapEnvContainers             = `.spec.template.spec.containers.*.envFrom.*.configMapRef.name`
	jsonPathConfigMapEnvContainersPod          = `.spec.containers.*.envFrom.*.configMapRef.name`
	jsonPathConfigMapEnvInitContainers         = `.spec.template.spec.initContainers.*.envFrom.*.configMapRef.name`
	jsonPathConfigMapEnvInitContainersPod      = `.spec.initContainers.*.envFrom.*.configMapRef.name`
	// Artifact type docker/image JSON path locations.
	jsonPathDockerImageContainers        = `.spec.template.spec.containers.*.image`
	jsonPathDockerImageContainersPod     = `.spec.containers.*.image`
	jsonPathDockerImageInitContainers    = `.spec.template.spec.initContainers.*.image`
	jsonPathDockerImageInitContainersPod = `.spec.initContainers.*.image`
	// Artifact type kubernetes/secret JSON path locations.
	jsonPathSecretVolume                    = `.spec.template.spec.volumes.*.secret.secretName`
	jsonPathSecretVolumePod                 = `.spec.volumes.*.secret.secretName`
	jsonPathSecretProjectedVolume           = `.spec.template.spec.volumes.*.projected.sources.*.secret.name`
	jsonPathSecretProjectedVolumePod        = `.spec.volumes.*.projected.sources.*.secret.name`
	jsonPathSecretKeyValueContainers        = `.spec.template.spec.containers.*.env.*.valueFrom.secretKeyRef.name`
	jsonPathSecretKeyValueContainersPod     = `.spec.containers.*.env.*.valueFrom.secretKeyRef.name`
	jsonPathSecretKeyValueInitContainers    = `.spec.template.spec.initContainers.*.env.*.valueFrom.secretKeyRef.name`
	jsonPathSecretKeyValueInitContainersPod = `.spec.initContainers.*.env.*.valueFrom.secretKeyRef.name`
	jsonPathSecretEnvContainers             = `.spec.template.spec.containers.*.envFrom.*.secretRef.name`
	jsonPathSecretEnvContainersPod          = `.spec.containers.*.envFrom.*.secretRef.name`
	jsonPathSecretEnvInitContainers         = `.spec.template.spec.initContainers.*.envFrom.*.secretRef.name`
	jsonPathSecretEnvInitContainersPod      = `.spec.initContainers.*.envFrom.*.secretRef.name`
	// Artifact type kubernetes/deployment and kubernetes/replicaSet JSON path locations.
	jsonPathHPAKind = `.spec.scaleTargetRef.kind`
	jsonPathHPAName = `.spec.scaleTargetRef.name`
)

// BindArtifacts takes in an unstructured Kubernetes object and a slice of artifacts
// then binds these artifacts to any applicable JSON path for the given artifact.
//
// For example, take the following manifest that references a configMap with name 'replace-me':
//
// apiVersion: v1
// kind: Pod
// metadata:
//   name: dapi-test-pod
// spec:
//   containers:
//     - name: test-container
//       image: k8s.gcr.io/busybox
//       volumeMounts:
//       - name: my-volume
//         mountPath: /etc/config
//   volumes:
//     - name: my-volume
//       configMap:
//         # Provide the name of the ConfigMap containing the files you want
//         # to add to the container
//         name: replace-me
//   restartPolicy: Never
//
// Now let's say we pass in the following Clouddriver Artifact:
// {
//   "name": "replace-me",
//   "reference": "my-config-map-v000"
// }
//
// This would result in the JSON path '.spec.volumes[0].configMap.name' changing from
// 'replace-me' to 'my-config-map-v000'.
//
// The source code for these Replacers can be found here:
// https://github.com/spinnaker/clouddriver/blob/4d4e01084ac5259792020e419b1af7686ab38019/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/artifact/Replacer.java#L150
func BindArtifacts(u *unstructured.Unstructured,
	artifacts []clouddriver.Artifact) {
	for _, a := range artifacts {
		switch a.Type {
		case artifact.TypeDockerImage:
			bindArtifact(u.Object, a, splitIterables(jsonPathDockerImageContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathDockerImageContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathDockerImageInitContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathDockerImageInitContainersPod)...)
		case artifact.TypeKubernetesConfigMap:
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapVolume)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapPodVolume)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapProjectedVolume)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapPodProjectedVolume)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapKeyValueContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapKeyValueContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapKeyValueInitContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapKeyValueInitContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapEnvContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapEnvContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapEnvInitContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathConfigMapEnvInitContainersPod)...)
		case artifact.TypeKubernetesSecret:
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretVolume)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretVolumePod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretProjectedVolume)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretProjectedVolumePod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretKeyValueContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretKeyValueContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretKeyValueInitContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretKeyValueInitContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretEnvContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretEnvContainersPod)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretEnvInitContainers)...)
			bindArtifact(u.Object, a, splitIterables(jsonPathSecretEnvInitContainersPod)...)
		case artifact.TypeKubernetesDeployment:
			if !strings.EqualFold(u.GetKind(), "horizontalPodAutoscaler") {
				continue
			}

			kind, found, _ := unstructured.NestedString(u.Object, fields(jsonPathHPAKind)...)
			if !found || !strings.EqualFold(kind, "deployment") {
				continue
			}

			bindArtifact(u.Object, a, jsonPathHPAName)
		case artifact.TypeKubernetesReplicaSet:
			if !strings.EqualFold(u.GetKind(), "horizontalPodAutoscaler") {
				continue
			}

			kind, found, _ := unstructured.NestedString(u.Object, fields(jsonPathHPAKind)...)
			if !found || !strings.EqualFold(kind, "replicaSet") {
				continue
			}

			bindArtifact(u.Object, a, jsonPathHPAName)
		}
	}
}

// bindArtifact is a recursive function that iterates through all slices found in given
// JSON path, then attempts to find and replace a nested string field, which is the final
// path passeed in.
//
// Take for example the following variadic arguments for paths passed into the function:
// [
//   '.spec.containers',
//   '.env',
//   '.field.name'
// ]
//
// In this case the first two args ('.spec.containers' and '.env') are of type []interface{}.
// We recursively loop through all iterations of each container's environment looking for
// the nested string' .field.name' and if found, replace with the given artifact reference if their names
// match.
func bindArtifact(obj map[string]interface{}, a clouddriver.Artifact, paths ...string) {
	if len(paths) > 1 {
		objs, found, err := unstructured.NestedSlice(obj, fields(paths[0])...)
		if !found || err != nil {
			return
		}

		for _, obj := range objs {
			o, ok := obj.(map[string]interface{})
			if !ok {
				continue
			}

			bindArtifact(o, a, paths[1:]...)
		}

		unstructured.SetNestedField(obj, objs, fields(paths[0])...)
	} else {
		name, found, _ := unstructured.NestedString(obj, fields(paths[0])...)
		if found && a.Name == name {
			_ = unstructured.SetNestedField(obj, a.Reference, fields(paths[0])...)
		}
	}
}

// splitIterables splits a string on the character '*' returning the resulting
// string slice.
func splitIterables(path string) []string {
	return strings.Split(path, "*")
}

// fields takes a path, removes any '.' prefix and suffix and return the resulting
// string slice.
func fields(path string) []string {
	path = strings.TrimPrefix(path, ".")
	path = strings.TrimSuffix(path, ".")

	return strings.Split(path, ".")
}
