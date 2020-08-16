package v0

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

var (
	decode = scheme.Codecs.UniversalDeserializer().Decode
)

func CreateKubernetesDeployment(c *gin.Context) {
	id := uuid.New().String()
	kor := kubernetes.OpsRequest{}

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
					Patch(context.TODO(), name, types.StrategicMergePatchType, b, metav1.PatchOptions{})
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
}

func GetManifest(c *gin.Context) {
	account := c.Param("account")
	namespace := c.Param("location")
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

	kmr := kubernetes.ManifestResponse{
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
		Status: kubernetes.ManifestStatus{
			Available: kubernetes.Available{
				State: true,
			},
			Failed: kubernetes.Failed{
				State: false,
			},
			Paused: kubernetes.Paused{
				State: false,
			},
			Stable: kubernetes.Stable{
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
