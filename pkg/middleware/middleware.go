package middleware

import (
	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/artifact"
	"github.com/billiford/go-clouddriver/pkg/helm"
	kube "github.com/billiford/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

func SetArcadeClient(a arcade.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(arcade.ClientInstanceKey, a)
		c.Next()
	}
}

func SetSQLClient(r sql.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(sql.ClientInstanceKey, r)
		c.Next()
	}
}

func SetKubeController(k kubernetes.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(kubernetes.ControllerInstanceKey, k)
		c.Next()
	}
}

func SetArtifactCredentialsController(a artifact.CredentialsController) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(artifact.CredentialsControllerInstanceKey, a)
		c.Next()
	}
}

func SetKubeActionHandler(k kube.ActionHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(kube.ActionHandlerInstanceKey, k)
		c.Next()
	}
}

func SetHelmClient(h helm.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(helm.ClientInstanceKey, h)
		c.Next()
	}
}
