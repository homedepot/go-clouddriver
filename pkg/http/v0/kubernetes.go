package v0

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/deployment"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8s.io/client-go/rest"
)

type OpsRequest []struct {
	DeployManifest *DeployManifest `json:"deployManifest"`
	ScaleManifest  *ScaleManifest  `json:"scaleManifest"`
}

type DeployManifest struct {
	EnableTraffic     bool                     `json:"enableTraffic"`
	NamespaceOverride string                   `json:"namespaceOverride"`
	OptionalArtifacts []interface{}            `json:"optionalArtifacts"`
	CloudProvider     string                   `json:"cloudProvider"`
	Manifests         []map[string]interface{} `json:"manifests"`
	TrafficManagement struct {
		Options struct {
			EnableTraffic bool `json:"enableTraffic"`
		} `json:"options"`
		Enabled bool `json:"enabled"`
	} `json:"trafficManagement"`
	Moniker struct {
		App string `json:"app"`
	} `json:"moniker"`
	Source                   string        `json:"source"`
	Account                  string        `json:"account"`
	SkipExpressionEvaluation bool          `json:"skipExpressionEvaluation"`
	RequiredArtifacts        []interface{} `json:"requiredArtifacts"`
}

type ScaleManifest struct {
	Replicas      string `json:"replicas"`
	ManifestName  string `json:"manifestName"`
	CloudProvider string `json:"cloudProvider"`
	Location      string `json:"location"`
	User          string `json:"user"`
	Account       string `json:"account"`
}

type OpsResponse struct {
	ID          string `json:"id"`
	ResourceURI string `json:"resourceUri"`
}

type ManifestResponse struct {
	Account string `json:"account"`
	// Artifacts []struct {
	// 	CustomKind bool `json:"customKind"`
	// 	Metadata   struct {
	// 	} `json:"metadata"`
	// 	Name      string `json:"name"`
	// 	Reference string `json:"reference"`
	// 	Type      string `json:"type"`
	// } `json:"artifacts"`
	Events   []interface{}           `json:"events"`
	Location string                  `json:"location"`
	Manifest map[string]interface{}  `json:"manifest"`
	Metrics  []interface{}           `json:"metrics"`
	Moniker  ManifestResponseMoniker `json:"moniker"`
	Name     string                  `json:"name"`
	Status   manifest.Status         `json:"status"`
	Warnings []interface{}           `json:"warnings"`
}

type ManifestResponseMoniker struct {
	App     string `json:"app"`
	Cluster string `json:"cluster"`
}

func CreateKubernetesDeployment(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)
	taskID := uuid.New().String()
	kor := OpsRequest{}

	err := c.ShouldBindJSON(&kor)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if len(kor) == 0 {
		or := OpsResponse{
			ID:          taskID,
			ResourceURI: "/task/" + taskID,
		}
		c.JSON(http.StatusOK, or)
		return
	}

	for _, req := range kor {
		if req.DeployManifest != nil {
			// TODO this is hacky - need to figure out how to handle providers.
			accountName := kor[0].DeployManifest.Account

			provider, err := sc.GetKubernetesProvider(kor[0].DeployManifest.Account)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			cd, err := base64.StdEncoding.DecodeString(provider.CAData)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			config := &rest.Config{
				Host:        provider.Host,
				BearerToken: os.Getenv("BEARER_TOKEN"),
				TLSClientConfig: rest.TLSClientConfig{
					CAData: cd,
				},
			}

			if err = kc.WithConfig(config); err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			for _, manifest := range req.DeployManifest.Manifests {
				b, err := json.Marshal(manifest)
				if err != nil {
					clouddriver.WriteError(c, http.StatusBadRequest, err)
					return
				}

				_, meta, err := kc.Apply(b, req.DeployManifest.Moniker.App)
				if err != nil {
					clouddriver.WriteError(c, http.StatusInternalServerError, err)
					return
				}

				kr := kubernetes.Resource{
					AccountName:  accountName,
					ID:           uuid.New().String(),
					TaskID:       taskID,
					APIGroup:     meta.Group,
					Name:         meta.Name,
					Namespace:    meta.Namespace,
					Resource:     meta.Resource,
					Version:      meta.Version,
					Kind:         meta.Kind,
					SpinnakerApp: req.DeployManifest.Moniker.App,
				}

				err = sc.CreateKubernetesResource(kr)
				if err != nil {
					clouddriver.WriteError(c, http.StatusInternalServerError, err)
					return
				}
			}
		}
		if req.ScaleManifest != nil {
			sm := req.ScaleManifest
			provider, err := sc.GetKubernetesProvider(sm.Account)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			cd, err := base64.StdEncoding.DecodeString(provider.CAData)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			config := &rest.Config{
				Host:        provider.Host,
				BearerToken: os.Getenv("BEARER_TOKEN"),
				TLSClientConfig: rest.TLSClientConfig{
					CAData: cd,
				},
			}

			if err = kc.WithConfig(config); err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			a := strings.Split(sm.ManifestName, " ")
			kind := a[0]
			name := a[1]

			unstructuredObj, err := kc.Get(kind, name, sm.Location)
			if err = kc.WithConfig(config); err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}

			switch strings.ToLower(kind) {
			case "deployment":
				d := deployment.New(unstructuredObj.Object)
				replicas, _ := strconv.Atoi(sm.Replicas)
				desiredReplicas := int32(replicas)

				d.Spec.Replicas = &desiredReplicas

				b, err := json.Marshal(d)
				if err != nil {
					log.Println("error marshaling")
					clouddriver.WriteError(c, http.StatusBadRequest, err)
					return
				}

				_, err = kc.Patch(b)
				if err != nil {
					clouddriver.WriteError(c, http.StatusInternalServerError, err)
					return
				}
			}
		}
	}

	or := OpsResponse{
		ID:          taskID,
		ResourceURI: "/task/" + taskID,
	}
	c.JSON(http.StatusOK, or)
}

func GetManifest(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)
	account := c.Param("account")
	namespace := c.Param("location")
	n := c.Param("name")
	a := strings.Split(n, " ")
	kind := a[0]
	name := a[1]

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: os.Getenv("BEARER_TOKEN"),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	if err = kc.WithConfig(config); err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	result, err := kc.Get(kind, name, namespace)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	app := "unknown"
	labels := result.GetLabels()
	if _, ok := labels["app.kubernetes.io/name"]; ok {
		app = labels["app.kubernetes.io/name"]
	}

	kmr := ManifestResponse{
		Account:  account,
		Events:   []interface{}{},
		Location: namespace,
		Manifest: result.Object,
		Metrics:  []interface{}{},
		Moniker: ManifestResponseMoniker{
			App:     app,
			Cluster: fmt.Sprintf("%s %s", kind, name),
		},
		Name: fmt.Sprintf("%s %s", kind, name),
		// The 'default' status of a kubernetes resource.
		Status:   kubernetes.GetStatus(kind, result.Object),
		Warnings: []interface{}{},
	}

	c.JSON(http.StatusOK, kmr)
}
