package v1

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	// Only need the kinds that show in the infrastructure pages.
	infrastructureKinds = []string{
		"daemonSets",
		"deployments",
		"ingresses",
		"replicaSets",
		"services",
		"statefulSets",
	}
)

// LoadKubernetesResources populates the Kubernetes resources table
// by querying the cluster for resources deployed by Spinnaker.
//
// The only resources populated are Kubernetes kinds that show in
// the infrastructure pages (clusters, load balancers, firewalls).
//   - daemonSets
//   - deployments
//   - ingresses
//   - replicaSets
//   - statefulSets
//   - services
func (cc *Controller) LoadKubernetesResources(c *gin.Context) {
	account := c.Param("name")

	// Grab the kube provider for the given account.
	provider, err := cc.KubernetesProviderWithTimeout(account, time.Second*internal.DefaultListTimeoutSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// First, run discovery on this dynamic client before listing resources
	// concurrently. This is necessary since the rest mapper for dynamic
	// clients uses a mutex lock. Failure to do this will make concurrent
	// requests appear to run serially. This is particularly bad if a cluster is not
	// reachable - even with a timeout of 10 seconds, a request for 4 resources
	// would take 40 seconds since the API cannot be discovered concurrently.
	//
	// See https://github.com/kubernetes/client-go/blob/f6ce18ae578c8cca64d14ab9687824d9e1305a67/restmapper/discovery.go#L194.
	if err = provider.Client.Discover(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Declare a waitgroup to wait on concurrent resource listing.
	wg := &sync.WaitGroup{}
	// Add the number of kinds we will be listing concurrently.
	wg.Add(len(infrastructureKinds))
	// Create channel of unstructured objects (manifests) to send to.
	uc := make(chan unstructured.Unstructured, internal.DefaultChanSize)
	// List all required kinds concurrently.
	for _, kind := range infrastructureKinds {
		go listKinds(wg, uc, provider, kind)
	}
	// Wait for the calls to finish.
	wg.Wait()
	// Close the channel.
	close(uc)

	// Remove existing entries from the DB.
	err = cc.SQLClient.DeleteKubernetesResourcesByAccountName(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use the same task ID for all entries.
	taskID := uuid.New().String()
	resources := []kubernetes.Resource{}

	// Create a resource entry in the DB for each unstructured object (manifest).
	for u := range uc {
		gvr, err := provider.Client.GVRForKind(u.GetKind())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Map values to kubernetes.Resource.
		nameWithoutVersion := kubernetes.NameWithoutVersion(u.GetName())
		kr := kubernetes.Resource{
			AccountName:  account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			Timestamp:    internal.CurrentTimeUTC(),
			APIGroup:     gvr.Group,
			Name:         u.GetName(),
			ArtifactName: nameWithoutVersion,
			Namespace:    u.GetNamespace(),
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         u.GetKind(),
			SpinnakerApp: kubernetes.SpinnakerMonikerApplication(u),
			Cluster:      kubernetes.Cluster(u.GetKind(), nameWithoutVersion),
		}
		// Insert row into kubernetes_resources table.
		err = cc.SQLClient.CreateKubernetesResource(kr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resources = append(resources, kr)
	}

	c.JSON(http.StatusOK, resources)
}

// DeleteKubernetesResources deletes the resources from the database
// for the given provider (account).
func (cc *Controller) DeleteKubernetesResources(c *gin.Context) {
	name := c.Param("name")

	_, err := cc.SQLClient.GetKubernetesProvider(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = cc.SQLClient.DeleteKubernetesResourcesByAccountName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// listKinds lists a given kind and sends to a channel of unstructured.Unstructured.
// It uses a context with a timeout of 10 seconds.
func listKinds(wg *sync.WaitGroup, uc chan unstructured.Unstructured,
	provider *kubernetes.Provider, kind string) {
	// Finish the wait group when we're done here.
	defer wg.Done()

	// Declare server side filtering options.
	lo := metav1.ListOptions{
		LabelSelector: kubernetes.DefaultLabelSelector(),
	}

	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
	defer cancel()

	var ul *unstructured.UnstructuredList

	var err error

	if provider.Namespace != nil {
		ul, err = provider.Client.ListResourcesByKindAndNamespaceWithContext(ctx, kind, *provider.Namespace, lo)
	} else {
		ul, err = provider.Client.ListResourceWithContext(ctx, kind, lo)
	}

	if err != nil {
		clouddriver.Log(err)
		return
	}

	// Send all unstructured objects to the channel.
	if ul == nil {
		return
	}

	for _, u := range ul.Items {
		uc <- u
	}
}
