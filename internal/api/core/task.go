package core

import (
	"fmt"
	"net/http"
	"strings"

	clouddriver "github.com/homedepot/go-clouddriver/pkg"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/iancoleman/strcase"
)

// GetTask gets a task - currently only associated with kubernetes 'tasks'.
func (cc *Controller) GetTask(c *gin.Context) {
	id := c.Param("id")
	manifests := []map[string]interface{}{}

	resources, err := cc.SQLClient.ListKubernetesResourcesByTaskID(id)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	if len(resources) == 0 {
		clouddriver.Error(c, http.StatusNotFound, fmt.Errorf("Task not found (id: %s)", id))
		return
	}

	accountName := resources[0].AccountName

	provider, err := cc.KubernetesProvider(accountName)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	for _, r := range resources {
		// Ignore getting the manifest if task type is "cleanup" or "noop".
		if strings.EqualFold(r.TaskType, clouddriver.TaskTypeCleanup) ||
			strings.EqualFold(r.TaskType, clouddriver.TaskTypeNoOp) {
			manifests = append(manifests, map[string]interface{}{})
			continue
		}

		result, err := provider.Client.Get(r.Resource, r.Name, r.Namespace)
		if err != nil {
			// If the task type is "delete" and the resource was not found,
			// append an empty manifest and continue.
			// I tried to use `errors.IsNotFound(err)` here to check
			// if the error was a not found error, but was unable to get the
			// test to work in doing this, so for now we are just checking
			// the string suffix.
			if strings.EqualFold(r.TaskType, clouddriver.TaskTypeDelete) &&
				strings.HasSuffix(err.Error(), "not found") {
				manifests = append(manifests, map[string]interface{}{})
				continue
			}

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

func buildCreatedArtifacts(resources []kubernetes.Resource) []clouddriver.Artifact {
	var (
		artifactVersion string
		lastIndex       int
	)

	cas := []clouddriver.Artifact{}

	for _, resource := range resources {
		artifactVersion = ""
		lastIndex = strings.LastIndex(resource.Name, "-v")

		if lastIndex != -1 {
			artifactVersion = resource.Name[lastIndex+1:]
		}

		ca := clouddriver.Artifact{
			CustomKind: false,
			Location:   resource.Namespace,
			Metadata: clouddriver.ArtifactMetadata{
				Account: resource.AccountName,
			},
			Name:      resource.ArtifactName,
			Reference: resource.Name,
			Type:      artifact.Type("kubernetes/" + strcase.ToLowerCamel(resource.Kind)),
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
