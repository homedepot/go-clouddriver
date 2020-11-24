package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

func Scale(c *gin.Context, sm ScaleManifestRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)

	provider, err := sc.GetKubernetesProvider(sm.Account)
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

	a := strings.Split(sm.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := client.Get(kind, name, sm.Location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// TODO need to allow scaling for other kinds.
	switch strings.ToLower(kind) {
	case "deployment":
		d := kubernetes.NewDeployment(u.Object)

		replicas, err := strconv.Atoi(sm.Replicas)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		desiredReplicas := int32(replicas)
		d.SetReplicas(&desiredReplicas)

		scaledManifestObject, err := d.ToUnstructured()
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		_, err = client.Apply(&scaledManifestObject)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	default:
		clouddriver.Error(c, http.StatusBadRequest,
			fmt.Errorf("scaling kind %s not currently supported", kind))
		return
	}
}
