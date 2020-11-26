package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/artifact"
	"github.com/homedepot/go-clouddriver/pkg/fiat"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
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

func SetFiatClient(f fiat.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(fiat.ClientInstanceKey, f)
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
