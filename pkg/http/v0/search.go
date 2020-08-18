package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, []string{})
}
