package kubernetes

import (
	"encoding/base64"
	"encoding/json"
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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func Patch(c *gin.Context, pm PatchManifestRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := sc.GetKubernetesProvider(pm.Account)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	token, err := ac.Token()
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
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
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(pm.PatchBody)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
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
		clouddriver.WriteError(c, http.StatusBadRequest,
			fmt.Errorf("invalid merge strategy %s", pm.Options.MergeStrategy))
		return
	}

	meta, _, err := client.PatchUsingStrategy(kind, name, pm.Location, b, strategy)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// TODO Record the applied patch in the kubernetes.io/change-cause annotation. If the annotation already exists, the contents are replaced.
	if pm.Options.Record {
	}

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

	err = sc.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}
}
