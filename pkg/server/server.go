package server

import (
	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/helm"
	"github.com/billiford/go-clouddriver/pkg/http"
	kube "github.com/billiford/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/middleware"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

type Config struct {
	ArcadeClient      arcade.Client
	SQLClient         sql.Client
	KubeController    kubernetes.Controller
	KubeActionHandler kube.ActionHandler
	HelmClient        helm.Client
}

// Define all middlewares to use then set up the API.
func Setup(r *gin.Engine, c *Config) {
	r.Use(middleware.SetArcadeClient(c.ArcadeClient))
	r.Use(middleware.SetSQLClient(c.SQLClient))
	r.Use(middleware.SetKubeController(c.KubeController))
	r.Use(middleware.SetKubeActionHandler(c.KubeActionHandler))
	r.Use(middleware.SetHelmClient(c.HelmClient))
	r.Use(middleware.HandleError())

	http.Initialize(r)
}
