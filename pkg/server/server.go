package server

import (
	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/artifact"
	"github.com/homedepot/go-clouddriver/pkg/fiat"
	"github.com/homedepot/go-clouddriver/pkg/http"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/middleware"
	"github.com/homedepot/go-clouddriver/pkg/sql"
)

type Config struct {
	ArcadeClient                  arcade.Client
	ArtifactCredentialsController artifact.CredentialsController
	SQLClient                     sql.Client
	FiatClient                    fiat.Client
	KubeController                kubernetes.Controller
	VerboseRequestLogging         bool
}

// Define all middlewares to use then set up the API.
func Setup(r *gin.Engine, c *Config) {
	r.Use(middleware.SetArcadeClient(c.ArcadeClient))
	r.Use(middleware.SetSQLClient(c.SQLClient))
	r.Use(middleware.SetKubeController(c.KubeController))
	r.Use(middleware.SetArtifactCredentialsController(c.ArtifactCredentialsController))
	r.Use(middleware.SetFiatClient(c.FiatClient))
	r.Use(middleware.HandleError())

	if c.VerboseRequestLogging {
		r.Use(middleware.LogRequest())
	}

	http.Initialize(r)
}
