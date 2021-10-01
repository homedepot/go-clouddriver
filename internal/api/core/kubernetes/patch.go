package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func (cc *Controller) Patch(c *gin.Context, pm PatchManifestRequest) {
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := pm.Location

	provider, err := cc.SQLClient.GetKubernetesProvider(pm.Account)
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

	b, err := json.Marshal(pm.PatchBody)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	// Manifest name is *really* the Spinnaker cluster - i.e. "deployment test-deployment", so we
	// need to split on a whitespace and get the actual name of the manifest.
	kind := ""
	name := pm.ManifestName

	a := strings.Split(pm.ManifestName, " ")
	if len(a) > 1 {
		kind = a[0]
		name = a[1]
	}

	err = provider.ValidateKindStatus(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	// Merge strategy can be "strategic", "json", or "merge".
	var strategy types.PatchType

	switch pm.Options.MergeStrategy {
	case "strategic":
		strategy = types.StrategicMergePatchType
	case "json":
		strategy = types.JSONPatchType
	case "merge":
		strategy = types.MergePatchType
	default:
		clouddriver.Error(c, http.StatusBadRequest,
			fmt.Errorf("invalid merge strategy %s", pm.Options.MergeStrategy))
		return
	}

	meta, _, err := client.PatchUsingStrategy(kind, name, namespace, b, strategy)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// TODO Record the applied patch in the kubernetes.io/change-cause annotation. If the annotation already exists, the contents are replaced.
	// if pm.Options.Record {
	// }

	kr := kubernetes.Resource{
		AccountName:  pm.Account,
		ID:           uuid.New().String(),
		TaskID:       taskID,
		APIGroup:     meta.Group,
		Name:         meta.Name,
		Namespace:    meta.Namespace,
		Resource:     meta.Resource,
		Version:      meta.Version,
		Kind:         meta.Kind,
		SpinnakerApp: pm.App,
	}

	err = cc.SQLClient.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
}
