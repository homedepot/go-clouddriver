package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

// RollingRestart performs a `kubectl rollout restart` by setting an annotation on a pod template
// to the current time in RFC3339.
func RollingRestart(c *gin.Context, rr RollingRestartManifestRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := sc.GetKubernetesProvider(rr.Account)
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

	a := strings.Split(rr.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := client.Get(kind, name, rr.Location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	var meta kube.Metadata

	switch strings.ToLower(kind) {
	case "deployment":
		// Add annotation to pod spec:
		// kubectl.kubernetes.io/restartedAt: "2020-08-21T03:56:27Z"
		kubernetes.AnnotateTemplate(u, "clouddriver.spinnaker.io/restartedAt",
			time.Now().In(time.UTC).Format(time.RFC3339))

		meta, err = client.Apply(u)
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

	err = sc.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
}
