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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// GetTask gets a task - currently only associated with kubernetes 'tasks'.
func (cc *Controller) GetTask(c *gin.Context) {
	id := c.Param("id")
	task := clouddriver.NewDefaultTask(id)
	manifests := []map[string]interface{}{}

	resources, err := cc.SQLClient.ListKubernetesResourcesByTaskID(id)
	if err != nil {
		clouddriver.Log(fmt.Errorf("GetTask: error listing resources for task (id: %s): %w", id, err))

		task.Status.Failed = true
		task.Status.Retryable = true
		task.Status.Status = fmt.Sprintf("Error listing resources for task (id: %s): %v", id, err)
		c.JSON(http.StatusInternalServerError, task)

		return
	}

	if len(resources) == 0 {
		task.Status.Failed = true
		task.Status.Status = fmt.Sprintf("Task not found (id: %s)", id)
		c.JSON(http.StatusNotFound, task)

		return
	}

	accountName := resources[0].AccountName

	provider, err := cc.KubernetesProvider(accountName)
	if err != nil {
		clouddriver.Log(fmt.Errorf("GetTask: error getting kubernetes provider for task (id: %s, account: %s): %w",
			id, accountName, err))

		task.Status.Failed = true
		task.Status.Retryable = true
		task.Status.Status = fmt.Sprintf("Error getting kubernetes provider %s for task (id: %s): %v",
			accountName, id, err)
		c.JSON(http.StatusInternalServerError, task)

		return
	}

	// Rows whose resource vanished after a successful create and is allowed to
	// vanish (see mayVanishAfterCreate). They are left out of the manifests
	// and of the name maps below.
	vanished := map[int]bool{}

	for i, r := range resources {
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
			if strings.EqualFold(r.TaskType, clouddriver.TaskTypeDelete) &&
				strings.HasSuffix(err.Error(), "not found") {
				manifests = append(manifests, map[string]interface{}{})

				continue
			}

			// A row only exists after the create succeeded, so a "not found"
			// here means the resource vanished afterwards. That is expected
			// for a few kinds that something else removes within seconds
			// (the Argo Rollouts canary Mapping, a Job removed by its TTL).
			// Every other kind keeps failing the task.
			if apierrors.IsNotFound(err) && mayVanishAfterCreate(r, resources) &&
				lookupGroupMatchesStored(provider.Client, r) {
				clouddriver.Log(fmt.Errorf("GetTask: tolerated vanished resource (id: %s, account: %s, "+
					"group: %s, resource: %s, name: %s, namespace: %s): %w",
					id, r.AccountName, r.APIGroup, r.Resource, r.Name, r.Namespace, err))

				// Leave it out of the manifests: Orca's cleanupArtifacts task
				// sends them back and needs a kind for every entry.
				vanished[i] = true

				continue
			}

			clouddriver.Log(fmt.Errorf("GetTask: error getting resource for task (id: %s, account: %s, task type: %s, "+
				"group: %s, version: %s, resource: %s, name: %s, namespace: %s): %w",
				id, r.AccountName, r.TaskType, r.APIGroup, r.Version, r.Resource, r.Name, r.Namespace, err))

			task.Status.Failed = true
			task.Status.Retryable = true
			task.Status.Status = fmt.Sprintf("Error getting resource for task (task ID: %s, kind: %s, name: %s, namespace: %s): %v",
				id, r.Resource, r.Name, r.Namespace, err)
			c.JSON(http.StatusInternalServerError, task)

			return
		} else if strings.EqualFold(r.TaskType, clouddriver.TaskTypeDelete) {
			task.Status.Complete = false
			task.Status.Completed = false
			task.Status.Status = "Orchestration in progress."
		}

		manifests = append(manifests, result.Object)
	}

	// The names Orca checks for stability leave out the resources that
	// vanished and the canary Mappings that Argo Rollouts manages, so
	// waitForManifestToStabilize never asks about them.
	mnr := buildMapOfNamespaceToResource(namesForStabilityCheck(id, resources, vanished))
	// Refactor bound artifact to get the list of bound artifacts as not all created artifacts need to be bound.
	createdArtifacts := buildCreatedArtifacts(resources)
	ro := clouddriver.TaskResultObject{
		BoundArtifacts:                    createdArtifacts,
		DeployedNamesByLocation:           mnr,
		CreatedArtifacts:                  createdArtifacts,
		Manifests:                         manifests,
		ManifestNamesByNamespace:          mnr,
		ManifestNamesByNamespaceToRefresh: mnr,
	}

	task.ResultObjects = []clouddriver.TaskResultObject{ro}

	c.JSON(http.StatusOK, task)
}

const (
	ambassadorAPIGroup = "getambassador.io"
	canarySuffix       = "-canary"
)

// isMapping returns true if the resource is an Ambassador (Emissary) Mapping.
func isMapping(r kubernetes.Resource) bool {
	return strings.EqualFold(r.Kind, "Mapping") && strings.EqualFold(r.APIGroup, ambassadorAPIGroup)
}

// isPairedCanaryMapping returns true if the resource is a Mapping named
// "<base>-canary" and the same task also has a Mapping named "<base>" in the
// same namespace. That pair is what Argo Rollouts creates: it clones the base
// Mapping into "<base>-canary" and deletes the clone when the rollout ends.
// A "-canary" Mapping without a base Mapping is owned by the application and
// is not paired.
func isPairedCanaryMapping(r kubernetes.Resource, resources []kubernetes.Resource) bool {
	if !isMapping(r) || !strings.HasSuffix(r.Name, canarySuffix) || len(r.Name) == len(canarySuffix) {
		return false
	}

	base := strings.TrimSuffix(r.Name, canarySuffix)

	for _, other := range resources {
		if isMapping(other) && other.Namespace == r.Namespace && other.Name == base {
			return true
		}
	}

	return false
}

// mayVanishAfterCreate returns true for the kinds that are allowed to be gone
// when the task is read back after a successful create: a paired canary
// Mapping and a Job (removed by its TTL). Everything else, such as workloads
// or custom kinds owned by other teams, must still exist.
func mayVanishAfterCreate(r kubernetes.Resource, resources []kubernetes.Resource) bool {
	if isPairedCanaryMapping(r, resources) {
		return true
	}

	return strings.EqualFold(r.Kind, "Job") && strings.EqualFold(r.APIGroup, "batch")
}

// lookupGroupMatchesStored returns true if the group the lookup by resource
// name resolves to is the group the resource was stored with. When a plural
// exists in more than one group (for example gateways), the lookup can pick
// the wrong group and report "not found" for a resource that exists.
func lookupGroupMatchesStored(client kubernetes.Client, r kubernetes.Resource) bool {
	gvk, err := client.GVKForKind(r.Resource)
	if err != nil {
		return false
	}

	return strings.EqualFold(gvk.Group, r.APIGroup)
}

// namesForStabilityCheck returns the rows that make up the name maps, without
// the rows that vanished and the paired canary Mappings.
func namesForStabilityCheck(id string, resources []kubernetes.Resource, vanished map[int]bool) []kubernetes.Resource {
	rows := make([]kubernetes.Resource, 0, len(resources))

	for i, r := range resources {
		if vanished[i] {
			continue
		}

		if isPairedCanaryMapping(r, resources) {
			clouddriver.Log(fmt.Errorf("GetTask: skipped canary Mapping (id: %s, account: %s, "+
				"group: %s, resource: %s, name: %s, namespace: %s)",
				id, r.AccountName, r.APIGroup, r.Resource, r.Name, r.Namespace))

			continue
		}

		rows = append(rows, r)
	}

	return rows
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
