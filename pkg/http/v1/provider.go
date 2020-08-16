package v1

import (
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

func CreateKubernetesProvider(c *gin.Context) {
	sc := sql.Instance(c)
	p := kubernetes.Provider{}

	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sc.GetKubernetesProvider(p.Name)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "provider already exists"})
		return
	}

	err = sc.CreateKubernetesProvider(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}
