package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Vary sets the 'Vary' header to the defined headers.
func Vary(headers ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Vary", strings.Join(headers, ", "))
		c.Next()
	}
}
