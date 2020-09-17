package kubernetes

import (
	"encoding/base64"

	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/google/uuid"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewDeployManifestAction(ac ActionConfig) Action {
	return &deployManfest{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		dm: ac.Operation.DeployManifest,
	}
}

type deployManfest struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	dm *DeployManifestRequest
}

func (d *deployManfest) Run() error {
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

	for _, manifest := range d.dm.Manifests {
		u, err := d.kc.ToUnstructured(manifest)
		if err != nil {
			return err
		}

		err = d.kc.AddSpinnakerAnnotations(u, d.dm.Moniker.App)
		if err != nil {
			return err
		}

		err = d.kc.AddSpinnakerLabels(u, d.dm.Moniker.App)
		if err != nil {
			return err
		}

		meta, err := client.ApplyWithNamespaceOverride(u, d.dm.NamespaceOverride)
		if err != nil {
			return err
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
		}

		err = d.sc.CreateKubernetesResource(kr)
		if err != nil {
			return err
		}
	}

	return nil
}
