package v0

import (
	"encoding/base64"
	"fmt"
	"net/http"

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

	resources, err := sc.ListKubernetesResourcesByTaskID(id)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	// If there were no kubernetes resources associated with this task ID,
	// return the default task.
	if len(resources) == 0 {
		c.JSON(http.StatusOK, clouddriver.NewDefaultTask(id))
		return
	}

	accountName := resources[0].AccountName

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
		BearerToken: provider.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	if err = kc.WithConfig(config); err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	for _, r := range resources {
		result, err := kc.Get(r.Resource, r.Name, r.Namespace)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

		manifests = append(manifests, result.Object)
	}

	ro := clouddriver.TaskResultObject{
		Manifests:                         manifests,
		CreatedArtifacts:                  buildCreatedArtifacts(resources),
		ManifestNamesByNamespace:          makeManifestNamesByNamespace(resources),
		ManifestNamesByNamespaceToRefresh: makeManifestNamesByNamespaceToRefresh(resources),
	}

	task := clouddriver.NewDefaultTask(id)
	task.ResultObjects = []clouddriver.TaskResultObject{ro}

	c.JSON(http.StatusOK, task)
}

func buildCreatedArtifacts(resources []kubernetes.Resource) []clouddriver.TaskCreatedArtifact {
	cas := []clouddriver.TaskCreatedArtifact{}

	for _, resource := range resources {
		ca := clouddriver.TaskCreatedArtifact{
			CustomKind: false,
			Location:   resource.Namespace,
			Metadata: clouddriver.TaskCreatedArtifactMetadata{
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
