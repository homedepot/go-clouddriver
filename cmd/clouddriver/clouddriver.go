package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	dependenciesDir = `/home/spinnaker/.hal/spinnaker-us-central1/staging/dependencies`
)

var (
	r              = gin.New()
	kubeconfigPath string
	kubeconfig     *rest.Config
	client         *kubernetes.Clientset
	// this WILL go away
	cache = map[string][]runtime.Object{}
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

	kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	client, err = kubernetes.NewForConfig(kubeconfig)
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
		"/spectator/metrics",
		"/applications", // TODO
	}}))

	r.Use(gin.Recovery())

	r.GET("/health", func(*gin.Context) {})
	r.GET("/credentials", func(c *gin.Context) {
		credentials := []Credential{}
		sca := Credential{
			AccountType:                 "spin-cluster-account",
			ChallengeDestructiveActions: false,
			CloudProvider:               "kubernetes",
			Environment:                 "spin-cluster-account",
			Name:                        "spin-cluster-account",
			Permissions: struct {
				READ  []string "json:\"READ\""
				WRITE []string "json:\"WRITE\""
			}{
				READ: []string{
					"gg_cloud_gcp_spinnaker_admins",
				},
				WRITE: []string{
					"gg_cloud_gcp_spinnaker_admins",
				},
			},
			PrimaryAccount:          false,
			ProviderVersion:         "v2",
			RequiredGroupMembership: []string{},
			Type:                    "kubernetes",
		}
		credentials = append(credentials, sca)
		c.JSON(http.StatusOK, credentials)
	})

	// Step 1 of a deploy (manifest) stage?
	r.POST("/kubernetes/ops", func(c *gin.Context) {
		id := uuid.New().String()
		kor := KubernetesOpsRequest{}

		err := c.BindJSON(&kor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, req := range kor {
			for _, manifest := range req.DeployManifest.Manifests {
				b, err := json.Marshal(manifest)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				obj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, b)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				log.Println("applying object", obj.GetObjectKind().GroupVersionKind().Kind)

				o, err := getObject(client, *kubeconfig, obj)
				if err != nil {
					o, err = createObject(client, *kubeconfig, obj)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				} else {
					o, err = patchObject(client, *kubeconfig, obj, b)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}

				if _, ok := cache[id]; !ok {
					cache[id] = []runtime.Object{}
				}

				ro := cache[id]
				ro = append(ro, o)
				cache[id] = ro
			}
		}

		or := KubernetesOpsResponse{
			ID:          id,
			ResourceURI: "/task/" + id,
		}
		c.JSON(http.StatusOK, or)
	})

	r.GET("/task/:id", func(c *gin.Context) {
		id := c.Param("id")
		manifests := []runtime.Object{}

		objs := cache[id]
		for _, o := range objs {
			obj, err := getObject(client, *kubeconfig, o)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			manifests = append(manifests, obj)
		}

		ro := ResultObject{
			Manifests: manifests,
		}
		tr := TaskResponse{
			ID:            id,
			ResultObjects: []ResultObject{ro},
			Status: Status{
				Complete:  true,
				Completed: true,
				Failed:    false,
				Phase:     "ORCHESTRATION",
				Retryable: false,
				Status:    "Orchestration completed.",
			},
		}
		c.JSON(http.StatusOK, tr)
	})

	r.Run(":7002")
}

type Credential struct {
	AccountType                 string      `json:"accountType"`
	ChallengeDestructiveActions bool        `json:"challengeDestructiveActions"`
	CloudProvider               string      `json:"cloudProvider"`
	Environment                 string      `json:"environment"`
	Name                        string      `json:"name"`
	Permissions                 Permissions `json:"permissions"`
	PrimaryAccount              bool        `json:"primaryAccount"`
	ProviderVersion             string      `json:"providerVersion"`
	RequiredGroupMembership     []string    `json:"requiredGroupMembership"`
	Type                        string      `json:"type"`
}

type Permissions struct {
	READ  []string `json:"READ"`
	WRITE []string `json:"WRITE"`
}

type KubernetesOpsResponse struct {
	ID          string `json:"id"`
	ResourceURI string `json:"resourceUri"`
}

type TaskResponse struct {
	ID string `json:"id"`
	// SagaIds []interface{} `json:"sagaIds"`
	// History []struct {
	// 	Phase  string `json:"phase"`
	// 	Status string `json:"status"`
	// } `json:"history"`
	// OwnerIDClouddriverSQL   string `json:"ownerId$clouddriver_sql"`
	// RequestIDClouddriverSQL string `json:"requestId$clouddriver_sql"`
	// Retryable                 bool  `json:"retryable"`
	// StartTimeMsClouddriverSQL int64 `json:"startTimeMs$clouddriver_sql"`
	ResultObjects []ResultObject `json:"resultObjects"`
	Status        Status         `json:"status"`
}

type Status struct {
	Complete  bool   `json:"complete"`
	Completed bool   `json:"completed"`
	Failed    bool   `json:"failed"`
	Phase     string `json:"phase"`
	Retryable bool   `json:"retryable"`
	Status    string `json:"status"`
}

type ResultObject struct {
	BoundArtifacts   []interface{} `json:"boundArtifacts"`
	CreatedArtifacts []struct {
		// CustomKind bool   `json:"customKind"`
		// Location   string `json:"location"`
		// Metadata   struct {
		// 	Account string `json:"account"`
		// } `json:"metadata"`
		// Name      string `json:"name"`
		// Reference string `json:"reference"`
		// Type      string `json:"type"`
	} `json:"createdArtifacts"`
	ManifestNamesByNamespace struct {
		// Default []string `json:"default"`
	} `json:"manifestNamesByNamespace"`
	Manifests []runtime.Object `json:"manifests"`
}

type KubernetesOpsRequest []struct {
	DeployManifest struct {
		EnableTraffic     bool                     `json:"enableTraffic"`
		NamespaceOverride string                   `json:"namespaceOverride"`
		OptionalArtifacts []interface{}            `json:"optionalArtifacts"`
		CloudProvider     string                   `json:"cloudProvider"`
		Manifests         []map[string]interface{} `json:"manifests"`
		TrafficManagement struct {
			Options struct {
				EnableTraffic bool `json:"enableTraffic"`
			} `json:"options"`
			Enabled bool `json:"enabled"`
		} `json:"trafficManagement"`
		Moniker struct {
			App string `json:"app"`
		} `json:"moniker"`
		Source                   string        `json:"source"`
		Account                  string        `json:"account"`
		SkipExpressionEvaluation bool          `json:"skipExpressionEvaluation"`
		RequiredArtifacts        []interface{} `json:"requiredArtifacts"`
	} `json:"deployManifest"`
}

func getObject(kubeClientset kubernetes.Interface,
	restConfig rest.Config, obj runtime.Object) (runtime.Object, error) {
	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
	if err != nil {
		return nil, err
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}

	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	name, err := meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, err
	}

	// Create a client specifically for creating the object.
	restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return nil, err
	}

	// Use the REST helper to create the object in the "default" namespace.
	restHelper := resource.NewHelper(restClient, mapping)

	return restHelper.Get("default", name, false)
}

func createObject(kubeClientset kubernetes.Interface,
	restConfig rest.Config, obj runtime.Object) (runtime.Object, error) {
	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
	if err != nil {
		return nil, err
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}

	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	_, err = meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, err
	}

	// Create a client specifically for creating the object.
	restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return nil, err
	}

	// Use the REST helper to create the object in the "default" namespace.
	restHelper := resource.NewHelper(restClient, mapping)

	return restHelper.Create("default", false, obj)
}

func patchObject(kubeClientset kubernetes.Interface,
	restConfig rest.Config, obj runtime.Object, b []byte) (runtime.Object, error) {
	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
	if err != nil {
		return nil, err
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}

	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	name, err := meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, err
	}

	// Create a client specifically for creating the object.
	restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return nil, err
	}

	// Use the REST helper to create the object in the "default" namespace.
	restHelper := resource.NewHelper(restClient, mapping)

	return restHelper.Patch("default", name, types.MergePatchType, b, nil)
}

func newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv

	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	return rest.RESTClientFor(&restConfig)
}
