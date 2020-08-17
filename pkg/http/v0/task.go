package v0

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"k8s.io/client-go/rest"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

// Get a task - currently only associated with kubernetes 'tasks'.
func GetTask(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)
	id := c.Param("id")
	manifests := []map[string]interface{}{}

	resources, err := sc.ListKubernetesResources(id)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	var accountName string
	// TODO create a separate table to associate a provider with a task ID.
	if len(resources) > 0 {
		accountName = resources[0].AccountName
	}

	provider, err := sc.GetKubernetesProvider(accountName)
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

	for _, r := range resources {
		result, err := kc.Get(r.Kind, r.Name, r.Namespace)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
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
