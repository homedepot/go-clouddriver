package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		// If an error occured during handling the request, write the error as a JSON response.
		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			statusCode := c.Writer.Status()
			text := http.StatusText(statusCode)

			if statusCode >= http.StatusInternalServerError {
				meta := clouddriver.Meta(err)
				clouddriver.Log(err, meta)
				text += " (error ID: " + meta.GUID + ")"
			}

			ce := clouddriver.NewError(
				text,
				err.Error(),
				statusCode,
			)

			c.JSON(c.Writer.Status(), ce)
		}
	}
}
