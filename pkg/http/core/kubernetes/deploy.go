package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/homedepot/go-clouddriver/pkg/artifact"
	"github.com/homedepot/go-clouddriver/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

var (
	listTimeout = int64(30)
)

// Deploy performs a "Deploy (Manifest)" Spinnaker operation.
// It takes in a list of manifest and the Kubernetes provider
// to "apply" them to. It adds Spinnaker annotations/labels,
// and handles any Spinnaker versioning, then applies each manifest
// one by one.
func Deploy(c *gin.Context, dm DeployManifestRequest) {
	kubeController := kube.ControllerInstance(c)
	sqlClient := sql.Instance(c)

	config, status, err := kubeClientConfig(c, sqlClient, dm.Account)
	if err != nil {
		clouddriver.Error(c, status, err)
		return
	}

	kubeClient, err := kubeController.NewClient(config)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
	// First, convert all manifests to unstructured objects.
	manifests, err := toUnstructured(dm.Manifests)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}
	// Merge all list element items into the manifest list.
	manifests, err = mergeManifests(manifests)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}
	// Sort the manifests by their kind's priority.
	manifests = kube.SortManifests(manifests)
	// Consolidate all deploy manifest request artifacts.
	artifacts := []clouddriver.Artifact{}
	for _, artifact := range dm.RequiredArtifacts {
		artifacts = append(artifacts, artifact)
	}

	for _, artifact := range dm.OptionalArtifacts {
		artifacts = append(artifacts, artifact)
	}

	application := dm.Moniker.App

	for _, manifest := range manifests {
		nameWithoutVersion := manifest.GetName()
		// If the kind is a job, its name is not set, and generateName is set,
		// generate a name for the job as `apply` will throw the error
		// `resource name may not be empty`.
		if strings.EqualFold(manifest.GetKind(), "job") && nameWithoutVersion == "" {
			generateName := manifest.GetGenerateName()
			manifest.SetName(generateName + rand.String(randNameNumber))
		}

		err = kube.AddSpinnakerAnnotations(&manifest, application)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = kube.AddSpinnakerLabels(&manifest, application)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		kube.BindArtifacts(&manifest, artifacts)

		if kube.IsVersioned(manifest) {
			err := handleVersionedManifest(kubeClient, &manifest, application)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		if kube.UseSourceCapacity(manifest) {
			err = handleUseSourceCapacity(kubeClient, &manifest, dm.NamespaceOverride)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		meta, err := kubeClient.ApplyWithNamespaceOverride(&manifest, dm.NamespaceOverride)
		if err != nil {
			e := fmt.Errorf("error applying manifest (kind: %s, apiVersion: %s, name: %s): %s",
				manifest.GetKind(), manifest.GroupVersionKind().Version, manifest.GetName(), err.Error())
			clouddriver.Error(c, http.StatusInternalServerError, e)

			return
		}

		taskID := clouddriver.TaskIDFromContext(c)
		kr := kube.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			Timestamp:    util.CurrentTimeUTC(),
			APIGroup:     meta.Group,
			Name:         meta.Name,
			ArtifactName: nameWithoutVersion,
			Namespace:    meta.Namespace,
			Resource:     meta.Resource,
			Version:      meta.Version,
			Kind:         meta.Kind,
			SpinnakerApp: dm.Moniker.App,
			Cluster:      cluster(meta.Kind, nameWithoutVersion),
		}

		annotations := manifest.GetAnnotations()
		artifactType := annotations[kube.AnnotationSpinnakerArtifactType]
		artifact := clouddriver.Artifact{
			Name:      nameWithoutVersion,
			Reference: meta.Name,
			Type:      artifact.Type(artifactType),
		}
		artifacts = append(artifacts, artifact)

		err = sqlClient.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	}
}

// Generate the cluster that a kind is a part of.
// A Kubernetes cluster is of kind deployment, statefulSet, replicaSet, ingress, service, and daemonSet
// so only generate a cluster for these kinds.
func cluster(kind, name string) string {
	cluster := ""

	if strings.EqualFold(kind, "deployment") ||
		strings.EqualFold(kind, "statefulSet") ||
		strings.EqualFold(kind, "replicaSet") ||
		strings.EqualFold(kind, "ingress") ||
		strings.EqualFold(kind, "service") ||
		strings.EqualFold(kind, "daemonSet") {
		cluster = fmt.Sprintf("%s %s", lowercaseFirst(kind), name)
	}

	return cluster
}

// toUnstructured converts a slice of map[string]interface{} to unstructured.Unstructured.
func toUnstructured(manifests []map[string]interface{}) ([]unstructured.Unstructured, error) {
	m := []unstructured.Unstructured{}

	for _, manifest := range manifests {
		u, err := kube.ToUnstructured(manifest)
		if err != nil {
			return nil, fmt.Errorf("kubernetes: unable to convert manifest to unstructured: %w", err)
		}

		m = append(m, u)
	}

	return m, nil
}

func lowercaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}

func getListOptions(app string) (metav1.ListOptions, error) {
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			kube.LabelKubernetesName: app,
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      kube.LabelSpinnakerMonikerSequence,
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}

	ls, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		return metav1.ListOptions{}, err
	}

	lo := metav1.ListOptions{
		LabelSelector:  ls.String(),
		TimeoutSeconds: &listTimeout,
	}

	return lo, err
}

// mergeManifests merges manifests of kind List into the parent list of manifets.
func mergeManifests(manifests []unstructured.Unstructured) ([]unstructured.Unstructured, error) {
	mergedManifests := []unstructured.Unstructured{}

	for _, manifest := range manifests {
		if manifest.IsList() {
			ul, err := manifest.ToList()
			if err != nil {
				return nil, fmt.Errorf("error converting manifest to list: %w", err)
			}

			mergedManifests = append(mergedManifests, ul.Items...)
		} else {
			mergedManifests = append(mergedManifests, manifest)
		}
	}

	return mergedManifests, nil
}

func kubeClientConfig(c *gin.Context, sqlClient sql.Client, account string) (*rest.Config, int, error) {
	config := &rest.Config{}

	provider, err := sqlClient.GetKubernetesProvider(account)
	if err != nil {
		return config, http.StatusBadRequest, err
	}

	arcadeClient := arcade.Instance(c)

	token, err := arcadeClient.Token(provider.TokenProvider)
	if err != nil {
		return config, http.StatusInternalServerError, err
	}

	caData, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return config, http.StatusBadRequest, err
	}

	config = &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caData,
		},
	}

	return config, -1, nil
}

func handleVersionedManifest(kubeClient kube.Client, u *unstructured.Unstructured, application string) error {
	lo, err := getListOptions(application)
	if err != nil {
		return err
	}

	kind := strings.ToLower(u.GetKind())
	namespace := u.GetNamespace()

	results, err := kubeClient.ListResourcesByKindAndNamespace(kind, namespace, lo)
	if err != nil {
		return err
	}

	nameWithoutVersion := u.GetName()
	currentVersion := kube.GetCurrentVersion(results, kind, nameWithoutVersion)
	latestVersion := kube.IncrementVersion(currentVersion)
	u.SetName(nameWithoutVersion + "-" + latestVersion.Long)

	err = kube.AddSpinnakerVersionAnnotations(u, latestVersion)
	if err != nil {
		return err
	}

	err = kube.AddSpinnakerVersionLabels(u, latestVersion)
	if err != nil {
		return err
	}

	return nil
}

func handleUseSourceCapacity(kubeClient kube.Client, u *unstructured.Unstructured, namespaceOverride string) error {
	current, err := kubeClient.Get(u.GetKind(), u.GetName(), namespaceOverride)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}
	// If the resource is currently deployed then replace the replicas value
	// with the current value, if it has one
	if current != nil {
		r, found, err := unstructured.NestedInt64(current.Object, "spec", "replicas")
		if err != nil {
			return err
		}

		if found {
			err = unstructured.SetNestedField(u.Object, r, "spec", "replicas")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
