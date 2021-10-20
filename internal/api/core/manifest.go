package core

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	ops "github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	manifestListTimeout = int64(30)
)

// GetManifest returns a manifest for a given account (cluster),
// namespace, kind, and name.
func (cc *Controller) GetManifest(c *gin.Context) {
	account := c.Param("account")
	namespace := c.Param("location")
	// The name of this param should really be "id" or "cluster" as it
	// represents a Spinnaker cluster, such as "deployment my-deployment".
	// However, we have to make this path param match because of an underlying
	// httprouter issue https://github.com/gin-gonic/gin/issues/2016.
	n := c.Param("kind")
	a := strings.Split(n, " ")
	kind := a[0]
	name := a[1]

	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	if strings.Contains(kind, ".") {
		a2 := strings.Split(kind, ".")
		kind = a2[0]
	}

	provider, err := cc.KubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	result, err := provider.Client.Get(kind, name, namespace)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cluster := fmt.Sprintf("%s %s", kind, name)
	app := "unknown"

	annotations := result.GetAnnotations()
	if annotations != nil {
		if _, ok := annotations[kubernetes.AnnotationSpinnakerMonikerApplication]; ok {
			app = annotations[kubernetes.AnnotationSpinnakerMonikerApplication]
		}

		if _, ok := annotations[kubernetes.AnnotationSpinnakerMonikerCluster]; ok {
			cluster = annotations[kubernetes.AnnotationSpinnakerMonikerCluster]
		}
	}

	kmr := ops.ManifestResponse{
		Account:  account,
		Events:   []interface{}{},
		Location: namespace,
		Manifest: result.Object,
		Metrics:  []interface{}{},
		Moniker: ops.ManifestResponseMoniker{
			App:     app,
			Cluster: cluster,
		},
		Name: fmt.Sprintf("%s %s", kind, name),
		// The 'default' status of a kubernetes resource.
		Status:   kubernetes.GetStatus(kind, result.Object),
		Warnings: []interface{}{},
	}

	c.JSON(http.StatusOK, kmr)
}

func (cc *Controller) GetManifestByTarget(c *gin.Context) {
	account := c.Param("account")
	application := c.Param("application")
	namespace := c.Param("location")
	kind := c.Param("kind")
	cluster := c.Param("cluster")
	// Target can be newest, second_newest, oldest, largest, smallest.
	target := c.Param("target")

	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	if strings.Contains(kind, ".") {
		a2 := strings.Split(kind, ".")
		kind = a2[0]
	}

	provider, err := cc.KubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	gvr, err := provider.Client.GVRForKind(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	lo := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: gvr.Group + "/" + gvr.Version,
		},
		LabelSelector:  kubernetes.DefaultLabelSelector(),
		FieldSelector:  "metadata.namespace=" + namespace,
		TimeoutSeconds: &manifestListTimeout,
	}

	list, err := provider.Client.ListByGVR(gvr, lo)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	items := kubernetes.FilterOnAnnotation(list.Items,
		kubernetes.AnnotationSpinnakerMonikerCluster, cluster)
	// Filter out all unassociated objects based on the moniker.spinnaker.io/application annotation.
	items = kubernetes.FilterOnAnnotation(items,
		kubernetes.AnnotationSpinnakerMonikerApplication, application)
	if len(items) == 0 {
		clouddriver.Error(c, http.StatusNotFound, errors.New("no resources found for cluster "+cluster))
		return
	}

	// For now, we sort on creation timestamp to grab the manifest.
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetCreationTimestamp().String() > items[j].GetCreationTimestamp().String()
	})

	var result = items[0]

	// Target can be newest, second_newest, oldest, largest, smallest.
	// TODO fill in for largest and smallest targets.
	switch strings.ToLower(target) {
	case "newest":
		result = items[0]
	case "second_newest":
		if len(items) < 2 {
			clouddriver.Error(c, http.StatusBadRequest,
				errors.New("requested target \"Second Newest\" for cluster "+cluster+", but only one resource was found"))
			return
		}

		result = items[1]
	case "oldest":
		if len(items) < 2 {
			clouddriver.Error(c, http.StatusBadRequest,
				errors.New("requested target \"Oldest\" for cluster "+cluster+", but only one resource was found"))
			return
		}

		result = items[len(items)-1]
	default:
		clouddriver.Error(c, http.StatusNotImplemented,
			errors.New("requested target \""+target+"\" for cluster "+cluster+" is not supported"))
		return
	}

	mcr := ops.ManifestCoordinatesResponse{
		Kind:      kind,
		Name:      result.GetName(),
		Namespace: result.GetNamespace(),
	}

	c.JSON(http.StatusOK, mcr)
}

// ListManifestsByCluster returns a list of manifest coordinates
// for a given account, namespace, location, kind, and cluster.
func (cc *Controller) ListManifestsByCluster(c *gin.Context) {
	account := c.Param("account")
	application := c.Param("application")
	namespace := c.Param("location")
	kind := c.Param("kind")
	cluster := c.Param("cluster")
	manifests := []ops.ManifestCoordinatesResponse{}

	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	if strings.Contains(kind, ".") {
		a := strings.Split(kind, ".")
		kind = a[0]
	}

	provider, err := cc.KubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	gvr, err := provider.Client.GVRForKind(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	lo := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: gvr.Group + "/" + gvr.Version,
		},
		LabelSelector:  kubernetes.DefaultLabelSelector(),
		FieldSelector:  "metadata.namespace=" + namespace,
		TimeoutSeconds: &manifestListTimeout,
	}

	list, err := provider.Client.ListByGVR(gvr, lo)
	if err != nil {
		// Do not error here, just log and return an empty list.
		// This is the expected response from OSS Clouddriver.
		clouddriver.Log(err)
		c.JSON(http.StatusOK, manifests)

		return
	}

	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	items := kubernetes.FilterOnAnnotation(list.Items,
		kubernetes.AnnotationSpinnakerMonikerCluster, cluster)
	// Filter out all unassociated objects based on the moniker.spinnaker.io/application annotation.
	items = kubernetes.FilterOnAnnotation(items,
		kubernetes.AnnotationSpinnakerMonikerApplication, application)

	// Sort by name ascending.
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetName() < items[j].GetName()
	})

	for _, item := range items {
		m := ops.ManifestCoordinatesResponse{
			Kind:      lowercaseFirst(item.GetKind()),
			Name:      item.GetName(),
			Namespace: item.GetNamespace(),
		}
		manifests = append(manifests, m)
	}

	c.JSON(http.StatusOK, manifests)
}
