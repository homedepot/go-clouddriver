package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	dependenciesDir = `/home/spinnaker/.hal/spinnaker-us-central1/staging/dependencies`
)

var (
	r              = gin.New()
	kubeconfigPath string
	config         *rest.Config
	client         dynamic.Interface
	// this WILL go away
	cache     = map[string][]unstructured.Unstructured{}
	namespace = "default"
)

func init() {
	files, err := ioutil.ReadDir(dependenciesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "spinnaker-us-central1.config") {
			kubeconfigPath = dependenciesDir + "/" + file.Name()
		}
	}

	if kubeconfigPath == "" {
		log.Fatal("unable to get spin-cluster-account config file")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	client, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	gin.ForceConsoleColor()

	p := ginprometheus.NewPrometheus("gin")
	p.MetricsPath = "/metrics"
	p.Use(r)

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{
		"/health",
		"/applications", // TODO
	}}))

	r.Use(gin.Recovery())

	r.Run(":7002")
}
