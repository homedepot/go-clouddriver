package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	kube "github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

func (cc *Controller) Scale(c *gin.Context, sm ScaleManifestRequest) {
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := sm.Location

	provider, err := cc.SQLClient.GetKubernetesProvider(sm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	if provider.Namespace != nil {
		namespace = *provider.Namespace
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	token, err := cc.ArcadeClient.Token(provider.TokenProvider)
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

	client, err := cc.KubernetesController.NewClient(config)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	a := strings.Split(sm.ManifestName, " ")
	kind := a[0]
	name := a[1]

	err = provider.ValidateKindStatus(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	u, err := client.Get(kind, name, namespace)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	var meta kube.Metadata

	switch strings.ToLower(kind) {
	case "deployment", "replicaset", "statefulset":
		r, err := strconv.Atoi(sm.Replicas)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		err = unstructured.SetNestedField(u.Object, int64(r), "spec", "replicas")
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		meta, err = client.ApplyWithNamespaceOverride(u, namespace)
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

	err = cc.SQLClient.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
}
