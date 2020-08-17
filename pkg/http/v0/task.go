package v0

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var bearerToken string

func init() {
	bearerToken = os.Getenv("BEARER_TOKEN")
	if bearerToken == "" {
		panic("bearer token not set")
	}
}

// Get a task - currently only associated with kubernetes 'tasks'.
func GetTask(c *gin.Context) {
	sc := sql.Instance(c)
	id := c.Param("id")
	manifests := []map[string]interface{}{}

	resources, err := sc.ListKubernetesResources(id)
	if err != nil {
		e := clouddriver.NewError(
			"BadRequest",
			"Error getting task: "+err.Error(),
			http.StatusBadRequest,
		)
		c.JSON(http.StatusBadRequest, e)

		return
	}

	var accountName string
	// TODO create a separate table to associate a provider with a task ID.
	if len(resources) > 0 {
		accountName = resources[0].AccountName
	}

	provider, err := sc.GetKubernetesProvider(accountName)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error getting provider: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error decoding provider CA data: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: os.Getenv("BEARER_TOKEN"),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error generating dynamic kubernetes client: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	for _, r := range resources {
		resource := schema.GroupVersionResource{
			Group:    r.Group,
			Version:  r.Version,
			Resource: r.Resource,
		}
		result, err := client.
			Resource(resource).
			Namespace(r.Namespace).
			Get(context.TODO(), r.Name, metav1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		manifests = append(manifests, result.Object)
	}

	ro := clouddriver.ResultObject{
		Manifests:                         manifests,
		CreatedArtifacts:                  buildCreatedArtifacts(resources),
		ManifestNamesByNamespace:          makeManifestNamesByNamespace(resources),
		ManifestNamesByNamespaceToRefresh: makeManifestNamesByNamespaceToRefresh(resources),
	}

	tr := clouddriver.TaskResponse{
		ID:            id,
		ResultObjects: []clouddriver.ResultObject{ro},
		Status: clouddriver.TaskStatus{
			Complete:  true,
			Completed: true,
			Failed:    false,
			Phase:     "ORCHESTRATION",
			Retryable: false,
			Status:    "Orchestration completed.",
		},
	}

	c.JSON(http.StatusOK, tr)
}

func buildCreatedArtifacts(resources []kubernetes.Resource) []clouddriver.CreatedArtifact {
	cas := []clouddriver.CreatedArtifact{}

	for _, resource := range resources {
		ca := clouddriver.CreatedArtifact{
			CustomKind: false,
			Location:   resource.Namespace,
			Metadata: clouddriver.CreatedArtifactMetadata{
				Account: resource.AccountName,
			},
			Name:      resource.Name,
			Reference: resource.Name,
			Type:      "kubernetes/" + resource.Kind,
			Version:   resource.Version,
		}
		cas = append(cas, ca)
	}

	return cas
}

func makeManifestNamesByNamespace(resources []kubernetes.Resource) map[string][]string {
	m := map[string][]string{}

	for _, resource := range resources {
		if _, ok := m[resource.Namespace]; !ok {
			m[resource.Namespace] = []string{}
		}

		a := m[resource.Namespace]
		a = append(a, fmt.Sprintf("%s %s", resource.Kind, resource.Name))
		m[resource.Namespace] = a
	}

	return m
}

func makeManifestNamesByNamespaceToRefresh(resources []kubernetes.Resource) map[string][]string {
	m := map[string][]string{}

	for _, resource := range resources {
		if _, ok := m[resource.Namespace]; !ok {
			m[resource.Namespace] = []string{}
		}

		a := m[resource.Namespace]
		a = append(a, fmt.Sprintf("%s %s", resource.Kind, resource.Name))
		m[resource.Namespace] = a
	}

	return m
}
