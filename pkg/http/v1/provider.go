package v1

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/jinzhu/gorm"
)

// CreateKubernetesProvider creates the kubernetes account (provider).
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

	c.JSON(http.StatusCreated, p)
}

// DeleteKubernetesProvider deletes the kubernetes account (provider).
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

// GetKubernetesProvider retrieves the kubernetes account (provider).
func GetKubernetesProvider(c *gin.Context) {
	var p kubernetes.Provider

	sc := sql.Instance(c)
	name := c.Param("name")

	p, err := sc.GetKubernetesProviderAndPermissions(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	if p.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// CreateOrReplaceKubernetesProvider creates the kubernetes account (provider),
// or if existing account, replaces it.
func CreateOrReplaceKubernetesProvider(c *gin.Context) {
	sc := sql.Instance(c)
	p := kubernetes.Provider{}

	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = base64.StdEncoding.DecodeString(p.CAData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error decoding base64 CA data: %s", err.Error())})
		return
	}

	err = sc.DeleteKubernetesProvider(p.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sc.CreateKubernetesProvider(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}
