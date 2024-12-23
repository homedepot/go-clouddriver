package kubernetes

import (
	"context"
	"fmt"

	gcpatcher "github.com/homedepot/go-clouddriver/internal/kubernetes/patcher"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
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
	Replace(*unstructured.Unstructured) (Metadata, error)
	DeleteResourceByKindAndNameAndNamespace(string, string, string, metav1.DeleteOptions) error
	Discover() error
	GVRForKind(string) (schema.GroupVersionResource, error)
	Get(string, string, string) (*unstructured.Unstructured, error)
	ListByGVR(schema.GroupVersionResource, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	ListByGVRWithContext(context.Context, schema.GroupVersionResource, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	ListResource(string, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	ListResourceWithContext(context.Context, string, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	Patch(string, string, string, []byte) (Metadata, *unstructured.Unstructured, error)
	PatchUsingStrategy(string, string, string, []byte, types.PatchType) (Metadata, *unstructured.Unstructured, error)
	ListResourcesByKindAndNamespace(string, string, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	ListResourcesByKindAndNamespaceWithContext(context.Context, string, string, metav1.ListOptions) (*unstructured.UnstructuredList, error)
}

type client struct {
	c      dynamic.Interface
	config *rest.Config
	mapper *restmapper.DeferredDiscoveryRESTMapper
}

// Apply a given manifest.
func (c *client) Apply(u *unstructured.Unstructured) (Metadata, error) {
	var serverSideApply bool

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

	info := &resource.Info{
		Client:          restClient,
		Mapping:         restMapping,
		Namespace:       u.GetNamespace(),
		Name:            u.GetName(),
		Source:          "",
		Object:          u,
		ResourceVersion: restMapping.Resource.Version,
	}

	patcher, err := gcpatcher.New(info, helper)
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

	// Check if server-side annotation is set.
	if AnnotationMatches(*u, AnnotationSpinnakerServerSideApply, "true") {
		serverSideApply = true
	}

	// Server-side annotation can also be set to force-conflicts which  will update your resources using server-side
	// apply and becomes the sole manager.
	if AnnotationMatches(*u, AnnotationSpinnakerServerSideApply, "force-conflicts") {
		serverSideApply = true
		patcher.Force = true
	}

	if !serverSideApply {
		if err := info.Get(); err != nil {
			if !errors.IsNotFound(err) {
				return metadata, err
			}

			// Create the resource if it doesn't exist
			// First, update the annotation used by kubectl apply
			if err := util.CreateApplyAnnotation(info.Object, unstructured.UnstructuredJSONScheme); err != nil {
				return metadata, err
			}

			// Then create the resource and skip the three-way merge if not a server-side apply
			obj, err := helper.Create(info.Namespace, true, info.Object)
			if err != nil {
				return metadata, err
			}

			_ = info.Refresh(obj, true)
		}
	}

	_, patchedObject, err := patcher.Patch(info.Object, modified, info.Namespace, info.Name, serverSideApply)
	if err != nil {
		return metadata, err
	}

	_ = info.Refresh(patchedObject, true)

	metadata.Name = u.GetName()
	metadata.Namespace = u.GetNamespace()
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource
	metadata.Kind = gvk.Kind
	metadata.Version = gvr.Version

	return metadata, nil
}

// Replace a given manifest.
func (c *client) Replace(u *unstructured.Unstructured) (Metadata, error) {
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

	info := &resource.Info{
		Client:          restClient,
		Mapping:         restMapping,
		Namespace:       u.GetNamespace(),
		Name:            u.GetName(),
		Source:          "",
		Object:          u,
		ResourceVersion: restMapping.Resource.Version,
	}

	// If annotation kubectl.kubernetes.io/last-applied-configuration exists, then update it.
	err = util.CreateOrUpdateAnnotation(false, info.Object, unstructured.UnstructuredJSONScheme)
	if err != nil {
		return metadata, err
	}

	exists := true
	// Determine if the resource currently exists.
	if _, err := helper.Get(info.Namespace, info.Name); err != nil {
		if !errors.IsNotFound(err) {
			return metadata, err
		}

		exists = false
	}

	if !exists {
		// Create the resource if it doesn't exist.
		obj, err := helper.Create(info.Namespace, true, info.Object)
		if err != nil {
			return metadata, err
		}

		_ = info.Refresh(obj, true)
	} else {
		// Replace the resource if it does exist.
		obj, err := helper.Replace(info.Namespace, info.Name, true, info.Object)
		if err != nil {
			return metadata, err
		}

		_ = info.Refresh(obj, true)
	}

	metadata.Name = u.GetName()
	metadata.Namespace = u.GetNamespace()
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource
	metadata.Kind = gvk.Kind
	metadata.Version = gvr.Version

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

func (c *client) DeleteResourceByKindAndNameAndNamespace(kind, name, namespace string, do metav1.DeleteOptions) error {
	gvk, err := c.mapper.KindFor(schema.GroupVersionResource{Resource: kind})
	if err != nil {
		return err
	}

	restMapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	restClient, err := newRestClient(*c.config, gvk.GroupVersion())
	if err != nil {
		return err
	}

	helper := resource.NewHelper(restClient, restMapping)
	if helper.NamespaceScoped {
		err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Delete(context.TODO(), name, do)
	} else {
		err = c.c.
			Resource(restMapping.Resource).
			Delete(context.TODO(), name, do)
	}

	return err
}

// Discover uses the resource singularize function of a client GVR mapping
// to initialize the API discovery cache.
//
// This should be ran before running any client request operations concurrently.
// First, it will initialize the cache making any future requests use the disk
// cache for API discovery instead of making requests to the cluster. Second,
// since the REST mapper has a mutex lock on API discovery, concurrent requests
// to grab the GVR from the mapper will appear to run serially.
//
// See https://github.com/kubernetes/client-go/blob/f6ce18ae578c8cca64d14ab9687824d9e1305a67/restmapper/discovery.go#L194.
func (c *client) Discover() error {
	// Just use this function call to cache the API discovery.
	_, err := c.mapper.ResourceSingularizer("pods")
	if err != nil {
		return fmt.Errorf("error discovering API: %w", err)
	}

	return nil
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

	var u *unstructured.Unstructured

	helper := resource.NewHelper(restClient, restMapping)
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

func (c *client) GVRForKind(kind string) (schema.GroupVersionResource, error) {
	return c.mapper.ResourceFor(schema.GroupVersionResource{Resource: kind})
}

// List all resources by their GVR and list options.
func (c *client) ListByGVR(gvr schema.GroupVersionResource, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.c.Resource(gvr).List(context.TODO(), lo)
}

// List all resources by their GVR and list options with context,
func (c *client) ListByGVRWithContext(ctx context.Context, gvr schema.GroupVersionResource, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.c.Resource(gvr).List(ctx, lo)
}

// ListResource lists all resources by their kind or resource (e.g. "replicaset" or "replicasets").
func (c *client) ListResource(resource string, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	gvr, err := c.GVRForKind(resource)
	if err != nil {
		return nil, err
	}

	return c.c.Resource(gvr).List(context.TODO(), lo)
}

// ListResourceWithContext lists all resources by their kind or resource (e.g. "replicaset" or "replicasets") with a context.
func (c *client) ListResourceWithContext(ctx context.Context,
	resource string, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	gvr, err := c.GVRForKind(resource)
	if err != nil {
		return nil, err
	}

	return c.c.Resource(gvr).List(ctx, lo)
}

func (c *client) ListResourcesByKindAndNamespace(kind, namespace string, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.ListResourcesByKindAndNamespaceWithContext(context.Background(), kind, namespace, lo)
}

func (c *client) ListResourcesByKindAndNamespaceWithContext(ctx context.Context,
	kind, namespace string, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
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

	var ul *unstructured.UnstructuredList

	helper := resource.NewHelper(restClient, restMapping)
	if helper.NamespaceScoped {
		ul, err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			List(ctx, lo)
	} else {
		ul, err = c.c.
			Resource(restMapping.Resource).
			List(ctx, lo)
	}

	return ul, err
}

func (c *client) Patch(kind, name, namespace string, p []byte) (Metadata, *unstructured.Unstructured, error) {
	return c.PatchUsingStrategy(kind, name, namespace, p, types.StrategicMergePatchType)
}

func (c *client) PatchUsingStrategy(kind, name, namespace string, p []byte, strategy types.PatchType) (Metadata, *unstructured.Unstructured, error) {
	metadata := Metadata{}

	gvk, err := c.mapper.KindFor(schema.GroupVersionResource{Resource: kind})
	if err != nil {
		return metadata, nil, err
	}

	restMapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return metadata, nil, err
	}

	restClient, err := newRestClient(*c.config, gvk.GroupVersion())
	if err != nil {
		return metadata, nil, err
	}

	var u *unstructured.Unstructured

	helper := resource.NewHelper(restClient, restMapping)
	if helper.NamespaceScoped {
		u, err = c.c.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Patch(context.TODO(), name, strategy, p, metav1.PatchOptions{})
	} else {
		u, err = c.c.
			Resource(restMapping.Resource).
			Patch(context.TODO(), name, strategy, p, metav1.PatchOptions{})
	}

	if err != nil {
		return metadata, nil, err
	}

	gvr := restMapping.Resource

	metadata.Name = u.GetName()
	metadata.Namespace = u.GetNamespace()
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource
	metadata.Kind = gvk.Kind
	metadata.Version = gvr.Version

	return metadata, u, err
}
