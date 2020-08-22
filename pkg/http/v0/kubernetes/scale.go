package kubernetes

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
)

func ScaleManifest(c *gin.Context, sm ScaleManifestRequest) error {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)

	provider, err := sc.GetKubernetesProvider(sm.Account)
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

	if err = kc.SetDynamicClientForConfig(config); err != nil {
		return err
	}

	a := strings.Split(sm.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := kc.Get(kind, name, sm.Location)
	if err != nil {
		return err
	}

	switch strings.ToLower(kind) {
	case "deployment":
		d := kubernetes.NewDeployment(u.Object)

		replicas, err := strconv.Atoi(sm.Replicas)
		if err != nil {
			return err
		}

		desiredReplicas := int32(replicas)

		d.SetReplicas(&desiredReplicas)

		scaledManifestObject, err := d.ToUnstructured()
		if err != nil {
			return err
		}

		_, err = kc.Apply(&scaledManifestObject)
		if err != nil {
			return err
		}
	}

	return nil
}
