package middleware

import (
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/gin-gonic/gin"
)

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		// If an error occured during handling the request, write the error as a JSON response.
		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			statusCode := c.Writer.Status()
			ce := clouddriver.NewError(
				http.StatusText(statusCode),
				err.Error(),
				statusCode,
			)
			c.JSON(c.Writer.Status(), ce)
		}
	}
}
