package core

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/client-go/rest"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/iancoleman/strcase"
)

// Get a task - currently only associated with kubernetes 'tasks'.
func GetTask(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	id := c.Param("id")
	manifests := []map[string]interface{}{}

	resources, err := sc.ListKubernetesResourcesByTaskID(id)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
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
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
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

	for _, r := range resources {
		// Ignore getting the manifest if task type is "cleanup".
		if strings.EqualFold(r.TaskType, "cleanup") {
			manifests = append(manifests, map[string]interface{}{})
			continue
		}

		result, err := client.Get(r.Resource, r.Name, r.Namespace)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		manifests = append(manifests, result.Object)
	}

	mnr := buildMapOfNamespaceToResource(resources)

	//Refactor bound artifact to get the list of bound artifacts as not all created artifacts need to be bound
	createdArtifacts := buildCreatedArtifacts(resources)

	ro := clouddriver.TaskResultObject{
		BoundArtifacts:                    createdArtifacts,
		DeployedNamesByLocation:           mnr,
		CreatedArtifacts:                  createdArtifacts,
		Manifests:                         manifests,
		ManifestNamesByNamespace:          mnr,
		ManifestNamesByNamespaceToRefresh: mnr,
	}

	task := clouddriver.NewDefaultTask(id)
	task.ResultObjects = []clouddriver.TaskResultObject{ro}

	c.JSON(http.StatusOK, task)
}

func buildCreatedArtifacts(resources []kubernetes.Resource) []clouddriver.TaskCreatedArtifact {
	var (
		artifactVersion string
		lastIndex       int
	)

	cas := []clouddriver.TaskCreatedArtifact{}

	for _, resource := range resources {
		artifactVersion = ""
		lastIndex = strings.LastIndex(resource.Name, "-v")

		if lastIndex != -1 {
			artifactVersion = resource.Name[lastIndex+1:]
		}

		ca := clouddriver.TaskCreatedArtifact{
			CustomKind: false,
			Location:   resource.Namespace,
			Metadata: clouddriver.TaskCreatedArtifactMetadata{
				Account: resource.AccountName,
			},
			Name:      resource.ArtifactName,
			Reference: resource.Name,
			Type:      "kubernetes/" + strcase.ToLowerCamel(resource.Kind),
			Version:   artifactVersion,
		}
		cas = append(cas, ca)
	}

	return cas
}

func buildMapOfNamespaceToResource(resources []kubernetes.Resource) map[string][]string {
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
