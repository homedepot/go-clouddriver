package middleware

import (
	kube "github.com/billiford/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

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

func SetKubeActionHandler(k kube.ActionHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(kube.ActionHandlerInstanceKey, k)
		c.Next()
	}
}
