package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/deployment"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/replicaset"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
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
	Apply([]byte, string) (*unstructured.Unstructured, Metadata, error)
	Patch([]byte) (*unstructured.Unstructured, error)
	Get(string, string, string) (*unstructured.Unstructured, error)
	List(schema.GroupVersionResource, metav1.ListOptions) (*unstructured.UnstructuredList, error)
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

// Apply a given manifest.
func (c *client) Apply(manifest []byte, application string) (*unstructured.Unstructured, Metadata, error) {
	metadata := Metadata{}

	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(manifest, nil, nil)
	if err != nil {
		return nil, metadata, err
	}

	// Convert the runtime.Object to unstructured.Unstructured.
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, metadata, err
	}

	unstructuredObj := &unstructured.Unstructured{
		Object: m,
	}

	name := unstructuredObj.GetName()
	namespace := unstructuredObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	unstructuredObj.SetNamespace(namespace)

	gvk := obj.GetObjectKind().GroupVersionKind()
	t := fmt.Sprintf("kubernetes/%s", strings.ToLower(gvk.Kind))
	cluster := fmt.Sprintf("%s %s", gvk.Kind, name)

	// Add reserved annotations.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-annotations
	annotate(unstructuredObj, "artifact.spinnaker.io/location", namespace)
	annotate(unstructuredObj, "artifact.spinnaker.io/name", name)
	annotate(unstructuredObj, "artifact.spinnaker.io/type", t)
	annotate(unstructuredObj, "moniker.spinnaker.io/application", application)
	annotate(unstructuredObj, "moniker.spinnaker.io/cluster", cluster)

	// Add reserved labels.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-labels
	label(unstructuredObj, "app.kubernetes.io/name", application)
	label(unstructuredObj, "app.kubernetes.io/managed-by", "spinnaker")

	// If this is a deployemnt, set the .spec.template.metadata.* info same as above.
	if strings.EqualFold(gvk.Kind, "deployment") {
		d := deployment.New(unstructuredObj.Object)

		// Add spinnaker annotations to the deployment pod template.
		deployment.AnnotateTemplate(d, "artifact.spinnaker.io/location", namespace)
		deployment.AnnotateTemplate(d, "artifact.spinnaker.io/name", name)
		deployment.AnnotateTemplate(d, "artifact.spinnaker.io/type", t)
		deployment.AnnotateTemplate(d, "moniker.spinnaker.io/application", application)
		deployment.AnnotateTemplate(d, "moniker.spinnaker.io/cluster", cluster)

		// Add reserved labels.
		deployment.LabelTemplate(d, "app.kubernetes.io/name", application)
		deployment.LabelTemplate(d, "app.kubernetes.io/managed-by", "spinnaker")

		unstructuredObj, err = deployment.ToUnstructured(d)
		if err != nil {
			return nil, metadata, err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := replicaset.New(unstructuredObj.Object)

		// Add spinnaker annotations to the replicaset pod template.
		rs.AnnotateTemplate("artifact.spinnaker.io/location", namespace)
		rs.AnnotateTemplate("artifact.spinnaker.io/name", name)
		rs.AnnotateTemplate("artifact.spinnaker.io/type", t)
		rs.AnnotateTemplate("moniker.spinnaker.io/application", application)
		rs.AnnotateTemplate("moniker.spinnaker.io/cluster", cluster)

		// Add reserved labels.
		rs.LabelTemplate("app.kubernetes.io/name", application)
		rs.LabelTemplate("app.kubernetes.io/managed-by", "spinnaker")

		unstructuredObj, err = rs.ToUnstructured()
		if err != nil {
			return nil, metadata, err
		}
	}

	restMapping, err := findGVR(&gvk, c.config)
	if err != nil {
		return nil, metadata, err
	}

	gvr := restMapping.Resource

	resource, err := c.apply(gvr, unstructuredObj)

	metadata.Name = name
	metadata.Namespace = namespace
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource
	metadata.Version = gvr.Version
	metadata.Kind = gvk.Kind

	return resource, metadata, nil
}

func (c *client) apply(gvr schema.GroupVersionResource, o *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	resource := &unstructured.Unstructured{}

	_, err := c.c.
		Resource(gvr).
		Namespace(o.GetNamespace()).
		Get(context.TODO(), o.GetName(), metav1.GetOptions{})
	if err != nil {
		resource, err = c.c.
			Resource(gvr).
			Namespace(o.GetNamespace()).
			Create(context.TODO(), o, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		b, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}
		resource, err = c.c.
			Resource(gvr).
			Namespace(o.GetNamespace()).
			Patch(context.TODO(), o.GetName(), types.StrategicMergePatchType, b, metav1.PatchOptions{})
		if err != nil {
			return nil, err
		}
	}

	return resource, nil
}

// Patch a given manifest
func (c *client) Patch(manifest []byte) (*unstructured.Unstructured, error) {
	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(manifest, nil, nil)
	if err != nil {
		return nil, err
	}

	name, err := meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, err
	}

	namespace, err := meta.NewAccessor().Namespace(obj)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = "default"
	}

	gvk := obj.GetObjectKind().GroupVersionKind()

	restMapping, err := findGVR(&gvk, c.config)
	if err != nil {
		return nil, err
	}

	return c.c.
		Resource(restMapping.Resource).
		Namespace(namespace).
		Patch(context.TODO(), name, types.StrategicMergePatchType, manifest, metav1.PatchOptions{})
}

// Get a manifest by resource (ex 'pods'), name (ex 'my-pod'), and namespace (ex 'my-namespace').
func (c *client) Get(resource, name, namespace string) (*unstructured.Unstructured, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(c.config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	gvk, err := mapper.KindFor(schema.GroupVersionResource{Resource: resource})
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
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *client) List(gvr schema.GroupVersionResource, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.c.Resource(gvr).List(context.TODO(), lo)
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
