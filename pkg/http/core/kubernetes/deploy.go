package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/homedepot/go-clouddriver/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"

	"k8s.io/apimachinery/pkg/util/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	listTimeout = int64(30)
)

func Deploy(c *gin.Context, dm DeployManifestRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := sc.GetKubernetesProvider(dm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	token, err := ac.Token()
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	manifests := []map[string]interface{}{}
	application := dm.Moniker.App
	// Merge all list element items into the manifest list.
	for _, manifest := range dm.Manifests {
		u, err := kc.ToUnstructured(manifest)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		if strings.EqualFold(u.GetKind(), "list") {
			listElement := kubernetes.ListElement{}

			b, err := json.Marshal(u.Object)
			if err != nil {
				clouddriver.Error(c, http.StatusBadRequest, err)
				return
			}

			err = json.Unmarshal(b, &listElement)
			if err != nil {
				clouddriver.Error(c, http.StatusBadRequest, err)
				return
			}

			manifests = append(manifests, listElement.Items...)
		} else {
			manifests = append(manifests, u.Object)
		}
	}

	// Sort the manifests by their kind's priority.
	manifests, err = kc.SortManifests(manifests)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	for _, manifest := range manifests {
		u, err := kc.ToUnstructured(manifest)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		// If the kind is a job, its name is not set, and generateName is set,
		// generate a name for the job as `apply` will throw the error
		// `resource name may not be empty`.

		if strings.EqualFold(u.GetKind(), "job") {
			name := u.GetName()

			generateName := u.GetGenerateName()

			if name == "" && generateName != "" {
				u.SetName(generateName + rand.String(randNameNumber))
			}
		}

		name := u.GetName()
		artifactName := name

		err = kc.AddSpinnakerAnnotations(u, application)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = kc.AddSpinnakerLabels(u, application)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = kc.VersionVolumes(u, dm.RequiredArtifacts)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		if kc.IsVersioned(u) {
			kind := strings.ToLower(u.GetKind())
			namespace := u.GetNamespace()
			labelSelector := metav1.LabelSelector{
				MatchLabels: map[string]string{
					kubernetes.LabelKubernetesName: application,
				},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      kubernetes.LabelSpinnakerMonikerSequence,
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			}

			ls, err := metav1.LabelSelectorAsSelector(&labelSelector)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			lo := metav1.ListOptions{
				LabelSelector:  ls.String(),
				TimeoutSeconds: &listTimeout,
			}

			results, err := client.ListResourcesByKindAndNamespace(kind, namespace, lo)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			currentVersion := kc.GetCurrentVersion(results, kind, name)
			latestVersion := kc.IncrementVersion(currentVersion)
			u.SetName(u.GetName() + "-" + latestVersion.Long)

			err = kc.AddSpinnakerVersionAnnotations(u, latestVersion)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			err = kc.AddSpinnakerVersionLabels(u, latestVersion)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		meta, err := client.ApplyWithNamespaceOverride(u, dm.NamespaceOverride)
		if err != nil {
			e := fmt.Errorf("error applying manifest (kind: %s, apiVersion: %s, name: %s): %s",
				u.GetKind(), u.GroupVersionKind().Version, u.GetName(), err.Error())
			clouddriver.Error(c, http.StatusInternalServerError, e)

			return
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			Timestamp:    util.CurrentTimeUTC(),
			APIGroup:     meta.Group,
			Name:         meta.Name,
			ArtifactName: artifactName,
			Namespace:    meta.Namespace,
			Resource:     meta.Resource,
			Version:      meta.Version,
			Kind:         meta.Kind,
			SpinnakerApp: dm.Moniker.App,
			Cluster:      cluster(meta.Kind, name),
		}

		err = sc.CreateKubernetesResource(kr)
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

func lowercaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}
