package kubernetes

import (
	"encoding/base64"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8s.io/client-go/rest"
)

func DeployManifests(c *gin.Context, taskID string, dm DeployManifestRequest) error {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)

	provider, err := sc.GetKubernetesProvider(dm.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: provider.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	kc.WithConfig(config)

	for _, manifest := range dm.Manifests {
		u, err := kubernetes.ToUnstructured(manifest)
		if err != nil {
			return err
		}

		err = kubernetes.AddSpinnakerAnnotations(u, dm.Moniker.App)
		if err != nil {
			return err
		}

		err = kubernetes.AddSpinnakerLabels(u, dm.Moniker.App)
		if err != nil {
			return err
		}

		meta, err := kc.Apply(u)
		if err != nil {
			return err
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			APIGroup:     meta.Group,
			Name:         meta.Name,
			Namespace:    meta.Namespace,
			Resource:     meta.Resource,
			Version:      meta.Version,
			Kind:         meta.Kind,
			SpinnakerApp: dm.Moniker.App,
		}

		err = sc.CreateKubernetesResource(kr)
		if err != nil {
			return err
		}
	}

	return nil
}
