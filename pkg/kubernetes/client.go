package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/deployment"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/replicaset"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/kubectl/pkg/util"
)

const (
	ClientInstanceKey                     = `KubeClient`
	AnnotationSpinnakerArtifactLocation   = `artifact.spinnaker.io/location`
	AnnotationSpinnakerArtifactName       = `artifact.spinnaker.io/name`
	AnnotationSpinnakerArtifactType       = `artifact.spinnaker.io/type`
	AnnotationSpinnakerMonikerApplication = `moniker.spinnaker.io/application`
	AnnotationSpinnakerMonikerCluster     = `moniker.spinnaker.io/cluster`
	LabelKubernetesSpinnakerApp           = `app.kubernetes.io/spinnaker-app`
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	LabelKubernetesName      = `app.kubernetes.io/name`
	LabelKubernetesManagedBy = `app.kubernetes.io/managed-by`
	spinnaker                = `spinnaker`
)

// Wrapper for kubernetes dynamic client to make testing easier.

//go:generate counterfeiter . Client
type Client interface {
	SetDynamicClientForConfig(*rest.Config) error
	WithConfig(*rest.Config)
	Apply([]byte, string) (Metadata, error)
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

func (c *client) SetDynamicClientForConfig(config *rest.Config) error {
	d, err := dynamic.NewForConfig(config)
	c.c = d
	c.config = config

	return err
}

func (c *client) WithConfig(config *rest.Config) {
	c.config = config
}

// Apply a given manifest.
func (c *client) Apply(manifest []byte, application string) (Metadata, error) {
	metadata := Metadata{}

	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(manifest, nil, nil)
	if err != nil {
		return metadata, err
	}

	// Convert the runtime.Object to unstructured.Unstructured.
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return metadata, err
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
	cluster := fmt.Sprintf("%s %s", strings.ToLower(gvk.Kind), name)

	// Add reserved annotations.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-annotations
	annotate(unstructuredObj, AnnotationSpinnakerArtifactLocation, namespace)
	annotate(unstructuredObj, AnnotationSpinnakerArtifactName, name)
	annotate(unstructuredObj, AnnotationSpinnakerArtifactType, t)
	annotate(unstructuredObj, AnnotationSpinnakerMonikerApplication, application)
	annotate(unstructuredObj, AnnotationSpinnakerMonikerCluster, cluster)

	// Add reserved labels. Had some trouble with setting the kubernetes name as
	// this interferes with label selectors, so I changed that to be spinnaker-app.
	//
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-labels
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	// label(unstructuredObj, LabelKubernetesName, application)
	label(unstructuredObj, LabelKubernetesSpinnakerApp, application)
	label(unstructuredObj, LabelKubernetesManagedBy, spinnaker)

	// If this is a deployemnt, set the .spec.template.metadata.* info same as above.
	if strings.EqualFold(gvk.Kind, "deployment") {
		d := deployment.New(unstructuredObj.Object)

		// Add spinnaker annotations to the deployment pod template.
		d.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		d.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		d.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		// Add reserved labels.
		// d.LabelTemplate(LabelKubernetesName, application)
		d.LabelTemplate(LabelKubernetesSpinnakerApp, application)
		d.LabelTemplate(LabelKubernetesManagedBy, spinnaker)

		unstructuredObj, err = d.ToUnstructured()
		if err != nil {
			return metadata, err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := replicaset.New(unstructuredObj.Object)

		// Add spinnaker annotations to the replicaset pod template.
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		// Add reserved labels.
		// rs.LabelTemplate(LabelKubernetesName, application)
		rs.LabelTemplate(LabelKubernetesSpinnakerApp, application)
		rs.LabelTemplate(LabelKubernetesManagedBy, spinnaker)

		unstructuredObj, err = rs.ToUnstructured()
		if err != nil {
			return metadata, err
		}
	}

	restMapping, err := findGVR(&gvk, c.config)
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
		Namespace:       unstructuredObj.GetNamespace(),
		Name:            unstructuredObj.GetName(),
		Source:          "",
		Object:          unstructuredObj,
		ResourceVersion: restMapping.Resource.Version,
		Export:          false,
	}

	patcher, err := newPatcher(info, helper)
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

	// func (p *Patcher) Patch(current runtime.Object, modified []byte, namespace, name string) ([]byte, runtime.Object, error) {
	_, patchedObject, err := patcher.Patch(info.Object, modified, info.Namespace, info.Name)
	if err != nil {
		return metadata, err
	}

	info.Refresh(patchedObject, true)

	metadata.Name = name
	metadata.Namespace = namespace
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

// Get a manifest by resource/kind (example: 'pods' or 'pod'),
// name (example: 'my-pod'), and namespace (example: 'my-namespace').
func (c *client) Get(resource, name, namespace string) (*unstructured.Unstructured, error) {
	log.Printf("getting resource (%s) name (%s) namespace (%s)\n", resource, name, namespace)

	dc, err := discovery.NewDiscoveryClientForConfig(c.config)
	if err != nil {
		log.Printf("error 1: error getting resource (%s) name (%s) namespace (%s)\n", resource, name, namespace)
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	gvk, err := mapper.KindFor(schema.GroupVersionResource{Resource: resource})
	if err != nil {
		log.Printf("error 2: error getting resource (%s) name (%s) namespace (%s)\n", resource, name, namespace)
		return nil, err
	}

	restMapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		log.Printf("error 3: error getting resource (%s) name (%s) namespace (%s)\n", resource, name, namespace)
		return nil, err
	}

	// Try to get the resource at the namespace scope.
	u, err := c.c.
		Resource(restMapping.Resource).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		// Try again at the cluster scope.
		u, err = c.c.
			Resource(restMapping.Resource).
			Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

// List all resources by their GVR and list options.
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
