package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
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
	config         *rest.Config
	client         dynamic.Interface
	// this WILL go away
	cache     = map[string][]unstructured.Unstructured{}
	decode    = scheme.Codecs.UniversalDeserializer().Decode
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
			RequiredGroupMembership: []interface{}{},
			Skin:                    "v2",
			SpinnakerKindMap: map[string]string{
				"apiService":                     "unclassified",
				"clusterRole":                    "unclassified",
				"clusterRoleBinding":             "unclassified",
				"configMap":                      "configs",
				"controllerRevision":             "unclassified",
				"cronJob":                        "serverGroups",
				"customResourceDefinition":       "unclassified",
				"daemonSet":                      "serverGroups",
				"deployment":                     "serverGroupManagers",
				"event":                          "unclassified",
				"horizontalpodautoscaler":        "unclassified",
				"ingress":                        "loadBalancers",
				"job":                            "serverGroups",
				"limitRange":                     "unclassified",
				"mutatingWebhookConfiguration":   "unclassified",
				"namespace":                      "unclassified",
				"networkPolicy":                  "securityGroups",
				"persistentVolume":               "configs",
				"persistentVolumeClaim":          "configs",
				"pod":                            "instances",
				"podDisruptionBudget":            "unclassified",
				"podPreset":                      "unclassified",
				"podSecurityPolicy":              "unclassified",
				"replicaSet":                     "serverGroups",
				"role":                           "unclassified",
				"roleBinding":                    "unclassified",
				"secret":                         "configs",
				"service":                        "loadBalancers",
				"serviceAccount":                 "unclassified",
				"statefulSet":                    "serverGroups",
				"storageClass":                   "unclassified",
				"validatingWebhookConfiguration": "unclassified",
			},
			Type: "kubernetes",
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

				obj, _, err := decode(b, nil, nil)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// convert the runtime.Object to unstructured.Unstructured
				m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				unstructuredObj := &unstructured.Unstructured{
					Object: m,
				}

				name, err := meta.NewAccessor().Name(obj)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				gvk := obj.GetObjectKind().GroupVersionKind()

				restMapping, err := findGVR(&gvk, config)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				result, err := client.
					Resource(restMapping.Resource).
					Namespace(namespace).
					Get(context.TODO(), name, metav1.GetOptions{})
				if err != nil {
					result, err = client.
						Resource(restMapping.Resource).
						Namespace(namespace).
						Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
					if err != nil {
						panic(err)
					}
				} else {
					result, err = client.
						Resource(restMapping.Resource).
						Namespace(namespace).
						Patch(context.TODO(), name, types.MergePatchType, b, metav1.PatchOptions{})
					if err != nil {
						panic(err)
					}
				}

				if _, ok := cache[id]; !ok {
					cache[id] = []unstructured.Unstructured{}
				}

				ro := cache[id]
				ro = append(ro, *result)
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
		manifests := []map[string]interface{}{}

		objs := cache[id]
		for _, u := range objs {
			obj := u.DeepCopyObject()
			name, err := meta.NewAccessor().Name(obj)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			gvk := obj.GetObjectKind().GroupVersionKind()

			restMapping, err := findGVR(&gvk, config)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result, err := client.
				Resource(restMapping.Resource).
				Namespace(namespace).
				Get(context.TODO(), name, metav1.GetOptions{})
			// obj, err := getObject(client, *kubeconfig, o)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			manifests = append(manifests, result.Object)
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
	AccountType                 string        `json:"accountType"`
	CacheThreads                int           `json:"cacheThreads"`
	ChallengeDestructiveActions bool          `json:"challengeDestructiveActions"`
	CloudProvider               string        `json:"cloudProvider"`
	DockerRegistries            []interface{} `json:"dockerRegistries"`
	Enabled                     bool          `json:"enabled"`
	Environment                 string        `json:"environment"`
	Name                        string        `json:"name"`
	Namespaces                  []string      `json:"namespaces"`
	Permissions                 struct {
		READ  []string `json:"READ"`
		WRITE []string `json:"WRITE"`
	} `json:"permissions"`
	PrimaryAccount          bool              `json:"primaryAccount"`
	ProviderVersion         string            `json:"providerVersion"`
	RequiredGroupMembership []interface{}     `json:"requiredGroupMembership"`
	Skin                    string            `json:"skin"`
	SpinnakerKindMap        map[string]string `json:"spinnakerKindMap"`
	Type                    string            `json:"type"`
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
	Manifests []map[string]interface{} `json:"manifests"`
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

// Find the corresponding GVR (available in *meta.RESTMapping) for gvk.
func findGVR(gvk *schema.GroupVersionKind, cfg *rest.Config) (*meta.RESTMapping, error) {
	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}
