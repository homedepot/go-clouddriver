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

func SetKubeClient(k kubernetes.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(kubernetes.ClientInstanceKey, k)
		c.Next()
	}
}

func SetKubeActionHandler(k kube.ActionHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(kube.ActionHandlerInstanceKey, k)
		c.Next()
	}
}
