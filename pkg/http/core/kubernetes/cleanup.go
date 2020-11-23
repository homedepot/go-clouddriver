package kubernetes

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

func CleanupArtifacts(c *gin.Context, ca CleanupArtifactsRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)

	for _, manifest := range ca.Manifests {
		u, err := kc.ToUnstructured(manifest)
		if err != nil {
			clouddriver.WriteError(c, http.StatusBadRequest, err)
			return
		}

		provider, err := sc.GetKubernetesProvider(ca.Account)
		if err != nil {
			clouddriver.WriteError(c, http.StatusBadRequest, err)
			return
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			clouddriver.WriteError(c, http.StatusBadRequest, err)
			return
		}

		token, err := ac.Token()
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
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
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

		gvr, err := client.GVRForKind(u.GetKind())
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

		kr := kubernetes.Resource{
			AccountName:  ca.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			TaskType:     "cleanup",
			APIGroup:     gvr.Group,
			Name:         u.GetName(),
			Namespace:    u.GetNamespace(),
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         u.GetKind(),
			SpinnakerApp: app,
			Cluster:      cluster(u.GetKind(), u.GetName()),
		}

		err = sc.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}
	}
}
