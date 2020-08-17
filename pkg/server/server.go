package server

import (
	"github.com/billiford/go-clouddriver/pkg/http"
	"github.com/billiford/go-clouddriver/pkg/middleware"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

type Config struct {
	SQLClient sql.Client
}

// Define all middlewares to use then set up the API.
func Setup(r *gin.Engine, c *Config) {
	r.Use(middleware.SetSQLClient(c.SQLClient))

	http.Initialize(r)
}
