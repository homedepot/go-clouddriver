package v1

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/jinzhu/gorm"
)

// CreateKubernetesProvider creates the kubernetes account (provider).
func (cc *Controller) CreateKubernetesProvider(c *gin.Context) {
	p := kubernetes.Provider{}

	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = cc.SQLClient.GetKubernetesProvider(p.Name)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "provider already exists"})
		return
	}

	err = cc.validate(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if p.Namespace != nil && strings.TrimSpace(*p.Namespace) == "" {
		p.Namespace = nil
	}

	err = cc.SQLClient.CreateKubernetesProvider(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// DeleteKubernetesProvider deletes the kubernetes account (provider).
func (cc *Controller) DeleteKubernetesProvider(c *gin.Context) {
	name := c.Param("name")

	_, err := cc.SQLClient.GetKubernetesProvider(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = cc.SQLClient.DeleteKubernetesProvider(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetKubernetesProvider retrieves the kubernetes account (provider).
func (cc *Controller) GetKubernetesProvider(c *gin.Context) {
	name := c.Param("name")

	p, err := cc.SQLClient.GetKubernetesProviderAndPermissions(name)
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

// ListKubernetesProvider retrieves all the kubernetes accounts (providers).
func (cc *Controller) ListKubernetesProvider(c *gin.Context) {
	providers, err := cc.SQLClient.ListKubernetesProvidersAndPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, providers)
}

// CreateOrReplaceKubernetesProvider creates the kubernetes account (provider),
// or if existing account, replaces it.
func (cc *Controller) CreateOrReplaceKubernetesProvider(c *gin.Context) {
	p := kubernetes.Provider{}

	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.validate(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if p.Namespace != nil && strings.TrimSpace(*p.Namespace) == "" {
		p.Namespace = nil
	}

	err = cc.SQLClient.DeleteKubernetesProvider(p.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = cc.SQLClient.CreateKubernetesProvider(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

// validates verifies the providers data.  Validations performed:
// - the CAData is base64 encoded
// - the TokenProvider is known/supported by arcade
// - every Permissions.Write entry exists in Permissions.Read
func (cc *Controller) validate(p kubernetes.Provider) error {
	_, err := base64.StdEncoding.DecodeString(p.CAData)
	if err != nil {
		return fmt.Errorf("error decoding base64 CA data: %s", err.Error())
	}

	_, err = cc.ArcadeClient.Token(p.TokenProvider)
	if err != nil {
		return fmt.Errorf("error getting token: %s", err.Error())
	}

	// Verify that each write group is included as a read group
	for _, wg := range p.Permissions.Write {
		found := false

		for _, rg := range p.Permissions.Read {
			if strings.EqualFold(wg, rg) {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("error in permissions: write group '%s' must be included as a read group", wg)
		}
	}

	return nil
}
