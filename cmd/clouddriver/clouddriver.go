package main

import (
	"log"
	"net/http"
	"os"

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
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{
		"/health",
		"/applications", // TODO
	}}))
	r.Use(gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		// log.Println("HEADERS:", c.Request.Header)
		// b, _ := ioutil.ReadAll(c.Request.Body)
		// log.Println("BODY:", string(b))
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	kubeClient := kubernetes.NewClient()

	sqlClient := sql.NewClient(mustDBConnect())
	c := &server.Config{
		SQLClient:  sqlClient,
		KubeClient: kubeClient,
	}
	server.Setup(r, c)
}

func mustGetenv(env string) (s string) {
	if s = os.Getenv(env); s == "" {
		log.Fatal(env + " not set; exiting.")
	}

	return
}

func mustDBConnect() *gorm.DB {
	sqlConfig := sql.Config{
		User:     mustGetenv("DB_USER"),
		Password: mustGetenv("DB_PASS"),
		Host:     mustGetenv("DB_HOST"),
		Name:     mustGetenv("DB_NAME"),
	}

	db, err := sql.Connect("mysql", sql.Connection(sqlConfig))
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
