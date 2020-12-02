package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/artifact"
	"github.com/homedepot/go-clouddriver/pkg/fiat"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/server"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
)

var (
	r = gin.New()
)

func main() {
	err := r.Run(":7002")
	if err != nil {
		panic(err)
	}
}

func init() {
	// Setup metrics.
	p := ginprometheus.NewPrometheus("gin")
	p.MetricsPath = "/metrics"
	p.Use(r)

	// Preserve low cardinality for the request counter.
	// See https://github.com/zsais/go-gin-prometheus#preserving-a-low-cardinality-for-the-request-counter.
	p.ReqCntURLLabelMappingFn = reqCntURLLabelMappingFn

	gin.ForceConsoleColor()
	// Ignore logging of certain endpoints.
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health"}}))
	r.Use(gin.Recovery())

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

	// Grab our artifact credentials from /opt/spinnaker/artifacts/config.
	artifactCredentialsController, err := artifact.NewDefaultCredentialsController()
	if err != nil {
		log.Println("[CLOUDDRIVER] error setting up artifact credentials controller:", err.Error())
	}

	sqlClient := sql.NewClient(db)
	fiatClient := fiat.NewDefaultClient()
	kubeController := kubernetes.NewController()
	arcadeClient := arcade.NewDefaultClient()

	arcadeAPIKey := os.Getenv("ARCADE_API_KEY")
	if arcadeAPIKey == "" {
		log.Println("[CLOUDDRIVER] WARNING: ARCADE_API_KEY not set")
	}

	arcadeClient.WithAPIKey(arcadeAPIKey)

	c := &server.Config{
		ArcadeClient:                  arcadeClient,
		ArtifactCredentialsController: artifactCredentialsController,
		SQLClient:                     sqlClient,
		FiatClient:                    fiatClient,
		KubeController:                kubeController,
	}
	if os.Getenv("VERBOSE_REQUEST_LOGGING") == "true" {
		c.VerboseRequestLogging = true
	}

	server.Setup(r, c)
}

func reqCntURLLabelMappingFn(c *gin.Context) string {
	// Setting the url to the path will remove query params, which sometimes have a GUID.
	url := c.Request.URL.Path

	for _, p := range c.Params {
		// The following replaces certain path params with a generic name.
		switch p.Key {
		case "application":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		case "location":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		case "name":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		case "kind":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		case "cluster":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		case "target":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		case "id":
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		}
	}

	return url
}
