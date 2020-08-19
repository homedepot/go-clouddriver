package v1

import (
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	for _, group := range p.Permissions.Read {
		rp := clouddriver.ReadPermission{
			ID:          uuid.New().String(),
			AccountName: p.Name,
			ReadGroup:   group,
		}
		err = sc.CreateReadPermission(rp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	for _, group := range p.Permissions.Write {
		wp := clouddriver.WritePermission{
			ID:          uuid.New().String(),
			AccountName: p.Name,
			WriteGroup:  group,
		}
		err = sc.CreateWritePermission(wp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, p)
}
