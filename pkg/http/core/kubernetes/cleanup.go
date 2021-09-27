package kubernetes

import (
	"encoding/base64"
	"net/http"

	"github.com/homedepot/go-clouddriver/pkg/util"

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
		u, err := kube.ToUnstructured(manifest)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		provider, err := sc.GetKubernetesProvider(ca.Account)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		token, err := ac.Token(provider.TokenProvider)
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

		gvr, err := client.GVRForKind(u.GetKind())
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		namespace := u.GetNamespace()
		if provider.Namespace != "" {
			namespace = provider.Namespace
		}

		kr := kubernetes.Resource{
			AccountName:  ca.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			TaskType:     clouddriver.TaskTypeCleanup,
			Timestamp:    util.CurrentTimeUTC(),
			APIGroup:     gvr.Group,
			Name:         u.GetName(),
			Namespace:    namespace,
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         u.GetKind(),
			SpinnakerApp: app,
			Cluster:      cluster(u.GetKind(), u.GetName()),
		}

		err = sc.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	}
}
