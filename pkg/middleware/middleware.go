package middleware

import (
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

func SetSQLClient(r sql.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(sql.ClientInstanceKey, r)
		c.Next()
	}
}
