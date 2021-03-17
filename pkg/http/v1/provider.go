package v1

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/jinzhu/gorm"
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

	_, err = base64.StdEncoding.DecodeString(p.CAData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error decoding base64 CA data: %s", err.Error())})
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

func DeleteKubernetesProvider(c *gin.Context) {
	sc := sql.Instance(c)
	name := c.Param("name")

	_, err := sc.GetKubernetesProvider(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = sc.DeleteKubernetesProvider(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
