package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal"
	ops "github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	defaultErrorChanSize    = 1
	defaultManifestChanSize = 1
	manifestListTimeout     = int64(30)
)

// GetManifest returns a manifest for a given account (cluster),
// namespace, kind, and name.
func (cc *Controller) GetManifest(c *gin.Context) {
	includeEvents := c.Query("includeEvents")
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

	events := []v1.Event{}
	errCh := make(chan error, defaultErrorChanSize)
	eventsCh := make(chan v1.Event, internal.DefaultChanSize)
	manifestCh := make(chan *unstructured.Unstructured, defaultManifestChanSize)

	wg := &sync.WaitGroup{}
	// Add 1 to the wait group for getting the manifest or error.
	wg.Add(1)

	go getManifest(provider, wg, manifestCh, errCh, kind, name, namespace)

	if includeEvents != "false" {
		wg.Add(1)

		go getEvents(provider, wg, eventsCh, kind, name, namespace)
	}

	wg.Wait()

	close(eventsCh)

	// Receive all events.
	for event := range eventsCh {
		events = append(events, event)
	}

	var manifest *unstructured.Unstructured
	// Receive either the manifest or the error getting the manifest.
	select {
	case manifest = <-manifestCh:
		break
	case err = <-errCh:
		break
	}

	close(errCh)
	close(manifestCh)

	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cluster := fmt.Sprintf("%s %s", kind, name)
	app := "unknown"

	annotations := manifest.GetAnnotations()
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
		Events:   events,
		Location: namespace,
		Manifest: internal.DeleteNilValues(manifest.Object),
		Metrics:  []interface{}{},
		Moniker: ops.ManifestResponseMoniker{
			App:     app,
			Cluster: cluster,
		},
		Name: fmt.Sprintf("%s %s", kind, name),
		// The 'default' status of a kubernetes resource.
		Status:   kubernetes.GetStatus(kind, manifest.Object),
		Warnings: []interface{}{},
	}

	c.JSON(http.StatusOK, kmr)
}

func getManifest(provider *kubernetes.Provider,
	wg *sync.WaitGroup, manifestCh chan *unstructured.Unstructured,
	errCh chan error, kind, name, namespace string) {
	defer wg.Done()

	manifest, err := provider.Client.Get(kind, name, namespace)
	if err != nil {
		errCh <- err
		return
	}

	manifestCh <- manifest
}

func getEvents(provider *kubernetes.Provider,
	wg *sync.WaitGroup, eventsCh chan v1.Event,
	kind, name, namespace string) {
	defer wg.Done()
	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
	defer cancel()

	events, err := provider.Clientset.Events(ctx, kind, name, namespace)
	if err != nil {
		return
	}

	for _, event := range events {
		eventsCh <- event
	}
}

func (cc *Controller) GetManifestByCriteria(c *gin.Context) {
	account := c.Param("account")
	application := c.Param("application")
	namespace := c.Param("location")
	kind := c.Param("kind")
	cluster := c.Param("cluster")
	// Criteria can be newest, second_newest, oldest, largest, smallest.
	criteria := c.Param("criteria")

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

	sortAscending(items, criteria)

	var manifest unstructured.Unstructured

	// Criteria can be newest, second_newest, oldest, largest, smallest.
	//
	// Java source code here: https://github.com/spinnaker/clouddriver/blob/0fb3e75faa586f213a39c9fd4145f08e519b2e97/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/controllers/ManifestController.java#L132-L148
	switch criteria {
	case "oldest", "smallest":
		manifest = items[0]
	case "newest", "largest":
		manifest = items[len(items)-1]
	case "second_newest":
		if len(items) < 2 {
			clouddriver.Error(c, http.StatusBadRequest,
				errors.New("requested target \"Second Newest\" for cluster "+cluster+", but only one resource was found"))
			return
		}

		manifest = items[len(items)-2]
	default:
		clouddriver.Error(c, http.StatusBadRequest,
			fmt.Errorf("unknown criteria: %s", criteria))
		return
	}

	mcr := ops.ManifestCoordinatesResponse{
		Kind:      kind,
		Name:      manifest.GetName(),
		Namespace: manifest.GetNamespace(),
	}

	c.JSON(http.StatusOK, mcr)
}

// sortAscending sorts an unstructured slice ascending based on criteria. For criteria
// of 'oldest', 'newest', and 'second_newest', it sorts by age: creation timestamp ascending.
// For criteria of 'largest' and 'smallest' it sorts by number of replicas at the JSON path
// `.spec.replicas`.
//
// Java source code comparators here: https://github.com/spinnaker/clouddriver/blob/0fb3e75faa586f213a39c9fd4145f08e519b2e97/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/KubernetesHandler.java#L172
func sortAscending(ul []unstructured.Unstructured, criteria string) {
	switch criteria {
	case "oldest", "newest", "second_newest":
		sort.Slice(ul, func(i, j int) bool {
			return ul[i].GetCreationTimestamp().String() < ul[j].GetCreationTimestamp().String()
		})
	case "largest", "smallest":
		sort.Slice(ul, func(i, j int) bool {
			ir, _, _ := unstructured.NestedInt64(ul[i].Object, "spec", "replicas")
			jr, _, _ := unstructured.NestedInt64(ul[j].Object, "spec", "replicas")

			return ir < jr
		})
	}
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
