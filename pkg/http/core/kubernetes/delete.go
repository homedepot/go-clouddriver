package kubernetes

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewDeleteManifestAction(ac ActionConfig) Action {
	return &deleteManfest{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		dm: ac.Operation.DeleteManifest,
	}
}

type deleteManfest struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	dm *DeleteManifestRequest
}

func (d *deleteManfest) Run() error {
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

	do := metav1.DeleteOptions{}

	if d.dm.Options.GracePeriodSeconds != nil {
		do.GracePeriodSeconds = d.dm.Options.GracePeriodSeconds
	}

	propagationPolicy := v1.DeletePropagationOrphan
	if d.dm.Options.Cascading {
		propagationPolicy = v1.DeletePropagationForeground
	}

	do.PropagationPolicy = &propagationPolicy

	// Default to the static mode.
	mode := "static"
	if d.dm.Mode != "" {
		mode = d.dm.Mode
	}

	switch strings.ToLower(mode) {
	// Both dynamic and static use the same logic. For 'dynamic' the manifest has already been resolved and passed in.
	case "dynamic", "static":
		a := strings.Split(d.dm.ManifestName, " ")
		kind := a[0]
		name := a[1]

		gvr, err := client.GVRForKind(kind)
		if err != nil {
			return err
		}

		err = client.DeleteResourceByKindAndNameAndNamespace(kind, name, d.dm.Location, do)
		if err != nil {
			return err
		}

		kr := kubernetes.Resource{
			AccountName:  d.dm.Account,
			ID:           uuid.New().String(),
			TaskID:       d.id,
			APIGroup:     gvr.Group,
			Name:         name,
			Namespace:    d.dm.Location,
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         kind,
			SpinnakerApp: d.dm.App,
			Cluster:      cluster(kind, name),
		}

		err = d.sc.CreateKubernetesResource(kr)
		if err != nil {
			return err
		}
	case "label":
		return fmt.Errorf("requested to delete manifest %s using mode %s which is not implemented", d.dm.ManifestName, mode)
	default:
		return fmt.Errorf("requested to delete manifest %s using mode %s which is not implemented", d.dm.ManifestName, mode)
	}

	return nil
}
