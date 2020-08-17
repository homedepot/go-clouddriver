package kubernetes

import (
	"context"

	"github.com/gin-gonic/gin"
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

const (
	ClientInstanceKey = `KubeClient`
)

// Wrapper for kubernetes dynamic client to make testing easier.

//go:generate counterfeiter . Client
type Client interface {
	WithConfig(*rest.Config) error
	Apply([]byte) (*unstructured.Unstructured, Metadata, error)
	Get(string, string, string) (*unstructured.Unstructured, error)
}

func NewClient() Client {
	return &client{}
}

type client struct {
	c      dynamic.Interface
	config *rest.Config
}

type Metadata struct {
	Name      string
	Namespace string
	Group     string
	Version   string
	Resource  string
	Kind      string
}

func (c *client) WithConfig(config *rest.Config) error {
	d, err := dynamic.NewForConfig(config)
	c.c = d
	c.config = config

	return err
}

func (c *client) Apply(manifest []byte) (*unstructured.Unstructured, Metadata, error) {
	metadata := Metadata{}

	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(manifest, nil, nil)
	if err != nil {
		return nil, metadata, err
	}

	// convert the runtime.Object to unstructured.Unstructured
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, metadata, err
	}

	unstructuredObj := &unstructured.Unstructured{
		Object: m,
	}

	name, err := meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, metadata, err
	}

	namespace, err := meta.NewAccessor().Namespace(obj)
	if err != nil {
		return nil, metadata, err
	}

	gvk := obj.GetObjectKind().GroupVersionKind()

	restMapping, err := findGVR(&gvk, c.config)
	if err != nil {
		return nil, metadata, err
	}

	gvr := restMapping.Resource

	resource := &unstructured.Unstructured{}

	_, err = c.c.
		Resource(gvr).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		resource, err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return nil, metadata, err
		}
	} else {
		resource, err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Patch(context.TODO(), name, types.StrategicMergePatchType, manifest, metav1.PatchOptions{})
		if err != nil {
			return nil, metadata, err
		}
	}

	metadata.Name = name
	metadata.Namespace = namespace
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource
	metadata.Version = gvr.Version
	metadata.Kind = gvk.Kind

	return resource, metadata, nil
}

// Get a manifest by kind (ex 'pod'), name (ex 'my-pod'), and namespace (ex 'my-namespace').
func (c *client) Get(kind, name, namespace string) (*unstructured.Unstructured, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(c.config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	gvk, err := mapper.KindFor(schema.GroupVersionResource{Resource: kind})
	if err != nil {
		return nil, err
	}

	restMapping, err := findGVR(&gvk, c.config)
	if err != nil {
		return nil, err
	}

	u, err := c.c.
		Resource(restMapping.Resource).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	// obj, err := getObject(client, *kubeconfig, o)
	if err != nil {
		return nil, err
	}

	return u, nil
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

func Instance(c *gin.Context) Client {
	return c.MustGet(ClientInstanceKey).(Client)
}
