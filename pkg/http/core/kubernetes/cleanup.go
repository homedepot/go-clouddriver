package kubernetes

import (
	"encoding/base64"

	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/google/uuid"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewCleanupArtifactsAction(ac ActionConfig) Action {
	return &cleanupArtifacts{
		ac:  ac.ArcadeClient,
		sc:  ac.SQLClient,
		kc:  ac.KubeController,
		id:  ac.ID,
		ca:  ac.Operation.CleanupArtifacts,
		app: ac.Application,
	}
}

type cleanupArtifacts struct {
	ac  arcade.Client
	sc  sql.Client
	kc  kubernetes.Controller
	id  string
	ca  *CleanupArtifactsRequest
	app string
}

func (c *cleanupArtifacts) Run() error {
	for _, manifest := range c.ca.Manifests {
		u, err := c.kc.ToUnstructured(manifest)
		if err != nil {
			return err
		}

		provider, err := c.sc.GetKubernetesProvider(c.ca.Account)
		if err != nil {
			return err
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			return err
		}

		token, err := c.ac.Token()
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

		client, err := c.kc.NewClient(config)
		if err != nil {
			return err
		}

		gvr, err := client.GVRForKind(u.GetKind())
		if err != nil {
			return err
		}

		kr := kubernetes.Resource{
			AccountName:  c.ca.Account,
			ID:           uuid.New().String(),
			TaskID:       c.id,
			APIGroup:     gvr.Group,
			Name:         u.GetName(),
			Namespace:    u.GetNamespace(),
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         u.GetKind(),
			SpinnakerApp: c.app,
			Cluster:      cluster(u.GetKind(), u.GetName()),
		}

		err = c.sc.CreateKubernetesResource(kr)
		if err != nil {
			return err
		}
	}

	return nil
}
