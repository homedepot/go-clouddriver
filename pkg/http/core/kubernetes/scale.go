package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)

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

	a := strings.Split(sm.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := client.Get(kind, name, sm.Location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	var meta kube.Metadata

	// TODO need to allow scaling for additional kinds.
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

		meta, err = client.Apply(&scaledManifestObject)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	case "statefulset":
		ss := kubernetes.NewStatefulSet(u.Object)

		replicas, err := strconv.Atoi(sm.Replicas)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		desiredReplicas := int32(replicas)

		ss.SetReplicas(&desiredReplicas)

		scaledManifestObject, err := ss.ToUnstructured()
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		meta, err = client.Apply(&scaledManifestObject)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	default:
		clouddriver.Error(c, http.StatusBadRequest,
			fmt.Errorf("scaling kind %s not currently supported", kind))
		return
	}

	kr := kubernetes.Resource{
		AccountName:  sm.Account,
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
