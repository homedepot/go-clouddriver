package kubernetes

import (
	"context"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/patcher"
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
	ApplyWithNamespaceOverride(*unstructured.Unstructured, string) (Metadata, error)
	DeleteResourceByKindAndNameAndNamespace(string, string, string, metav1.DeleteOptions) error
	GVRForKind(string) (schema.GroupVersionResource, error)
	Get(string, string, string) (*unstructured.Unstructured, error)
	ListByGVR(schema.GroupVersionResource, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	ListByGVRWithContext(context.Context, schema.GroupVersionResource, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	ListResource(string, metav1.ListOptions) (*unstructured.UnstructuredList, error)
	Patch(string, string, string, []byte) (Metadata, *unstructured.Unstructured, error)
	PatchUsingStrategy(string, string, string, []byte, types.PatchType) (Metadata, *unstructured.Unstructured, error)
	ListResourcesByKindAndNamespace(string, string, metav1.ListOptions) (*unstructured.UnstructuredList, error)
}

type client struct {
	c      dynamic.Interface
	config *rest.Config
	mapper *restmapper.DeferredDiscoveryRESTMapper
}

// Apply a given manifest.
func (c *client) Apply(u *unstructured.Unstructured) (Metadata, error) {
	return c.ApplyWithNamespaceOverride(u, "")
}

// Apply a given manifest with an optional namespace to override.
// If no namespace is set on the manifest and no namespace override is passed in then we set the namespace to 'default'.
// If namespaceOverride is empty it will NOT override the namespace set on the manifest.
// We only override the namespace if the manifest is NOT cluster scoped (i.e. a ClusterRole) and namespaceOverride is NOT an
// empty string.
func (c *client) ApplyWithNamespaceOverride(u *unstructured.Unstructured, namespaceOverride string) (Metadata, error) {
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

	if namespaceOverride == "" {
		SetDefaultNamespaceIfScopedAndNoneSet(u, helper)
	} else {
		SetNamespaceIfScoped(namespaceOverride, u, helper)
	}

	info := &resource.Info{
		Client:          restClient,
		Mapping:         restMapping,
		Namespace:       u.GetNamespace(),
		Name:            u.GetName(),
		Source:          "",
		Object:          u,
		ResourceVersion: restMapping.Resource.Version,
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

		_ = info.Refresh(obj, true)
	}

	_, patchedObject, err := patcher.Patch(info.Object, modified, info.Namespace, info.Name)
	if err != nil {
		return metadata, err
	}

	_ = info.Refresh(patchedObject, true)

	metadata.Name = u.GetName()
	metadata.Namespace = u.GetNamespace()
	metadata.Group = gvr.Group
	metadata.Resource = gvr.Resource

	annotations := u.GetAnnotations()
	if annotations != nil {
		if _, ok := annotations[AnnotationSpinnakerArtifactVersion]; ok {
			metadata.Version = annotations[AnnotationSpinnakerArtifactVersion]
		}
	}
	
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

// List all resources by their kind or resource (e.g. "replicaset" or "replicasets")
func (c *client) ListResource(resource string, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	gvr, err := c.GVRForKind(resource)
	if err != nil {
		return nil, err
	}

	return c.c.Resource(gvr).List(context.TODO(), lo)
}

func (c *client) ListResourcesByKindAndNamespace(kind, namespace string, lo metav1.ListOptions) (*unstructured.UnstructuredList, error) {
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
			List(context.TODO(), lo)
	} else {
		ul, err = c.c.
			Resource(restMapping.Resource).
			List(context.TODO(), lo)
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
	annotations := u.GetAnnotations()
	if annotations != nil {
		if _, ok := annotations[AnnotationSpinnakerArtifactVersion]; ok {
			metadata.Version = annotations[AnnotationSpinnakerArtifactVersion]
		}
	}
	metadata.Kind = gvk.Kind

	return metadata, u, err
}
