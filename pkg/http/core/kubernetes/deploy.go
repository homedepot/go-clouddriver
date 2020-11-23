package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/google/uuid"
<<<<<<< HEAD
	"k8s.io/apimachinery/pkg/util/rand"
=======
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
>>>>>>> go-clouddriver-1/spinnaker-versioning
	"k8s.io/client-go/rest"
)

var (
	listTimeout = int64(30)
)

func (ah *actionHandler) NewDeployManifestAction(ac ActionConfig) Action {
	return &deployManifest{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		dm: ac.Operation.DeployManifest,
	}
}

type deployManifest struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	dm *DeployManifestRequest
}

func (d *deployManifest) Run() error {
	provider, err := d.sc.GetKubernetesProvider(d.dm.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	token, err := d.ac.Token()
	if err != nil {
		return err
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := d.kc.NewClient(config)
	if err != nil {
		return err
	}

	manifests := []map[string]interface{}{}

	// Merge all list element items into the manifest list.
	for _, manifest := range d.dm.Manifests {
		u, err := d.kc.ToUnstructured(manifest)
		if err != nil {
			return err
		}
		application := d.dm.Moniker.App

<<<<<<< HEAD
		if strings.EqualFold(u.GetKind(), "list") {
			listElement := kubernetes.ListElement{}

			b, err := json.Marshal(u.Object)
			if err != nil {
				return err
			}

			err = json.Unmarshal(b, &listElement)
			if err != nil {
				return err
			}

			manifests = append(manifests, listElement.Items...)
		} else {
			manifests = append(manifests, u.Object)
		}
	}

	for _, manifest := range manifests {
		u, err := d.kc.ToUnstructured(manifest)
		if err != nil {
			return err
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

		err = d.kc.AddSpinnakerAnnotations(u, d.dm.Moniker.App)
=======
		err = d.kc.AddSpinnakerAnnotations(u, application)
>>>>>>> go-clouddriver-1/spinnaker-versioning
		if err != nil {
			return err
		}

		err = d.kc.AddSpinnakerLabels(u, application)
		if err != nil {
			return err
		}

		name := u.GetName()

		if d.kc.IsVersioned(u) {
			kind := strings.ToLower(u.GetKind())
			namespace := u.GetNamespace()
			lo := metav1.ListOptions{
				LabelSelector:  kubernetes.LabelKubernetesName + "=" + application,
				TimeoutSeconds: &listTimeout,
			}

			results, err := client.ListResourcesByKindAndNamespace(kind, namespace, lo)
			if err != nil {
				return err
			}

			currentVersion := d.kc.GetCurrentVersion(results, kind, name)
			latestVersion := d.kc.IncrementVersion(currentVersion)

			u.SetName(u.GetName() + "-" + latestVersion.Long)

			err = d.kc.AddSpinnakerVersionAnnotations(u, latestVersion)
			if err != nil {
				return err
			}

			err = d.kc.AddSpinnakerVersionLabels(u, latestVersion)
			if err != nil {
				return err
			}
		}

		meta, err := client.ApplyWithNamespaceOverride(u, d.dm.NamespaceOverride)
		if err != nil {
			return fmt.Errorf("error applying manifest (kind: %s, apiVersion: %s, name: %s): %s",
				u.GetKind(), u.GroupVersionKind().Version, u.GetName(), err.Error())
		}

		kr := kubernetes.Resource{
			AccountName:  d.dm.Account,
			ID:           uuid.New().String(),
			TaskID:       d.id,
			APIGroup:     meta.Group,
			Name:         meta.Name,
			Namespace:    meta.Namespace,
			Resource:     meta.Resource,
			Version:      meta.Version,
			Kind:         meta.Kind,
			SpinnakerApp: d.dm.Moniker.App,
			Cluster:      cluster(meta.Kind, name),
		}

		err = d.sc.CreateKubernetesResource(kr)
		if err != nil {
			return err
		}
	}

	return nil
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
