package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
)

// CreateKubernetesOperation is the main function that starts a kubernetes operation.
//
// Kubernetes operations are things like deploy/delete manifest or perform
// a rolling restart. Spinnaker sends *all* of these types of events to the
// same endpoint (/kubernetes/ops), so we have to unmarshal and check which
// kind of operation we are performing.
//
// The actual actions have been moved to the kubernetes subfolder to make
// this function a bit more readable.
func CreateKubernetesOperation(c *gin.Context) {
	// All operations are bound to a task ID and stored in the database.
	ko := kubernetes.Operations{}
	taskID := clouddriver.TaskIDFromContext(c)

	err := c.ShouldBindJSON(&ko)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	// Loop through each request in the kubernetes operations and perform
	// each requested action.
	for _, req := range ko {
		if req.DeployManifest != nil {
			kubernetes.Deploy(c, *req.DeployManifest)
		}

		if req.DeleteManifest != nil {
			kubernetes.Delete(c, *req.DeleteManifest)
		}

		if req.ScaleManifest != nil {
			kubernetes.Scale(c, *req.ScaleManifest)
		}

		if req.CleanupArtifacts != nil {
			kubernetes.CleanupArtifacts(c, *req.CleanupArtifacts)
		}

		if req.RollingRestartManifest != nil {
			kubernetes.RollingRestart(c, *req.RollingRestartManifest)
		}

		if req.RunJob != nil {
			kubernetes.RunJob(c, *req.RunJob)
		}

		if req.UndoRolloutManifest != nil {
			kubernetes.Rollback(c, *req.UndoRolloutManifest)
		}

		if req.PatchManifest != nil {
			kubernetes.Patch(c, *req.PatchManifest)
		}

		if c.Errors != nil {
			return
		}
	}

	or := kubernetes.OperationsResponse{
		ID:          taskID,
		ResourceURI: "/task/" + taskID,
	}
	c.JSON(http.StatusOK, or)
}
