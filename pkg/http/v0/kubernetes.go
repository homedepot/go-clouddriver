package v0

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/pod"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/replicaset"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
)

func CreateKubernetesDeployment(c *gin.Context) {
	sc := sql.Instance(c)
	taskID := uuid.New().String()
	kor := kubernetes.OpsRequest{}

	err := c.BindJSON(&kor)
	if err != nil {
		e := clouddriver.NewError(
			"BadRequest",
			"Error binding request json: "+err.Error(),
			http.StatusBadRequest,
		)
		c.JSON(http.StatusBadRequest, e)

		return
	}

	// TODO this is hacky - need to figure out how to handle providers.
	accountName := kor[0].DeployManifest.Account

	provider, err := sc.GetKubernetesProvider(kor[0].DeployManifest.Account)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error getting provider: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error decoding provider CA data: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: os.Getenv("BEARER_TOKEN"),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error generating dynamic kubernetes client: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	for _, req := range kor {
		for _, manifest := range req.DeployManifest.Manifests {
			b, err := json.Marshal(manifest)
			if err != nil {
				e := clouddriver.NewError(
					"BadRequest",
					"Error marshaling manifest: "+err.Error(),
					http.StatusBadRequest,
				)
				c.JSON(http.StatusBadRequest, e)

				return
			}

			obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(b, nil, nil)
			if err != nil {
				e := clouddriver.NewError(
					"InternalServerError",
					"Error decoding manifest: "+err.Error(),
					http.StatusInternalServerError,
				)
				c.JSON(http.StatusInternalServerError, e)

				return
			}

			// convert the runtime.Object to unstructured.Unstructured
			m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			if err != nil {
				e := clouddriver.NewError(
					"InternalServerError",
					"Error converting manifest: "+err.Error(),
					http.StatusInternalServerError,
				)
				c.JSON(http.StatusInternalServerError, e)

				return
			}

			unstructuredObj := &unstructured.Unstructured{
				Object: m,
			}

			name, err := meta.NewAccessor().Name(obj)
			if err != nil {
				e := clouddriver.NewError(
					"InternalServerError",
					"Error getting manifest name: "+err.Error(),
					http.StatusInternalServerError,
				)
				c.JSON(http.StatusInternalServerError, e)

				return
			}

			namespace, err := meta.NewAccessor().Namespace(obj)
			if err != nil {
				e := clouddriver.NewError(
					"InternalServerError",
					"Error getting manifest namespace: "+err.Error(),
					http.StatusInternalServerError,
				)
				c.JSON(http.StatusInternalServerError, e)

				return
			}

			gvk := obj.GetObjectKind().GroupVersionKind()

			restMapping, err := findGVR(&gvk, config)
			if err != nil {
				e := clouddriver.NewError(
					"InternalServerError",
					"Error getting rest mapping: "+err.Error(),
					http.StatusInternalServerError,
				)
				c.JSON(http.StatusInternalServerError, e)

				return
			}

			gvr := restMapping.Resource

			_, err = client.
				Resource(gvr).
				Namespace(namespace).
				Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				_, err = client.
					Resource(restMapping.Resource).
					Namespace(namespace).
					Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
				if err != nil {
					panic(err)
				}
			} else {
				_, err = client.
					Resource(restMapping.Resource).
					Namespace(namespace).
					Patch(context.TODO(), name, types.StrategicMergePatchType, b, metav1.PatchOptions{})
				if err != nil {
					panic(err)
				}
			}

			kr := kubernetes.Resource{
				AccountName: accountName,
				ID:          uuid.New().String(),
				TaskID:      taskID,
				Group:       gvr.Group,
				Name:        name,
				Namespace:   namespace,
				Resource:    gvr.Resource,
				Version:     gvr.Version,
				Kind:        gvk.Kind,
			}

			err = sc.CreateKubernetesResource(kr)
			if err != nil {
				e := clouddriver.NewError(
					"InternalServerError",
					"Error creating kubernetes resource in db: "+err.Error(),
					http.StatusInternalServerError,
				)
				c.JSON(http.StatusInternalServerError, e)

				return
			}
		}
	}

	or := kubernetes.OpsResponse{
		ID:          taskID,
		ResourceURI: "/task/" + taskID,
	}
	c.JSON(http.StatusOK, or)
}

func GetManifest(c *gin.Context) {
	sc := sql.Instance(c)
	account := c.Param("account")
	namespace := c.Param("location")
	n := c.Param("name")
	a := strings.Split(n, " ")
	resource := a[0]
	name := a[1]

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error getting provider: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error decoding provider CA data: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: os.Getenv("BEARER_TOKEN"),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error generating dynamic kubernetes client: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error getting new discovery client: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	gvk, err := mapper.KindFor(schema.GroupVersionResource{Resource: resource})
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error getting kind: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	restMapping, err := findGVR(&gvk, config)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error finding GVR: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

		return
	}

	result, err := client.
		Resource(restMapping.Resource).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	// obj, err := getObject(client, *kubeconfig, o)
	if err != nil {
		e := clouddriver.NewError(
			"InternalServerError",
			"Error getting resource: "+err.Error(),
			http.StatusInternalServerError,
		)
		c.JSON(http.StatusInternalServerError, e)

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
		Status:   kubernetes.DefaultStatus,
		Warnings: nil,
	}

	// status https://github.com/spinnaker/clouddriver/tree/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler
	// pod status check https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/KubernetesPodHandler.java
	switch strings.ToLower(gvk.GroupKind().Kind) {
	case "pod":
		kmr.Status = pod.Status(result.Object)
	case "replicaset":
		kmr.Status = replicaset.Status(result.Object)
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
