package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

// TaskID attaches a unique GUID to the context.
func TaskID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(clouddriver.TaskIDKey, uuid.New().String())
		c.Next()
	}
}
