package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/mitchellh/mapstructure"
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

		b, _ := json.Marshal(kor)
		log.Println(string(b))

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

	// force cache refresh
	// {"account":"spin-cluster-account","location":"default","name":"pod rss-site"}
	r.POST("/cache/kubernetes/manifest", func(c *gin.Context) {
		b, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println("CACHE:", string(b))
	})

	// monitor deploy
	r.GET("/manifests/:account/:location/:name", func(c *gin.Context) {
		account := c.Param("account")
		// default
		namespace := c.Param("location")
		// pod rss-site
		n := c.Param("name")
		a := strings.Split(n, " ")
		resource := a[0]
		name := a[1]

		dc, err := discovery.NewDiscoveryClientForConfig(config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
		gvk, err := mapper.KindFor(schema.GroupVersionResource{Resource: resource})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		kmr := KubernetesManifestResponse{
			Account:  account,
			Events:   nil,
			Location: namespace,
			Manifest: result.Object,
			Metrics:  nil,
			Moniker: struct {
				App     string "json:\"app\""
				Cluster string "json:\"cluster\""
			}{
				App:     "TODO",
				Cluster: "TODO",
			},
			Name: name,
			// The 'default' status of a kubernetes resource.
			Status: KubernetesManifestStatus{
				Available: Available{
					State: true,
				},
				Failed: Failed{
					State: false,
				},
				Paused: Paused{
					State: false,
				},
				Stable: Stable{
					State: true,
				},
			},
			Warnings: nil,
		}

		// status https://github.com/spinnaker/clouddriver/blob/900f2b1013781b290a9d0db96ce1dd964917382f/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/model/Manifest.java#L48
		// pod status check https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/KubernetesPodHandler.java
		switch strings.ToLower(gvk.GroupKind().Kind) {
		case "pod":
			var pod struct {
				Status struct {
					Phase string `json:"phase"`
				} `json:"status"`
			}
			err := mapstructure.Decode(result.Object, &pod)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if strings.EqualFold(pod.Status.Phase, "pending") ||
				strings.EqualFold(pod.Status.Phase, "failed") ||
				strings.EqualFold(pod.Status.Phase, "unknown") {
				kmr.Status.Stable.State = false
				kmr.Status.Stable.Message = "Pod is " + strings.ToLower(pod.Status.Phase)
				kmr.Status.Available.State = false
				kmr.Status.Available.Message = "Pod is " + strings.ToLower(pod.Status.Phase)
			}
		default:
			c.JSON(http.StatusInternalServerError, nil)
		}

		c.JSON(http.StatusOK, kmr)
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
			CreatedArtifacts: []CreatedArtifact{{
				CustomKind: false,
				Location:   "",
				Metadata: struct {
					Account string "json:\"account\""
				}{
					Account: "spin-cluster-account",
				},
				Name:      "rss-site",
				Reference: "rss-site",
				Type:      "kubernetes/pod",
				Version:   "",
			},
			},
			ManifestNamesByNamespace: map[string][]string{
				"default": {"pod rss-site"},
			},
			ManifestNamesByNamespaceToRefresh: map[string][]string{
				"default": {"pod rss-site"},
			},
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
	BoundArtifacts                    []interface{}            `json:"boundArtifacts"`
	CreatedArtifacts                  []CreatedArtifact        `json:"createdArtifacts"`
	ManifestNamesByNamespace          map[string][]string      `json:"manifestNamesByNamespace"`
	ManifestNamesByNamespaceToRefresh map[string][]string      `json:"manifestNamesByNamespaceToRefresh"`
	Manifests                         []map[string]interface{} `json:"manifests"`
}

type CreatedArtifact struct {
	CustomKind bool   `json:"customKind"`
	Location   string `json:"location"`
	Metadata   struct {
		Account string `json:"account"`
	} `json:"metadata"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Type      string `json:"type"`
	Version   string `json:"version"`
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

type KubernetesManifestResponse struct {
	Account string `json:"account"`
	// Artifacts []struct {
	// 	CustomKind bool `json:"customKind"`
	// 	Metadata   struct {
	// 	} `json:"metadata"`
	// 	Name      string `json:"name"`
	// 	Reference string `json:"reference"`
	// 	Type      string `json:"type"`
	// } `json:"artifacts"`
	Events   []interface{}          `json:"events"`
	Location string                 `json:"location"`
	Manifest map[string]interface{} `json:"manifest"`
	Metrics  []interface{}          `json:"metrics"`
	Moniker  struct {
		App     string `json:"app"`
		Cluster string `json:"cluster"`
	} `json:"moniker"`
	Name     string                   `json:"name"`
	Status   KubernetesManifestStatus `json:"status"`
	Warnings []interface{}            `json:"warnings"`
}

type KubernetesManifestStatus struct {
	Available Available `json:"available"`
	Failed    Failed    `json:"failed"`
	Paused    Paused    `json:"paused"`
	Stable    Stable    `json:"stable"`
}

type Available struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type Failed struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type Paused struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type Stable struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}
