package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

	token, err := ac.Token()
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

	switch strings.ToLower(kind) {
	case "deployment":
		d := kubernetes.NewDeployment(u.Object)

		// add annotations to pod spec
		// kubectl.kubernetes.io/restartedAt: "2020-08-21T03:56:27Z"
		d.AnnotateTemplate("clouddriver.spinnaker.io/restartedAt",
			time.Now().In(time.UTC).Format(time.RFC3339))

		annotatedObject, err := d.ToUnstructured()
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		_, err = client.Apply(&annotatedObject)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	default:
		clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("restarting kind %s not currently supported", kind))
		return
	}
}
