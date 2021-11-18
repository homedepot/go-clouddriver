package kubernetes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

// RollingRestart performs a `kubectl rollout restart` by setting an annotation on a pod template
// to the current time in RFC3339.
func (cc *Controller) RollingRestart(c *gin.Context, rr RollingRestartManifestRequest) {
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := rr.Location

	provider, err := cc.KubernetesProvider(rr.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	if provider.Namespace != nil {
		namespace = *provider.Namespace
	}

	a := strings.Split(rr.ManifestName, " ")
	kind := a[0]
	name := a[1]

	err = provider.ValidateKindStatus(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	u, err := provider.Client.Get(kind, name, namespace)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	var meta kubernetes.Metadata

	switch strings.ToLower(kind) {
	case "deployment":
		// Add annotation to pod spec:
		// kubectl.kubernetes.io/restartedAt: "2020-08-21T03:56:27Z"
		err = kubernetes.AnnotateTemplate(u, "clouddriver.spinnaker.io/restartedAt",
			time.Now().In(time.UTC).Format(time.RFC3339))
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		meta, err = provider.Client.Apply(u)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	default:
		clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("restarting kind %s not currently supported", kind))
		return
	}

	kr := kubernetes.Resource{
		AccountName:  rr.Account,
		ID:           uuid.New().String(),
		TaskID:       taskID,
		APIGroup:     meta.Group,
		Name:         meta.Name,
		Namespace:    meta.Namespace,
		Resource:     meta.Resource,
		Version:      meta.Version,
		Kind:         meta.Kind,
		SpinnakerApp: app,
	}

	err = cc.SQLClient.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
}
