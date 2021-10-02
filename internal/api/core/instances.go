package core

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	defaultGetTimeoutSeconds = 10
)

// InstanceRepsponse represents the HTTP response
// when requesting instance information. For Kubernetes,
// an "instance" is Pod.
type InstanceResponse struct {
	Account           string                   `json:"account"`
	Apiversion        string                   `json:"apiVersion"`
	Cloudprovider     string                   `json:"cloudProvider"`
	Createdtime       int64                    `json:"createdTime"`
	Displayname       string                   `json:"displayName"`
	Health            []InstanceResponseHealth `json:"health"`
	Healthstate       string                   `json:"healthState"`
	Humanreadablename string                   `json:"humanReadableName"`
	Kind              string                   `json:"kind"`
	Labels            map[string]string        `json:"labels"`
	Moniker           Moniker                  `json:"moniker"`
	Name              string                   `json:"name"`
	Namespace         string                   `json:"namespace"`
	Providertype      string                   `json:"providerType"`
	Zone              string                   `json:"zone"`
}

// InstanceResponseHealth represents health of an instance,
// which is Kubernetes is a Pod.
type InstanceResponseHealth struct {
	Platform string `json:"platform"`
	Source   string `json:"source"`
	State    string `json:"state"`
	Type     string `json:"type"`
}

// GetInstance grabs an instance by account, location, and name.
// It builds the instance response and calculates health status
// of the instance.
func (cc *Controller) GetInstance(c *gin.Context) {
	account := c.Param("account")
	namespace := c.Param("location")
	n := c.Param("name")
	a := strings.Split(n, " ")
	kind := a[0]
	name := a[1]
	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	//
	// This should never happen on the instances endpoint, but just to be safe!
	if strings.Contains(kind, ".") {
		a2 := strings.Split(kind, ".")
		kind = a2[0]
	}
	// Grab the Kubernetes provider with a default timeout.
	provider, err := cc.KubernetesProviderWithTimeout(account, time.Second*defaultGetTimeoutSeconds)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	// Get the instance (Pod) from the cluster.
	instance, err := provider.Client.Get(kind, name, namespace)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	annotations := instance.GetAnnotations()
	cluster := annotations["moniker.spinnaker.io/cluster"]
	app := annotations["moniker.spinnaker.io/application"]
	// Get overall pod health state and the health status
	// of all its containers.
	healthState, health := calculateInstanceHealth(instance)

	ir := InstanceResponse{
		Account:       account,
		Apiversion:    instance.GetAPIVersion(),
		Cloudprovider: typeKubernetes,
		Createdtime:   instance.GetCreationTimestamp().Unix() * 1000,
		Displayname:   instance.GetName(),
		Health:        health,
		Healthstate:   healthState,
		Humanreadablename: fmt.Sprintf("%s %s",
			lowercaseFirst(instance.GetKind()), instance.GetName()),
		Kind:   lowercaseFirst(instance.GetKind()),
		Labels: instance.GetLabels(),
		Moniker: Moniker{
			App:     app,
			Cluster: cluster,
		},
		// Name is for some reason the UID?
		Name:         string(instance.GetUID()),
		Namespace:    instance.GetNamespace(),
		Providertype: typeKubernetes,
		Zone:         instance.GetNamespace(),
	}

	c.JSON(http.StatusOK, ir)
}

// calculateInstanceHealth returns the health slice of an
// instance. This contains health information for a Kubernetes
// pod and all its containers. The first return argument is the health
// state of the pod and the second is a slice of health
// information for the pod and each container.
func calculateInstanceHealth(instance *unstructured.Unstructured) (string, []InstanceResponseHealth) {
	healthState := stateDown
	health := []InstanceResponseHealth{}
	// Only calculate health info if we know the instance
	// is of kind Pod.
	if strings.EqualFold(instance.GetKind(), "pod") {
		p := kubernetes.NewPod(instance.Object)
		status := p.Object().Status
		// healthState represents the state of the whole Pod.
		healthState = podState(status)
		// Define the Pod's health.
		podHealth := InstanceResponseHealth{
			Platform: "platform",
			Source:   "Pod",
			State:    healthState,
			Type:     "kubernetes/pod",
		}
		health = append(health, podHealth)
		// Get all the Pod's cantainer statuses.
		containerStatuses := status.ContainerStatuses
		for _, cs := range containerStatuses {
			containerState := stateDown
			// For now, we define a healthy container as one
			// that has details about it's running state.
			if cs.State.Running != nil {
				containerState = stateUp
			}
			// Define the container's health info.
			h := InstanceResponseHealth{
				Platform: "platform",
				Source:   fmt.Sprintf("Container %s", cs.Name),
				State:    containerState,
				Type:     "kubernetes/container",
			}
			health = append(health, h)
		}
	}

	return healthState, health
}

// podState returns the "state" of the Pod, which has been simplified
// here to either be "Up" or "Down".
//
// Source code for instance health here:
// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/provider/KubernetesModelUtil.java
func podState(status v1.PodStatus) string {
	if status.Phase == v1.PodRunning ||
		status.Phase == v1.PodSucceeded {
		return stateUp
	}

	return stateDown
}

// console represents a container's name and log output.
type console struct {
	Name   string `json:"name"`
	Output string `json:"output"`
}

// GetInstanceConsole returns the "console" of an instance. In the case for Kubernetes,
// a "console" is the logs of a given Pod.
func (cc *Controller) GetInstanceConsole(c *gin.Context) {
	account := c.Param("account")
	namespace := c.Param("location")
	n := c.Param("name")
	a := strings.Split(n, " ")
	kind := a[0]
	name := a[1]

	// If the provider is not kubernetes, fail as we cannot generate a console for
	// other providers yet.
	qProvider := c.Query("provider")
	if qProvider != "kubernetes" {
		clouddriver.Error(c, http.StatusNotImplemented, fmt.Errorf("provider %s console not implemented",
			qProvider))
		return
	}

	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	if strings.Contains(kind, ".") {
		a2 := strings.Split(kind, ".")
		kind = a2[0]
	}
	// If the requested Kubernetes kind is not a Pod, return status not implemented.
	if !strings.EqualFold(kind, "pod") {
		clouddriver.Error(c, http.StatusNotImplemented, fmt.Errorf("kind %s console not implemented",
			kind))
		return
	}
	// Grab the Kubernetes provider with a deefault timeout.
	provider, err := cc.KubernetesProviderWithTimeout(account, time.Second*defaultGetTimeoutSeconds)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}
	// Get the instance.
	instance, err := provider.Client.Get(kind, name, namespace)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
	// Declare a new pod structure and grab all containers
	// and init containers from the pod.
	p := kubernetes.NewPod(instance.Object)
	o := p.Object()
	// Combine the containers and init containers into
	// one object.
	containers := []v1.Container{}
	containers = append(containers, o.Spec.Containers...)
	containers = append(containers, o.Spec.InitContainers...)
	// Declare a wait group for all the concurrent calls
	// to make.
	wg := &sync.WaitGroup{}
	// Increment the wait group count to the total number of
	// containers.
	wg.Add(len(containers))
	// Create a channel of console to send to.
	cCh := make(chan console, len(containers))
	// Grab logs for all containers concurrently. I could not
	// find a way to grab all logs for all of a pod's containers,
	// but this works.
	//
	// Unlike the dynamic client, the Kubernetes clientset
	// does not have any hidden mutex locks and can run requests concurrently.
	for _, container := range containers {
		go getLogs(wg, cCh, provider.Clientset, instance, container)
	}
	// Wait for all concurrent calls to finish.
	wg.Wait()

	// Close the console channel.
	close(cCh)

	// Receive all console logs from the console channel.
	consoles := []console{}
	for console := range cCh {
		consoles = append(consoles, console)
	}
	// Sort console logs by name descending.
	sort.Slice(consoles, func(i, j int) bool {
		return consoles[i].Name < consoles[j].Name
	})

	c.JSON(http.StatusOK, gin.H{"output": consoles})
}

// getLogs grabs the logs from a given Pod container and sends them
// to a channel of logs.
func getLogs(wg *sync.WaitGroup, cc chan console, clientset kubernetes.Clientset,
	pod *unstructured.Unstructured, container v1.Container) {
	defer wg.Done()

	// This make a call to the following endpoint on the Kubernetes API server:
	//
	// - /api/v1/namespaces/:namespace/pods/:pod/log?container=:containerName
	//
	// Since this is a direct call and does not need to do API discovery, it is safe
	// to run this call concurrently.
	output, err := clientset.PodLogs(pod.GetName(), pod.GetNamespace(), container.Name)
	if err != nil {
		// If there was an error, log and return.
		clouddriver.Log(err)
		return
	}

	cc <- console{Name: container.Name, Output: output}
}
