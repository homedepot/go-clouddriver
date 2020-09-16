package kubernetes

import (
	"context"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/patcher"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/kubectl/pkg/util"
)

const (
	spinnaker = `spinnaker`
)

type Metadata struct {
	Name      string
	Namespace string
	Group     string
	Version   string
	Resource  string
	Kind      string
}

// Generate a new client using the kubernetes controller.

//go:generate counterfeiter . Client
type Client interface {
	Apply(*unstructured.Unstructured) (Metadata, error)
	Get(string, string, string) (*unstructured.Unstructured, error)
	List(schema.GroupVersionResource, metav1.ListOptions) (*unstructured.UnstructuredList, error)
}

type client struct {
	c      dynamic.Interface
	config *rest.Config
	mapper *restmapper.DeferredDiscoveryRESTMapper
}

// Apply a given manifest.
func (c *client) Apply(u *unstructured.Unstructured) (Metadata, error) {
	metadata := Metadata{}
	gvk := u.GroupVersionKind()

	restMapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return metadata, err
	}

	gvr := restMapping.Resource
	gv := gvk.GroupVersion()
	c.config.GroupVersion = &gv

	restClient, err := newRestClient(*c.config, gv)
	if err != nil {
		return metadata, err
	}

	helper := resource.NewHelper(restClient, restMapping)
	SetDefaultNamespaceIfScopedAndNoneSet(u, helper)

	info := &resource.Info{
		Client:          restClient,
		Mapping:         restMapping,
		Namespace:       u.GetNamespace(),
		Name:            u.GetName(),
		Source:          "",
		Object:          u,
		ResourceVersion: restMapping.Resource.Version,
		Export:          false,
	}

	patcher, err := patcher.New(info, helper)
	if err != nil {
		return metadata, err
	}

	// Get the modified configuration of the object. Embed the result
	// as an annotation in the modified configuration, so that it will appear
	// in the patch sent to the server.
	modified, err := util.GetModifiedConfiguration(info.Object, true, unstructured.UnstructuredJSONScheme)
	if err != nil {
		return metadata, err
	}

	if err := info.Get(); err != nil {
		if !errors.IsNotFound(err) {
			return metadata, err
		}

		// Create the resource if it doesn't exist
		// First, update the annotation used by kubectl apply
		if err := util.CreateApplyAnnotation(info.Object, unstructured.UnstructuredJSONScheme); err != nil {
			return metadata, err
		}

		// Then create the resource and skip the three-way merge
		obj, err := helper.Create(info.Namespace, true, info.Object)
		if err != nil {
			return metadata, err
		}
		info.Refresh(obj, true)
	}

	_, patchedObject, err := patcher.Patch(info.Object, modified, info.Namespace, info.Name)
	if err != nil {
		return metadata, err
	}

	info.Refresh(patchedObject, true)

	metadata.Name = u.GetName()
	metadata.Namespace = u.GetNamespace()
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource
	metadata.Version = gvr.Version
	metadata.Kind = gvk.Kind

	return metadata, nil
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

// Get a manifest by resource/kind (example: 'pods' or 'pod'),
// name (example: 'my-pod'), and namespace (example: 'my-namespace').
func (c *client) Get(kind, name, namespace string) (*unstructured.Unstructured, error) {
	gvk, err := c.mapper.KindFor(schema.GroupVersionResource{Resource: kind})
	if err != nil {
		return nil, err
	}

	restMapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	restClient, err := newRestClient(*c.config, gvk.GroupVersion())
	if err != nil {
		return nil, err
	}

	helper := resource.NewHelper(restClient, restMapping)

	var u *unstructured.Unstructured

	if helper.NamespaceScoped {
		u, err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Get(context.TODO(), name, metav1.GetOptions{})
	} else {
		u, err = c.c.
			Resource(restMapping.Resource).
			Get(context.TODO(), name, metav1.GetOptions{})
	}

	return u, err
}

// List all resources by their GVR and list options.
func (c *client) List(gvr schema.GroupVersionResource, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.c.Resource(gvr).List(context.TODO(), lo)
}
