package kubernetes

import (
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

const (
	ControllerInstanceKey = `KubeController`
)

//go:generate counterfeiter . Controller
type Controller interface {
	NewClient(*rest.Config) (Client, error)
	ToUnstructured(map[string]interface{}) (*unstructured.Unstructured, error)
	AddSpinnakerAnnotations(u *unstructured.Unstructured, application string) error
	AddSpinnakerLabels(u *unstructured.Unstructured, application string) error
	AddSpinnakerVersionAnnotations(u *unstructured.Unstructured, application string, version SpinnakerVersion) error
	GetCurrentVersion(ul *unstructured.UnstructuredList, kind, name string) string
	IsVersioned(u *unstructured.Unstructured) bool
	IncrementVersion(currentVersion string) SpinnakerVersion
}

func NewController() Controller {
	return &controller{}
}

type controller struct{}

func (c *controller) NewClient(config *rest.Config) (Client, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return &client{
		c:      dynamicClient,
		config: config,
		mapper: mapper,
	}, nil
}

func ControllerInstance(c *gin.Context) Controller {
	return c.MustGet(ControllerInstanceKey).(Controller)
}
