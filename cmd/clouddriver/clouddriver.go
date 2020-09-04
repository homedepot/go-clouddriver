package main

import (
	"log"
	"os"

	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/helm"
	kube "github.com/billiford/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/server"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
)

var (
	r = gin.New()
)

func main() {
	r.Run(":7002")
}

func init() {
	// Setup metrics.
	p := ginprometheus.NewPrometheus("gin")
	p.MetricsPath = "/metrics"
	p.Use(r)

	gin.ForceConsoleColor()
	// Ignore logging of certain endpoints.
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health"}}))
	r.Use(gin.Recovery())

	kubeController := kubernetes.NewController()
	sqlClient := sql.NewClient(mustDBConnect())
	helmClient := helm.NewClient("https://kubernetes-charts.storage.googleapis.com")
	arcadeClient := arcade.NewDefaultClient()

	arcadeAPIKey := mustGetenv("ARCADE_API_KEY")
	arcadeClient.WithAPIKey(arcadeAPIKey)

	c := &server.Config{
		ArcadeClient:      arcadeClient,
		SQLClient:         sqlClient,
		KubeController:    kubeController,
		KubeActionHandler: kube.NewActionHandler(),
		HelmClient:        helmClient,
	}
	server.Setup(r, c)
}

func mustDBConnect() *gorm.DB {
	sqlConfig := sql.Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
	}

	db, err := sql.Connect(sql.Connection(sqlConfig))
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func mustGetenv(env string) (s string) {
	if s = os.Getenv(env); s == "" {
		log.Fatal(env + " not set; exiting.")
	}

	return
}
