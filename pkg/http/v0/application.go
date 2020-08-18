package v0

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ServerGroupManagersResponse []ServerGroupManager

type ServerGroupManager struct {
	Account       string                    `json:"account"`
	AccountName   string                    `json:"accountName"`
	CloudProvider string                    `json:"cloudProvider"`
	CreatedTime   int64                     `json:"createdTime"`
	Key           ServerGroupManagerKey     `json:"key"`
	Kind          string                    `json:"kind"`
	Labels        map[string]string         `json:"labels"`
	Manifest      map[string]interface{}    `json:"manifest"`
	Moniker       ServerGroupManagerMoniker `json:"moniker"`
	Name          string                    `json:"name"`
	ProviderType  string                    `json:"providerType"`
	Region        string                    `json:"region"`
	ServerGroups  []ServerGroup             `json:"serverGroups"`
	Type          string                    `json:"type"`
	UID           string                    `json:"uid"`
	Zone          string                    `json:"zone"`
}

type ServerGroupManagerKey struct {
	Account        string `json:"account"`
	Group          string `json:"group"`
	KubernetesKind string `json:"kubernetesKind"`
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Provider       string `json:"provider"`
}

type ServerGroupManagerMoniker struct {
	App     string `json:"app"`
	Cluster string `json:"cluster"`
}

type ServerGroup struct {
	Account string `json:"account"`
	Moniker struct {
		App      string `json:"app"`
		Cluster  string `json:"cluster"`
		Sequence int    `json:"sequence"`
	} `json:"moniker"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Region    string `json:"region"`
}

// Server Group Managers for a kubernetes target are deployments.
func ListServerGroupManagers(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)
	spinnakerApp := c.Param("application")
	response := ServerGroupManagersResponse{}

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(spinnakerApp)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Don't actually return while attempting to create a list of server group managers.
	// We want to avoid the situation where a user cannot perform operations when any
	// cluster is not available.
	for _, account := range accounts {
		provider, err := sc.GetKubernetesProvider(account)
		if err != nil {
			log.Println("unable to get kubernetes provider for account", account)
			continue
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			log.Println("error decoding ca data for account", account)
			continue
		}

		config := &rest.Config{
			Host:        provider.Host,
			BearerToken: os.Getenv("BEARER_TOKEN"),
			TLSClientConfig: rest.TLSClientConfig{
				CAData: cd,
			},
		}

		if err = kc.WithConfig(config); err != nil {
			log.Println("error creating dynamic client for account", account)
			continue
		}

		// Create a GVR for deployments.
		gvr := schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "deployments",
		}
		gvk := schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    "deployment",
		}
		lo := metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=" + spinnakerApp,
		}

		deployments, err := kc.List(gvr, lo)
		if err != nil {
			log.Println("error listing deployments:", err.Error())
			continue
		}

		for _, deployment := range deployments.Items {
			sgr := ServerGroupManager{
				Account:       account,
				AccountName:   account,
				CloudProvider: "kubernetes",
				CreatedTime:   deployment.GetCreationTimestamp().Unix() * 1000,
				Key: ServerGroupManagerKey{
					Account:        account,
					Group:          gvr.Group,
					KubernetesKind: gvk.Kind,
					Name:           deployment.GetName(),
					Namespace:      deployment.GetNamespace(),
					Provider:       "kubernetes",
				},
				Kind:     gvk.Kind,
				Labels:   deployment.GetLabels(),
				Manifest: deployment.Object,
				Moniker: ServerGroupManagerMoniker{
					App:     spinnakerApp,
					Cluster: fmt.Sprintf("%s %s", gvk.Kind, deployment.GetName()),
				},
				Name:         fmt.Sprintf("%s %s", gvk.Kind, deployment.GetName()),
				ProviderType: "kubernetes",
				Region:       spinnakerApp,
				ServerGroups: nil,
				Type:         "kubernetes",
				UID:          string(deployment.GetUID()),
				Zone:         spinnakerApp,
			}
			response = append(response, sgr)
		}
	}

	c.JSON(http.StatusOK, response)
}
