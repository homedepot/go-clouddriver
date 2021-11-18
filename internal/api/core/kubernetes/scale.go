package kubernetes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (cc *Controller) Scale(c *gin.Context, sm ScaleManifestRequest) {
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := sm.Location

	provider, err := cc.KubernetesProvider(sm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	if provider.Namespace != nil {
		namespace = *provider.Namespace
	}

	a := strings.Split(sm.ManifestName, " ")
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

		meta, err = provider.Client.Apply(u)
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
