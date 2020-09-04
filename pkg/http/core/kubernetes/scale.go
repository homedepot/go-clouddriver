package kubernetes

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewScaleManifestAction(ac ActionConfig) Action {
	return &scaleManifest{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		sm: ac.Operation.ScaleManifest,
	}
}

type scaleManifest struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	sm *ScaleManifestRequest
}

func (s *scaleManifest) Run() error {
	provider, err := s.sc.GetKubernetesProvider(s.sm.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	token, err := s.ac.Token()
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

	client, err := s.kc.NewClient(config)
	if err != nil {
		return err
	}

	a := strings.Split(s.sm.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := client.Get(kind, name, s.sm.Location)
	if err != nil {
		return err
	}

	// TODO need to allow scaling for other kinds.
	switch strings.ToLower(kind) {
	case "deployment":
		d := kubernetes.NewDeployment(u.Object)

		replicas, err := strconv.Atoi(s.sm.Replicas)
		if err != nil {
			return err
		}

		desiredReplicas := int32(replicas)
		d.SetReplicas(&desiredReplicas)

		scaledManifestObject, err := d.ToUnstructured()
		if err != nil {
			return err
		}

		_, err = client.Apply(&scaledManifestObject)
		if err != nil {
			return err
		}
	}

	return nil
}
