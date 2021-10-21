package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/api"
	"github.com/homedepot/go-clouddriver/internal/arcade"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/fiat"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/sql"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	mysqlDefaultStringSize = 256
)

var (
	r = gin.New()
)

func main() {
	if err := r.Run(":7002"); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Setup metrics.
	p := ginprometheus.NewPrometheus("clouddriver")
	p.MetricsPath = "/metrics"
	p.Use(r)

	// Preserve low cardinality for the request counter.
	// See https://github.com/zsais/go-gin-prometheus#preserving-a-low-cardinality-for-the-request-counter.
	p.ReqCntURLLabelMappingFn = reqCntURLLabelMappingFn

	gin.ForceConsoleColor()
	// Ignore logging of certain endpoints.
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health"}}))
	r.Use(gin.Recovery())

	sqlClient := sql.NewClient(dialector())
	if err := sqlClient.Connect(); err != nil {
		log.Fatal(err)
	}

	artifactCredentialsController := getArtifactsCredentialsController()
	fiatClient := fiat.NewDefaultClient()
	kubeController := kubernetes.NewController()
	arcadeClient := arcade.NewDefaultClient()

	arcadeAPIKey := os.Getenv("ARCADE_API_KEY")
	if arcadeAPIKey == "" {
		log.Println("[CLOUDDRIVER] WARNING: ARCADE_API_KEY not set")
	}

	arcadeClient.WithAPIKey(arcadeAPIKey)

	ic := &internal.Controller{
		ArcadeClient:                  arcadeClient,
		ArtifactCredentialsController: artifactCredentialsController,
		SQLClient:                     sqlClient,
		FiatClient:                    fiatClient,
		KubernetesController:          kubeController,
	}

	server := api.NewServer(r)
	server.WithController(ic)

	if os.Getenv("VERBOSE_REQUEST_LOGGING") == "true" {
		server.WithVerboseRequestLogging()
	}

	server.Setup()
}

func reqCntURLLabelMappingFn(c *gin.Context) string {
	// Setting the url to the path will remove query params, which sometimes have a GUID.
	url := c.Request.URL.Path

	for _, p := range c.Params {
		// The following replaces certain path params with a generic name.
		switch p.Key {
		case "account":
			// Leave account information if this is the Manifests API.
			if !strings.HasPrefix(url, "/manifests") {
				url = strings.Replace(url, p.Value, ":"+p.Key, 1)
			}
		case "application":
			// Leave application information if this is the Applications API.
			if !strings.HasPrefix(url, "/applications") {
				url = strings.Replace(url, p.Value, ":"+p.Key, 1)
			}
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

// getArtifactsCredentialsController gets an artifacts credentials controller
// either using the default directory (/opt/spinnaker/artifacts/config)
// or, if it exists, using the value of the ARTIFACTS_CREDENTIALS_CONFIG_DIR
// environment variable.
func getArtifactsCredentialsController() artifact.CredentialsController {
	var (
		artifactCredentialsController artifact.CredentialsController
		err                           error
	)

	artifactsCredentialsConfigDir := os.Getenv("ARTIFACTS_CREDENTIALS_CONFIG_DIR")
	if artifactsCredentialsConfigDir == "" {
		// Use default directory /opt/spinnaker/artifacts/config.
		artifactCredentialsController, err = artifact.NewDefaultCredentialsController()
	} else {
		artifactCredentialsController, err = artifact.NewCredentialsController(artifactsCredentialsConfigDir)
	}

	if err != nil {
		log.Println("[CLOUDDRIVER] error setting up artifact credentials controller:", err.Error())
	}

	return artifactCredentialsController
}

// dialector defines the SQL dialector.
//
// Defaults to sqlite if env vars DB_HOST, DB_NAME, DB_PASS, and DB_USER
// are not all defined. Otherwise it creates a MySQL dialctor with a DSN in the format
// `<USER>:<PASS>@tcp(<HOST>)/<NAME>?timeout=30s&charset=utf8&parseTime=True&loc=UTC`.
//
// See https://gorm.io/docs/connecting_to_the_database.html for more info.
func dialector() gorm.Dialector {
	var dialector gorm.Dialector

	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")
	user := os.Getenv("DB_USER")

	// Default to SQLite, else define a MySQL connection.
	if host == "" || name == "" || pass == "" || user == "" {
		log.Println("[CLOUDDRIVER] DB_HOST, DB_NAME, DB_PASS, or DB_USER not defined; defaulting to local SQLite DB")

		dialector = sqlite.Open("clouddriver.db")
	} else {
		dialector = mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s&charset=utf8&parseTime=True&loc=UTC",
				user, pass, host, name),
			DefaultStringSize: mysqlDefaultStringSize,
		})
	}

	return dialector
}
