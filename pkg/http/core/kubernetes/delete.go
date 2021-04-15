package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/homedepot/go-clouddriver/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func Delete(c *gin.Context, dm DeleteManifestRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := sc.GetKubernetesProvider(dm.Account)
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

	do := metav1.DeleteOptions{}

	if dm.Options.GracePeriodSeconds != nil {
		do.GracePeriodSeconds = dm.Options.GracePeriodSeconds
	}

	propagationPolicy := v1.DeletePropagationOrphan
	if dm.Options.Cascading {
		propagationPolicy = v1.DeletePropagationForeground
	}

	do.PropagationPolicy = &propagationPolicy

	// Default to the static mode.
	mode := "static"
	if dm.Mode != "" {
		mode = dm.Mode
	}

	switch strings.ToLower(mode) {
	// Both dynamic and static use the same logic. For 'dynamic' the manifest has already been resolved and passed in.
	case "dynamic", "static":
		a := strings.Split(dm.ManifestName, " ")
		kind := a[0]
		name := a[1]

		gvr, err := client.GVRForKind(kind)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = client.DeleteResourceByKindAndNameAndNamespace(kind, name, dm.Location, do)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			TaskType:     clouddriver.TaskTypeDelete,
			Timestamp:    util.CurrentTimeUTC(),
			APIGroup:     gvr.Group,
			Name:         name,
			Namespace:    dm.Location,
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         kind,
			SpinnakerApp: dm.App,
			Cluster:      cluster(kind, name),
		}

		err = sc.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	case "label":
		clouddriver.Error(c, http.StatusNotImplemented,
			fmt.Errorf("requested to delete manifest %s using mode %s which is not implemented", dm.ManifestName, mode))
		return
	default:
		clouddriver.Error(c, http.StatusNotImplemented,
			fmt.Errorf("requested to delete manifest %s using mode %s which is not implemented", dm.ManifestName, mode))
		return
	}
}
