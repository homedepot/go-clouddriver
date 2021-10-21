package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
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
func (cc *Controller) CreateKubernetesOperation(c *gin.Context) {
	// All operations are bound to a task ID and stored in the database.
	ko := kubernetes.Operations{}
	taskID := clouddriver.TaskIDFromContext(c)

	if err := c.ShouldBindBodyWith(&ko, binding.JSON); err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	kc := kubernetes.Controller{
		Controller: cc.Controller,
	}
	// Loop through each request in the kubernetes operations and perform
	// each requested action.
	for _, req := range ko {
		if req.DeployManifest != nil {
			kc.Deploy(c, *req.DeployManifest)
		}

		if req.DeleteManifest != nil {
			kc.Delete(c, *req.DeleteManifest)
		}

		if req.ScaleManifest != nil {
			kc.Scale(c, *req.ScaleManifest)
		}

		if req.CleanupArtifacts != nil {
			kc.CleanupArtifacts(c, *req.CleanupArtifacts)
		}

		if req.RollingRestartManifest != nil {
			kc.RollingRestart(c, *req.RollingRestartManifest)
		}

		if req.RunJob != nil {
			kc.RunJob(c, *req.RunJob)
		}

		if req.UndoRolloutManifest != nil {
			kc.Rollback(c, *req.UndoRolloutManifest)
		}

		if req.PatchManifest != nil {
			kc.Patch(c, *req.PatchManifest)
		}

		if c.Errors != nil && len(c.Errors) > 0 {
			return
		}
	}

	or := kubernetes.OperationsResponse{
		ID:          taskID,
		ResourceURI: "/task/" + taskID,
	}
	c.JSON(http.StatusOK, or)
}
