package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CacheControl sets the 'Cache-Control' header to the defined
// max-age.
func CacheControl(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Next()
	}
}
