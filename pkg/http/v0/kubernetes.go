package v0

import (
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/http/v0/kubernetes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateKubernetesOperation(c *gin.Context) {
	taskID := uuid.New().String()
	ko := kubernetes.Operations{}

	err := c.ShouldBindJSON(&ko)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if len(ko) == 0 {
		or := kubernetes.OperationsResponse{
			ID:          taskID,
			ResourceURI: "/task/" + taskID,
		}
		c.JSON(http.StatusOK, or)
		return
	}

	for _, req := range ko {
		if req.DeployManifest != nil {
			err = kubernetes.DeployManifests(c, taskID, *req.DeployManifest)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}
		}
		if req.ScaleManifest != nil {
			err = kubernetes.ScaleManifest(c, *req.ScaleManifest)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	or := kubernetes.OperationsResponse{
		ID:          taskID,
		ResourceURI: "/task/" + taskID,
	}
	c.JSON(http.StatusOK, or)
}
