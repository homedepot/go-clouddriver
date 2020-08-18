package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
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
	Apply([]byte, string) (*unstructured.Unstructured, Metadata, error)
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

func (c *client) Apply(manifest []byte, spinnakerApp string) (*unstructured.Unstructured, Metadata, error) {
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

	name, err := meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, metadata, err
	}

	namespace, err := meta.NewAccessor().Namespace(obj)
	if err != nil {
		return nil, metadata, err
	}

	if namespace == "" {
		namespace = "default"
	}

	gvk := obj.GetObjectKind().GroupVersionKind()

	// Add reserved annotations.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-annotations
	{
		annotations := unstructuredObj.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}

		annotations["artifact.spinnaker.io/location"] = namespace
		annotations["artifact.spinnaker.io/name"] = name
		annotations["artifact.spinnaker.io/type"] = fmt.Sprintf("kubernetes/%s", strings.ToLower(gvk.Kind))
		unstructuredObj.SetAnnotations(annotations)
	}

	// Add reserved labels.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-labels
	{
		labels := unstructuredObj.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}

		labels["app.kubernetes.io/name"] = spinnakerApp
		labels["app.kubernetes.io/managed-by"] = "spinnaker"
		unstructuredObj.SetLabels(labels)
	}

	// If this is a deployemnt, set the .spec.template.metadata.* info same as above.
	if strings.EqualFold(gvk.Kind, "deployment") {
		d := &v1.Deployment{}

		{
			b, err := json.Marshal(unstructuredObj.Object)
			if err != nil {
				return nil, metadata, err
			}

			err = json.Unmarshal(b, &d)
			if err != nil {
				return nil, metadata, err
			}
		}

		// Add reserved annotations.
		annotations := d.Spec.Template.ObjectMeta.Annotations
		if annotations == nil {
			annotations = map[string]string{}
		}

		annotations["artifact.spinnaker.io/location"] = namespace
		annotations["artifact.spinnaker.io/name"] = name
		annotations["artifact.spinnaker.io/type"] = fmt.Sprintf("kubernetes/%s", strings.ToLower(gvk.Kind))
		d.Spec.Template.ObjectMeta.Annotations = annotations

		// Add reserved labels.
		labels := d.Spec.Template.ObjectMeta.Labels
		if labels == nil {
			labels = map[string]string{}
		}

		labels["app.kubernetes.io/name"] = spinnakerApp
		labels["app.kubernetes.io/managed-by"] = "spinnaker"
		d.Spec.Template.ObjectMeta.Labels = labels

		{
			b, err := json.Marshal(d)
			if err != nil {
				return nil, metadata, err
			}

			err = json.Unmarshal(b, &unstructuredObj.Object)
			if err != nil {
				return nil, metadata, err
			}
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := &v1.ReplicaSet{}

		{
			b, err := json.Marshal(unstructuredObj.Object)
			if err != nil {
				return nil, metadata, err
			}

			err = json.Unmarshal(b, &rs)
			if err != nil {
				return nil, metadata, err
			}
		}

		// Add reserved annotations.
		annotations := rs.Spec.Template.ObjectMeta.Annotations
		if annotations == nil {
			annotations = map[string]string{}
		}

		annotations["artifact.spinnaker.io/location"] = namespace
		annotations["artifact.spinnaker.io/name"] = name
		annotations["artifact.spinnaker.io/type"] = fmt.Sprintf("kubernetes/%s", strings.ToLower(gvk.Kind))
		rs.Spec.Template.ObjectMeta.Annotations = annotations

		// Add reserved labels.
		labels := rs.Spec.Template.ObjectMeta.Labels
		if labels == nil {
			labels = map[string]string{}
		}

		labels["app.kubernetes.io/name"] = spinnakerApp
		labels["app.kubernetes.io/managed-by"] = "spinnaker"
		rs.Spec.Template.ObjectMeta.Labels = labels

		{
			b, err := json.Marshal(rs)
			if err != nil {
				return nil, metadata, err
			}

			err = json.Unmarshal(b, &unstructuredObj.Object)
			if err != nil {
				return nil, metadata, err
			}
		}
	}

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
		b, err := json.Marshal(unstructuredObj.Object)
		if err != nil {
			return nil, metadata, err
		}
		resource, err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Patch(context.TODO(), name, types.StrategicMergePatchType, b, metav1.PatchOptions{})
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
	// obj, err := getObject(client, *kubeconfig, o)
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
