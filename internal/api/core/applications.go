package core

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/client-go/rest"
)

const (
	defaultChanSize           = 100000
	defaultListTimeoutSeconds = 10
	stateUp                   = "Up"
	stateDown                 = "Down"
	statusRunning             = "Running"
	typeKubernetes            = "kubernetes"
)

var (
	errCancelJobNotImplemented = errors.New("cancelJob is not implemented for the Kubernetes provider")
	// serverGroupManagerResources consist of Kubernetes kinds Deployments
	// and ReplicaSets.
	serverGroupManagerResources = []string{
		"deployments",
		"replicaSets",
	}
	// serverGroupResources consist of Kubernetes kinds Pods, ReplicaSets,
	// DaemonSets, StatefulSets, and Services.
	serverGroupResources = []string{
		"pods",
		"replicaSets",
		"daemonSets",
		"statefulSets",
		"services",
	}
	// loadBalancerResources consist of Kubernetes kinds Services,
	// Ingresses, Pods, and ReplicaSets.
	loadBalancerResources = []string{
		"pods",
		"replicaSets",
		"statefulSets",
		"services",
		"ingresses",
	}
)

type Applications []Application

type Application struct {
	Attributes   ApplicationAttributes `json:"attributes"`
	ClusterNames map[string][]string   `json:"clusterNames"`
	Name         string                `json:"name"`
}

type ApplicationAttributes struct {
	Name string `json:"name"`
}

const KeyAllApplications = `AllApplications`

// ListApplications returns a list of applications and their associated
// accounts and clusters.
func (cc *Controller) ListApplications(c *gin.Context) {
	rs, err := cc.SQLClient.ListKubernetesClustersByFields("account_name", "kind", "name", "spinnaker_app")
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	response := Applications{}
	apps := uniqueSpinnakerApps(rs)

	for _, app := range apps {
		application := Application{
			Attributes: ApplicationAttributes{
				Name: app,
			},
			ClusterNames: clusterNamesForSpinnakerApp(app, rs),
			Name:         app,
		}

		response = append(response, application)
	}

	// Sort applications by name descending.
	sort.Slice(response, func(i, j int) bool {
		return response[i].Name < response[j].Name
	})

	c.Set(KeyAllApplications, response)
}

// uniqueSpinnakerApps returns a slice of unique Spinnaker
// applications associated with a given list of kubernetes
// resources.
func uniqueSpinnakerApps(rs []kubernetes.Resource) []string {
	apps := []string{}

	for _, r := range rs {
		if !contains(apps, r.SpinnakerApp) {
			apps = append(apps, r.SpinnakerApp)
		}
	}

	return apps
}

// contains returns true if slice s contains element e.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

// clusterNamesForSpinnakerApp returns a map of Kubernetes provider account names
// to a list of Kubernetes resources in the format `<KIND> <NAME>`.
func clusterNamesForSpinnakerApp(application string, rs []kubernetes.Resource) map[string][]string {
	clusterNames := map[string][]string{}

	for _, r := range rs {
		if r.SpinnakerApp == application {
			if _, ok := clusterNames[r.AccountName]; !ok {
				clusterNames[r.AccountName] = []string{}
			}

			resources := clusterNames[r.AccountName]
			resources = append(resources, fmt.Sprintf("%s %s", r.Kind, r.Name))
			clusterNames[r.AccountName] = resources
		}
	}

	return clusterNames
}

type ServerGroupManagers []ServerGroupManager

// ServerGroupManager is a Kubernetes kind "Deployment".
type ServerGroupManager struct {
	Account       string                          `json:"account"`
	APIVersion    string                          `json:"apiVersion"`
	CloudProvider string                          `json:"cloudProvider"`
	CreatedTime   int64                           `json:"createdTime"`
	Kind          string                          `json:"kind"`
	Labels        map[string]string               `json:"labels"`
	Moniker       Moniker                         `json:"moniker"`
	Name          string                          `json:"name"`
	DisplayName   string                          `json:"displayName"`
	Namespace     string                          `json:"namespace"`
	Region        string                          `json:"region"`
	ServerGroups  []ServerGroupManagerServerGroup `json:"serverGroups"`
}

type Key struct {
	Account        string `json:"account"`
	Group          string `json:"group"`
	KubernetesKind string `json:"kubernetesKind"`
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Provider       string `json:"provider"`
}

type Moniker struct {
	App     string `json:"app"`
	Cluster string `json:"cluster"`
}

type ServerGroupManagerServerGroup struct {
	Account   string                               `json:"account"`
	Moniker   ServerGroupManagerServerGroupMoniker `json:"moniker"`
	Name      string                               `json:"name"`
	Namespace string                               `json:"namespace"`
	Region    string                               `json:"region"`
}

type ServerGroupManagerServerGroupMoniker struct {
	App      string `json:"app"`
	Cluster  string `json:"cluster"`
	Sequence int    `json:"sequence"`
}

// ListServerGroupManagers returns a list of Kubernetes
// Deployments and their associated ReplicaSets for a given
// Spinnaker application.
func (cc *Controller) ListServerGroupManagers(c *gin.Context) {
	response := ServerGroupManagers{}
	application := c.Param("application")
	// List all request resources for the given application.
	rs, err := cc.listApplicationResources(c, serverGroupManagerResources, application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)

		return
	}
	// Declare slices to hold resources.
	deployments := filterResourcesByKind(rs, "deployment")
	replicaSets := filterResourcesByKind(rs, "replicaSet")
	// Make a server group manager map of the replicaSets.
	serverGroupManagerMap := makeServerGroupManagerMap(replicaSets, application)
	// Create a new server group manager for each deployment.
	for _, deployment := range deployments {
		sgm := newServerGroupManager(deployment.u, deployment.account, application)
		sgm.ServerGroups = []ServerGroupManagerServerGroup{}

		uid := string(deployment.u.GetUID())
		if v, ok := serverGroupManagerMap[uid]; ok {
			sgm.ServerGroups = v
		}

		response = append(response, sgm)
	}
	// Sort by account (cluster), then namespace, then kind, then name.
	sort.Slice(response, func(i, j int) bool {
		if response[i].Account != response[j].Account {
			return response[i].Account < response[j].Account
		}

		if response[i].Namespace != response[j].Namespace {
			return response[i].Namespace < response[j].Namespace
		}

		if response[i].Kind != response[j].Kind {
			return response[i].Kind < response[j].Kind
		}

		return response[i].Name < response[j].Name
	})

	c.JSON(http.StatusOK, response)
}

// resource represents a Kubernetes resource and holds
// it's unstructured object and associated account.
type resource struct {
	account string
	u       unstructured.Unstructured
}

// filterResourcesByKind filters a resource slice by the given kind and returns the
// resulting slice.
func filterResourcesByKind(rs []resource, kind string) []resource {
	var filtered []resource

	for _, r := range rs {
		if strings.EqualFold(r.u.GetKind(), kind) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// makeServerGroupManagerMap returns a map of a server group manager's (Deployment)
// UID to a list of ReplicaSets that the Deployment owns.
func makeServerGroupManagerMap(replicaSets []resource, application string) map[string][]ServerGroupManagerServerGroup {
	// Map of server group manager's UID to replica sets.
	serverGroupManagerMap := map[string][]ServerGroupManagerServerGroup{}
	// Loop through each pod.
	for _, replicaSet := range replicaSets {
		// Loop through each replica set's owner reference.
		for _, ownerReference := range replicaSet.u.GetOwnerReferences() {
			uid := string(ownerReference.UID)
			if uid == "" {
				continue
			}
			// Build the server group.
			annotations := replicaSet.u.GetAnnotations()
			sequence := sequence(annotations)
			s := newServerGroupManagerServerGroup(replicaSet.u, replicaSet.account,
				application, ownerReference.Name, sequence)
			// Append the server group to the list of server groups at the manager's UID.
			serverGroupManagerMap[uid] = append(serverGroupManagerMap[uid], s)
		}
	}

	return serverGroupManagerMap
}

// newServerGroupManagerServerGroup returns a generated instance of ServerGroupManagerServerGroup.
func newServerGroupManagerServerGroup(replicaSet unstructured.Unstructured, account,
	application, ownerReferenceName string, sequence int) ServerGroupManagerServerGroup {
	s := ServerGroupManagerServerGroup{
		Account: account,
		Moniker: ServerGroupManagerServerGroupMoniker{
			App:      application,
			Cluster:  fmt.Sprintf("%s %s", "deployment", ownerReferenceName),
			Sequence: sequence,
		},
		Name:      fmt.Sprintf("%s %s", "replicaSet", replicaSet.GetName()),
		Namespace: replicaSet.GetNamespace(),
		Region:    replicaSet.GetNamespace(),
	}

	return s
}

// newServerGroupManager returns an instance of ServerGroupManager.
func newServerGroupManager(deployment unstructured.Unstructured,
	account, application string) ServerGroupManager {
	return ServerGroupManager{
		Account:       account,
		APIVersion:    deployment.GetAPIVersion(),
		CloudProvider: typeKubernetes,
		CreatedTime:   deployment.GetCreationTimestamp().Unix() * 1000,
		Kind:          "deployment",
		Labels:        deployment.GetLabels(),
		Moniker: Moniker{
			App:     application,
			Cluster: fmt.Sprintf("%s %s", "deployment", deployment.GetName()),
		},
		Name:        fmt.Sprintf("%s %s", "deployment", deployment.GetName()),
		DisplayName: deployment.GetName(),
		Namespace:   deployment.GetNamespace(),
		Region:      deployment.GetNamespace(),
	}
}

// LoadBalancers is a slice of LoadBalancer.
type LoadBalancers []LoadBalancer

// LoadBalancer represents Kubernetes kinds Service and Ingress.
type LoadBalancer struct {
	Account       string                    `json:"account"`
	Apiversion    string                    `json:"apiVersion"`
	CloudProvider string                    `json:"cloudProvider"`
	CreatedTime   int64                     `json:"createdTime,omitempty"`
	DisplayName   string                    `json:"displayName"`
	Kind          string                    `json:"kind,omitempty"`
	Labels        map[string]string         `json:"labels,omitempty"`
	Moniker       Moniker                   `json:"moniker"`
	Name          string                    `json:"name"`
	Namespace     string                    `json:"namespace"`
	Region        string                    `json:"region"`
	ServerGroups  []LoadBalancerServerGroup `json:"serverGroups"`
	Type          string                    `json:"type"`
}

// LoadBalancerServer groups are ReplicaSets that are fronted by the LoadBalancer.
type LoadBalancerServerGroup struct {
	Account           string                 `json:"account"`
	Cloudprovider     string                 `json:"cloudProvider"`
	Detachedinstances []LoadBalancerInstance `json:"detachedInstances"`
	Instances         []LoadBalancerInstance `json:"instances"`
	Isdisabled        bool                   `json:"isDisabled"`
	Name              string                 `json:"name"`
	Region            string                 `json:"region"`
}

// LoadBalancerInstance represents Pods in a ReplicaSet fronted by a LoadBalancer.
type LoadBalancerInstance struct {
	Health LoadBalancerInstanceHealth `json:"health"`
	ID     string                     `json:"id"`
	Name   string                     `json:"name"`
	Zone   string                     `json:"zone"`
}

// LoadBalancerInstanceHealth is the health of a Pod.
type LoadBalancerInstanceHealth struct {
	Platform string `json:"platform"`
	Source   string `json:"source"`
	State    string `json:"state"`
	Type     string `json:"type"`
}

// ListLoadBalancers lists kubernetes "ingresses" and "services".
func (cc *Controller) ListLoadBalancers(c *gin.Context) {
	response := LoadBalancers{}
	application := c.Param("application")
	// List all request resources for the given application.
	rs, err := cc.listApplicationResources(c, loadBalancerResources, application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)

		return
	}
	// Declare slices to hold resources.
	services := filterResourcesByKind(rs, "service")
	ingresses := filterResourcesByKind(rs, "ingress")
	replicaSets := filterResourcesByKind(rs, "replicaSet")
	statefulSets := filterResourcesByKind(rs, "statefulSet")
	pods := filterResourcesByKind(rs, "pod")
	// Make a map of resources to Pods they own.
	instances := makeLoadBalancerInstanceMap(pods)
	// Make a map of ReplicaSet/StatefulSet UIDs to Services that front them.
	frontableServerGroups := []resource{}
	frontableServerGroups = append(frontableServerGroups, replicaSets...)
	frontableServerGroups = append(frontableServerGroups, statefulSets...)
	loadBalancerServerGroups := makeLoadBalancerServerGroupsMap(frontableServerGroups, services, instances)

	// Create a new Load Balancer for each Service.
	for _, r := range services {
		lb := newLoadBalancer(r.u, r.account, application)
		// If this Service fronts some set of Server Groups
		// (ReplicaSets), then set these as the Load Balancer's
		// Server Groups.
		uid := string(r.u.GetUID())
		if _, ok := loadBalancerServerGroups[uid]; ok {
			lb.ServerGroups = loadBalancerServerGroups[uid]
		}

		response = append(response, lb)
	}
	// Create a new Load Balancer for each Ingress.
	for _, i := range ingresses {
		lb := newLoadBalancer(i.u, i.account, application)
		response = append(response, lb)
	}

	// Sort by account (cluster), then region (namespace), then kind, then name.
	sort.Slice(response, func(i, j int) bool {
		if response[i].Account != response[j].Account {
			return response[i].Account < response[j].Account
		}

		if response[i].Region != response[j].Region {
			return response[i].Region < response[j].Region
		}

		if response[i].Kind != response[j].Kind {
			return response[i].Kind < response[j].Kind
		}

		return response[i].Name < response[j].Name
	})

	c.JSON(http.StatusOK, response)
}

// makeLoadBalancerInstanceMap returns a map of resource UIDs to
// a list of instances (pods) they own.
func makeLoadBalancerInstanceMap(pods []resource) map[string][]LoadBalancerInstance {
	instances := map[string][]LoadBalancerInstance{}

	for _, r := range pods {
		p := kubernetes.NewPod(r.u.Object)

		state := stateUp
		if p.Object().Status.Phase != statusRunning {
			state = stateDown
		}

		or := r.u.GetOwnerReferences()
		for _, o := range or {
			i := LoadBalancerInstance{
				Health: LoadBalancerInstanceHealth{
					Platform: "platform",
					// TODO get the container name that is fronted by the service.
					Source: fmt.Sprintf("Container %s", "TODO"),
					State:  state,
					Type:   "kubernetes/container",
				},
				ID:   string(r.u.GetUID()),
				Name: fmt.Sprintf("pod %s", r.u.GetName()),
				Zone: r.u.GetNamespace(),
			}
			uid := string(o.UID)
			instances[uid] = append(instances[uid], i)
		}
	}

	return instances
}

// makeLoadBalancerServerGroupsMap generates a map of Service UIDs
// to resources fronted by that service.
func makeLoadBalancerServerGroupsMap(serverGroups, services []resource,
	serverGroupInstances map[string][]LoadBalancerInstance) map[string][]LoadBalancerServerGroup {
	loadBalancerServerGroups := map[string][]LoadBalancerServerGroup{}
	// Loop through the resources and find the matching labels.
	for _, serverGroup := range serverGroups {
		// Define the resource and get the pod template labels.
		// Only certain kinds of resources can be fronted by
		// a service.
		kind := serverGroup.u.GetKind()
		if !strings.EqualFold(kind, "replicaSet") &&
			!strings.EqualFold(kind, "statefulSet") {
			continue
		}

		labels, found, err := unstructured.NestedStringMap(serverGroup.u.Object, "spec", "template", "metadata", "labels")
		if err != nil || !found {
			continue
		}
		// Loop through the services and check if a service
		// is fronting this server group.
		for _, service := range services {
			// If the resource and Service are not in the same
			// namespace or cluster then skip.
			if serverGroup.u.GetNamespace() != service.u.GetNamespace() ||
				serverGroup.account != service.account {
				continue
			}
			// Define the Service and get the selector.
			selector, found, err := unstructured.NestedStringMap(service.u.Object, "spec", "selector")
			if err != nil || !found {
				continue
			}
			// If there are no selectors, continue.
			if len(selector) == 0 {
				continue
			}
			// Define if the current resource is "fronted" by the service.
			// If the number of label key/value pairs matches that of the
			// Service's selector, then it is fronted.
			matching := 0
			// Loop through the selectors. A Service only fronts
			// pods owned by a resource if *all* selector key/value pairs
			// are present in the pod's labels.
			for k, v := range selector {
				// If the selector key is not a label key then this Pod Template
				// is not fronted by the service.
				if _, ok := labels[k]; !ok {
					break
				} else if labels[k] == v {
					matching++
				}
			}
			// If the number of matching labels in the resource's Pod Template
			// label's equals the number of selectors in the Service, then
			// this resource is fronted by the Service.
			if len(selector) == matching {
				sg := LoadBalancerServerGroup{
					Account:           serverGroup.account,
					Cloudprovider:     typeKubernetes,
					Detachedinstances: []LoadBalancerInstance{},
					Instances:         []LoadBalancerInstance{},
					Isdisabled:        false,
					Name: fmt.Sprintf("%s %s",
						lowercaseFirst(serverGroup.u.GetKind()), serverGroup.u.GetName()),
					Region: serverGroup.u.GetNamespace(),
				}
				// If there are Pod instances associated with this resource,
				// assign them to the LB Server Group here.
				uid := string(serverGroup.u.GetUID())
				if _, ok := serverGroupInstances[uid]; ok {
					sg.Instances = serverGroupInstances[uid]
				}

				serviceUID := string(service.u.GetUID())
				loadBalancerServerGroups[serviceUID] = append(loadBalancerServerGroups[serviceUID], sg)
			}
		}
	}

	return loadBalancerServerGroups
}

// lowercaseFirst lowercases the first letter of a string.
func lowercaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}

// newLoadBalancer returns an instance of LoadBalancer.
func newLoadBalancer(u unstructured.Unstructured, account, application string) LoadBalancer {
	kind := lowercaseFirst(u.GetKind())

	return LoadBalancer{
		Account:       account,
		Apiversion:    u.GetAPIVersion(),
		CloudProvider: typeKubernetes,
		CreatedTime:   u.GetCreationTimestamp().Unix() * 1000,
		DisplayName:   u.GetName(),
		Kind:          kind,
		Labels:        u.GetLabels(),
		Moniker: Moniker{
			App:     application,
			Cluster: fmt.Sprintf("%s %s", kind, u.GetName()),
		},
		Name:         fmt.Sprintf("%s %s", kind, u.GetName()),
		Namespace:    u.GetNamespace(),
		Region:       u.GetNamespace(),
		ServerGroups: []LoadBalancerServerGroup{},
		Type:         typeKubernetes,
	}
}

type Clusters map[string][]string

// ListClusters returns a list of clusters for a given application,
// which for kubernetes is a map of provider names to kubernetes deployment
// kinds and names.
//
// Clusters are kinds deployment, statefulSet, replicaSet, ingress, service, and daemonSet.
func (cc *Controller) ListClusters(c *gin.Context) {
	application := c.Param("application")

	rs, err := cc.SQLClient.ListKubernetesClustersByApplication(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	response := Clusters{}

	for _, resource := range rs {
		if resource.Cluster != "" {
			if _, ok := response[resource.AccountName]; !ok {
				response[resource.AccountName] = []string{}
			}

			kr := response[resource.AccountName]
			kr = append(kr, resource.Cluster)
			response[resource.AccountName] = kr
		}
	}

	c.JSON(http.StatusOK, response)
}

type ServerGroups []ServerGroup

type ServerGroup struct {
	Account        string            `json:"account"`
	AccountName    string            `json:"accountName"`
	BuildInfo      BuildInfo         `json:"buildInfo"`
	Capacity       Capacity          `json:"capacity"`
	CloudProvider  string            `json:"cloudProvider"`
	Cluster        string            `json:"cluster,omitempty"`
	CreatedTime    int64             `json:"createdTime"`
	Disabled       bool              `json:"disabled"`
	DisplayName    string            `json:"displayName"`
	InstanceCounts InstanceCounts    `json:"instanceCounts"`
	Instances      []Instance        `json:"instances"`
	IsDisabled     bool              `json:"isDisabled"`
	Key            Key               `json:"key"`
	Kind           string            `json:"kind"`
	Labels         map[string]string `json:"labels"`
	// LaunchConfig struct {} `json:"launchConfig"`
	LoadBalancers       []string                        `json:"loadBalancers"`
	Manifest            map[string]interface{}          `json:"manifest"`
	Moniker             ServerGroupMoniker              `json:"moniker"`
	Name                string                          `json:"name"`
	Namespace           string                          `json:"namespace"`
	ProviderType        string                          `json:"providerType"`
	Region              string                          `json:"region"`
	SecurityGroups      []interface{}                   `json:"securityGroups"`
	ServerGroupManagers []ServerGroupServerGroupManager `json:"serverGroupManagers"`
	Type                string                          `json:"type"`
	UID                 string                          `json:"uid"`
	Zone                string                          `json:"zone"`
	Zones               []interface{}                   `json:"zones"`
	InsightActions      []interface{}                   `json:"insightActions"`
}

type ServerGroupServerGroupManager struct {
	Account  string `json:"account"`
	Location string `json:"location"`
	Name     string `json:"name"`
}

type ServerGroupMoniker struct {
	App      string `json:"app"`
	Cluster  string `json:"cluster"`
	Sequence int    `json:"sequence"`
}

type BuildInfo struct {
	Images []string `json:"images"`
}

type Capacity struct {
	Desired int32 `json:"desired"`
	Pinned  bool  `json:"pinned"`
}

type InstanceCounts struct {
	Down         int   `json:"down"`
	OutOfService int   `json:"outOfService"`
	Starting     int   `json:"starting"`
	Total        int32 `json:"total"`
	Unknown      int   `json:"unknown"`
	Up           int32 `json:"up"`
}

// Instance if a Kuberntes kind "Pod".
type Instance struct {
	Account           string                 `json:"account,omitempty"`
	AccountName       string                 `json:"accountName,omitempty"`
	AvailabilityZone  string                 `json:"availabilityZone,omitempty"`
	CloudProvider     string                 `json:"cloudProvider,omitempty"`
	CreatedTime       int64                  `json:"createdTime,omitempty"`
	Health            []InstanceHealth       `json:"health,omitempty"`
	HealthState       string                 `json:"healthState,omitempty"`
	HumanReadableName string                 `json:"humanReadableName,omitempty"`
	ID                string                 `json:"id,omitempty"`
	Key               Key                    `json:"key,omitempty"`
	Kind              string                 `json:"kind,omitempty"`
	Labels            map[string]string      `json:"labels,omitempty"`
	Manifest          map[string]interface{} `json:"manifest,omitempty"`
	Moniker           Moniker                `json:"moniker,omitempty"`
	Name              string                 `json:"name,omitempty"`
	ProviderType      string                 `json:"providerType,omitempty"`
	Region            string                 `json:"region,omitempty"`
	Type              string                 `json:"type,omitempty"`
	UID               string                 `json:"uid,omitempty"`
	Zone              string                 `json:"zone,omitempty"`
}

type InstanceHealth struct {
	Platform string `json:"platform,omitempty"`
	Source   string `json:"source,omitempty"`
	State    string `json:"state"`
	Type     string `json:"type"`
}

// ListServerGroups returns a list of Kubernetes kinds ReplicaSets, DaemonSets,
// StatefulSets and their associated Pods.
func (cc *Controller) ListServerGroups(c *gin.Context) {
	response := ServerGroups{}
	application := c.Param("application")
	// List all request resources for the given application.
	rs, err := cc.listApplicationResources(c, serverGroupResources, application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)

		return
	}
	// Declare slices to hold resources.
	pods := filterResourcesByKind(rs, "pod")
	replicaSets := filterResourcesByKind(rs, "replicaSet")
	daemonSets := filterResourcesByKind(rs, "daemonSet")
	statefulSets := filterResourcesByKind(rs, "statefulSet")
	services := filterResourcesByKind(rs, "service")
	// Make a map of a Pod's owner reference to the list of pods
	// it owns.
	serverGroupMap := makeServerGroupMap(pods)
	// Make a map of ReplicaSet/StatefulSet UIDs to Services that front them.
	frontableServerGroups := []resource{}
	frontableServerGroups = append(frontableServerGroups, replicaSets...)
	frontableServerGroups = append(frontableServerGroups, statefulSets...)
	serverGroupLoadBalancers := makeServerGroupLoadBalancersMap(frontableServerGroups, services)
	// Combine the resources into one server group slice.
	serverGroups := []resource{}
	serverGroups = append(serverGroups, replicaSets...)
	serverGroups = append(serverGroups, daemonSets...)
	serverGroups = append(serverGroups, statefulSets...)

	// Create a new server group for each of these resources.
	for _, sg := range serverGroups {
		_sg := newServerGroup(sg.u, serverGroupMap, sg.account)
		if _, ok := serverGroupLoadBalancers[string(sg.u.GetUID())]; ok {
			_sg.LoadBalancers = serverGroupLoadBalancers[string(sg.u.GetUID())]
		}

		response = append(response, _sg)
	}
	// Sort by account (cluster), then namespace, then kind, then name.
	sort.Slice(response, func(i, j int) bool {
		if response[i].Account != response[j].Account {
			return response[i].Account < response[j].Account
		}

		if response[i].Namespace != response[j].Namespace {
			return response[i].Namespace < response[j].Namespace
		}

		if response[i].Kind != response[j].Kind {
			return response[i].Kind < response[j].Kind
		}

		return response[i].Name < response[j].Name
	})

	c.JSON(http.StatusOK, response)
}

// makeServerGroupLoadBalancersMap returns a map of resource UIDs to a slice
// of Service names that front the reource.
func makeServerGroupLoadBalancersMap(serverGroups, services []resource) map[string][]string {
	serverGroupLoadBalancers := map[string][]string{}
	// Loop through the resources and find the matching labels.
	for _, serverGroup := range serverGroups {
		// Define the resource and get the pod template labels.
		// Only certain kinds of resources can be fronted by
		// a service.
		kind := serverGroup.u.GetKind()
		if !strings.EqualFold(kind, "replicaSet") &&
			!strings.EqualFold(kind, "statefulSet") {
			continue
		}

		labels, found, err := unstructured.NestedStringMap(serverGroup.u.Object,
			"spec", "template", "metadata", "labels")
		if err != nil || !found {
			continue
		}
		// Loop through the services and check if a service
		// is fronting this server group.
		for _, service := range services {
			// If the resource and Service are not in the same
			// namespace or cluster then skip.
			if serverGroup.u.GetNamespace() != service.u.GetNamespace() ||
				serverGroup.account != service.account {
				continue
			}
			// Define the Service and get the selector.
			selector, found, err := unstructured.NestedStringMap(service.u.Object, "spec", "selector")
			if err != nil || !found {
				continue
			}
			// If there are no selectors, continue.
			if len(selector) == 0 {
				continue
			}
			// Define if the current resource is "fronted" by the service.
			// If the number of label key/value pairs matches that of the
			// Service's selector, then it is fronted.
			matching := 0
			// Loop through the selectors. A Service only fronts
			// pods owned by a resource if *all* selector key/value pairs
			// are present in the pod's labels.
			for k, v := range selector {
				// If the selector key is not a label key then this Pod Template
				// is not fronted by the service.
				if _, ok := labels[k]; !ok {
					break
				} else if labels[k] == v {
					matching++
				}
			}
			// If the number of matching labels in the resource's Pod Template
			// label's equals the number of selectors in the Service, then
			// this resource is fronted by the Service.
			if len(selector) == matching {
				uid := string(serverGroup.u.GetUID())
				serverGroupLoadBalancers[uid] =
					append(serverGroupLoadBalancers[uid], fmt.Sprintf("service %s", service.u.GetName()))
			}
		}
	}

	return serverGroupLoadBalancers
}

// makeServerGroupMap returns a map of a server group's (replicaSet, daemonSet, statefulSet)
// UID to a list of instances (pods) that this server group owns.
func makeServerGroupMap(pods []resource) map[string][]Instance {
	// Map of server group to instances (pods)
	serverGroupMap := map[string][]Instance{}
	// Sort the pods.
	sort.Slice(pods, func(i, j int) bool {
		if pods[i].u.GetNamespace() != pods[j].u.GetNamespace() {
			return pods[i].u.GetNamespace() < pods[j].u.GetNamespace()
		}

		return pods[i].u.GetName() < pods[j].u.GetName()
	})
	// Loop through each pod.
	for _, pod := range pods {
		// Loop through each pod's owner reference.
		for _, ownerReference := range pod.u.GetOwnerReferences() {
			uid := string(ownerReference.UID)
			if uid == "" {
				continue
			}

			serverGroupMap[uid] = append(serverGroupMap[uid], newInstance(pod.u))
		}
	}

	return serverGroupMap
}

// newInstance returns an "Instance" object from a given
// pod struct.
func newInstance(pod unstructured.Unstructured) Instance {
	state := stateUp

	p := kubernetes.NewPod(pod.Object)
	if p.Object().Status.Phase != statusRunning {
		state = stateDown
	}

	instance := Instance{
		AvailabilityZone: pod.GetNamespace(),
		Health: []InstanceHealth{
			{
				State: state,
				Type:  "kubernetes/pod",
			},
			{
				State: state,
				Type:  "kubernetes/container",
			},
		},
		HealthState: state,
		ID:          string(pod.GetUID()),
		Name:        fmt.Sprintf("%s %s", "pod", pod.GetName()),
	}

	return instance
}

// newServerGroup builds an instance of ServerGroup, which is of Kubernetes kind ReplicaSet, DaemonSet, or StatefulSet.
// It references the given resources owner reference to determine which resource owns it (for example, a ReplicaSet
// is owned by a given Deployment).
func newServerGroup(result unstructured.Unstructured, serverGroupMap map[string][]Instance, account string) ServerGroup {
	images := listImages(&result)
	desired := getDesiredReplicasCount(&result)

	serverGroupManagers := []ServerGroupServerGroupManager{}
	instances := []Instance{}
	annotations := result.GetAnnotations()
	// Get the instances from the instance map.
	uid := string(result.GetUID())
	if v, ok := serverGroupMap[uid]; ok {
		instances = v
	}

	// Build server group manager.
	ownerReferences := result.GetOwnerReferences()
	for _, ownerReference := range ownerReferences {
		// If the owner of the server group is a deployment.
		if strings.EqualFold(ownerReference.Kind, "deployment") {
			// Define a new server group manager from the owner reference.
			sgm := ServerGroupServerGroupManager{
				Account:  account,
				Location: result.GetNamespace(),
				Name:     ownerReference.Name,
			}
			serverGroupManagers = append(serverGroupManagers, sgm)
		}
	}

	cluster := annotations["moniker.spinnaker.io/cluster"]
	app := annotations["moniker.spinnaker.io/application"]
	sequence := sequence(annotations)

	return ServerGroup{
		Account: account,
		BuildInfo: BuildInfo{
			Images: images,
		},
		Capacity: Capacity{
			Desired: desired,
			Pinned:  false,
		},
		CloudProvider: typeKubernetes,
		Cluster:       cluster,
		CreatedTime:   result.GetCreationTimestamp().Unix() * 1000,
		// Include for sorting.
		Kind: lowercaseFirst(result.GetKind()),
		InstanceCounts: InstanceCounts{
			Down:         0,
			OutOfService: 0,
			Starting:     0,
			Total:        getTotalReplicasCount(&result),
			Unknown:      0,
			Up:           getReadyReplicasCount(&result),
		},
		Instances:     instances,
		IsDisabled:    false,
		DisplayName:   result.GetName(),
		LoadBalancers: nil,
		Moniker: ServerGroupMoniker{
			App:      app,
			Cluster:  cluster,
			Sequence: sequence,
		},
		Name:                fmt.Sprintf("%s %s", lowercaseFirst(result.GetKind()), result.GetName()),
		Namespace:           result.GetNamespace(),
		Region:              result.GetNamespace(),
		SecurityGroups:      nil,
		ServerGroupManagers: serverGroupManagers,
		Type:                typeKubernetes,
		Labels:              result.GetLabels(),
	}
}

// sequence returns the sequence of a given resource.
// A versioned resource contains its sequence in the
// `moniker.spinnaker.io/sequence` annotation.
// A resource which is owned by some deployment defines its sequence in
// the `deployment.kubernetes.io/revision` annotation.
func sequence(annotations map[string]string) int {
	if annotations == nil {
		return 0
	}

	if _, ok := annotations["moniker.spinnaker.io/sequence"]; ok {
		sequence, _ := strconv.Atoi(annotations["moniker.spinnaker.io/sequence"])
		return sequence
	}

	sequence, _ := strconv.Atoi(annotations["deployment.kubernetes.io/revision"])

	return sequence
}

// List images for replicaSets, statefulSets, and daemonSets.
func listImages(result *unstructured.Unstructured) []string {
	images := []string{}

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		o := rs.Object()

		for _, container := range o.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}
	case "daemonset":
		ds := kubernetes.NewDaemonSet(result.Object)
		o := ds.Object()

		for _, container := range o.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}
	case "statefulset":
		sts := kubernetes.NewStatefulSet(result.Object)

		o := sts.Object()
		for _, container := range o.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}
	}

	return images
}

// getDesiredReplicasCount returns the desired replicas for
// replicaSets, statefulSets, and daemonSets.
func getDesiredReplicasCount(result *unstructured.Unstructured) int32 {
	desired := int32(0)

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		if rs.Object().Spec.Replicas != nil {
			desired = *rs.Object().Spec.Replicas
		}
	case "daemonset":
		ds := kubernetes.NewDaemonSet(result.Object)
		o := ds.Object()
		desired = o.Status.DesiredNumberScheduled
	case "statefulset":
		sts := kubernetes.NewStatefulSet(result.Object)
		o := sts.Object()

		if o.Spec.Replicas != nil {
			desired = *o.Spec.Replicas
		}
	}

	return desired
}

// getTotalReplicasCount returns total desired replicas for
// replicaSets, statefulSets, and daemonSets.
func getTotalReplicasCount(result *unstructured.Unstructured) int32 {
	total := int32(0)

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		total = rs.Object().Status.Replicas
	case "daemonset":
		ds := kubernetes.NewDaemonSet(result.Object)
		o := ds.Object()
		total = o.Status.DesiredNumberScheduled
	case "statefulset":
		sts := kubernetes.NewStatefulSet(result.Object)
		o := sts.Object()
		total = o.Status.Replicas
	}

	return total
}

// getTotalReplicasCount returns total replicas in a ready state for replicaSets,
// statefulSets, and daemonSets.
func getReadyReplicasCount(result *unstructured.Unstructured) int32 {
	ready := int32(0)

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		if rs.Object().Spec.Replicas != nil {
			ready = rs.Object().Status.ReadyReplicas
		}
	case "daemonset":
		ds := kubernetes.NewDaemonSet(result.Object)
		o := ds.Object()
		ready = o.Status.NumberReady
	case "statefulset":
		sts := kubernetes.NewStatefulSet(result.Object)
		o := sts.Object()
		ready = o.Status.ReadyReplicas
	}

	return ready
}

// GetServerGroup returns a specific server group (Kubernetes kind ReplicaSet,
// DaemonSet, or StatefulSet) for a given cluster, namespace, name and Spinnaker application.
// This endpoint is called when clicking on a given resource in the "Clusters" tab in Deck.
func (cc *Controller) GetServerGroup(c *gin.Context) {
	account := c.Param("account")
	application := c.Param("application")
	location := c.Param("location")
	nameArray := strings.Split(c.Param("name"), " ")
	kind := nameArray[0]
	name := nameArray[1]

	client, err := cc.kubeConfigClient(c.Copy(), account)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	lo := metav1.ListOptions{
		LabelSelector: kubernetes.LabelKubernetesName + "=" + application,
	}

	result, err := client.Get(kind, name, location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defaultListTimeoutSeconds)
	defer cancel()
	// "Instances" in kubernetes are pods.
	pods, err := client.ListResourceWithContext(ctx, "pods", lo)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	instanceCounts := InstanceCounts{}
	images := listImages(result)
	desired := getDesiredReplicasCount(result)
	instanceCounts.Total = getTotalReplicasCount(result)
	instanceCounts.Up = getReadyReplicasCount(result)
	instances := []Instance{}

	for _, v := range pods.Items {
		p := kubernetes.NewPod(v.Object)
		for _, ownerReference := range p.Object().ObjectMeta.OwnerReferences {
			if ownerReference.UID == result.GetUID() {
				instance := newPodInstance(p, application, account)
				instance.Manifest = v.Object
				instances = append(instances, instance)
			}
		}
	}

	annotations := result.GetAnnotations()
	cluster := annotations["moniker.spinnaker.io/cluster"]
	app := annotations["moniker.spinnaker.io/application"]
	sequence := sequence(annotations)

	if app == "" {
		app = application
	}

	response := ServerGroup{
		Account:     account,
		AccountName: account,
		BuildInfo: BuildInfo{
			Images: images,
		},
		Capacity: Capacity{
			Desired: desired,
			Pinned:  false,
		},
		CloudProvider:  typeKubernetes,
		CreatedTime:    result.GetCreationTimestamp().Unix() * 1000,
		Disabled:       false,
		DisplayName:    result.GetName(),
		InstanceCounts: instanceCounts,
		Instances:      instances,
		Key: Key{
			Account:        account,
			Group:          lowercaseFirst(result.GetKind()),
			KubernetesKind: lowercaseFirst(result.GetKind()),
			Name:           result.GetName(),
			Namespace:      result.GetNamespace(),
			Provider:       typeKubernetes,
		},
		Kind:          lowercaseFirst(result.GetKind()),
		Labels:        result.GetLabels(),
		LoadBalancers: []string{},
		Manifest:      result.Object,
		Moniker: ServerGroupMoniker{
			App:      app,
			Cluster:  cluster,
			Sequence: sequence,
		},
		Name:                fmt.Sprintf("%s %s", lowercaseFirst(result.GetKind()), result.GetName()),
		Namespace:           result.GetNamespace(),
		ProviderType:        typeKubernetes,
		Region:              result.GetNamespace(),
		SecurityGroups:      []interface{}{},
		ServerGroupManagers: []ServerGroupServerGroupManager{},
		Type:                typeKubernetes,
		UID:                 string(result.GetUID()),
		Zone:                result.GetNamespace(),
		Zones:               []interface{}{},
		InsightActions:      []interface{}{},
	}

	c.JSON(http.StatusOK, response)
}

// newPodInstance returns a new instance that represents a kind Kubernetes
// pod.
func newPodInstance(p *kubernetes.Pod, application, account string) Instance {
	state := stateUp
	if p.Object().Status.Phase != statusRunning {
		state = stateDown
	}

	annotations := p.Object().ObjectMeta.Annotations
	cluster := annotations["moniker.spinnaker.io/cluster"]
	app := annotations["moniker.spinnaker.io/application"]

	if app == "" {
		app = application
	}

	instance := Instance{
		Account:          account,
		AccountName:      account,
		AvailabilityZone: p.Object().ObjectMeta.Namespace,
		CloudProvider:    typeKubernetes,
		CreatedTime:      p.Object().ObjectMeta.CreationTimestamp.Unix() * 1000,
		Health: []InstanceHealth{
			{
				State: state,
				Type:  "kubernetes/pod",
			},
			{
				State: state,
				Type:  "kubernetes/container",
			},
		},
		HealthState:       state,
		HumanReadableName: fmt.Sprintf("%s %s", "pod", p.Object().ObjectMeta.Name),
		ID:                string(p.Object().ObjectMeta.UID),
		Key: Key{
			Account:        account,
			Group:          "pod",
			KubernetesKind: "pod",
			Name:           p.Object().ObjectMeta.Name,
			Namespace:      p.Object().ObjectMeta.Namespace,
			Provider:       typeKubernetes,
		},
		Kind:   "pod",
		Labels: p.Object().ObjectMeta.Labels,
		Moniker: Moniker{
			App:     app,
			Cluster: cluster,
		},
		Name:         fmt.Sprintf("%s %s", "pod", p.Object().ObjectMeta.Name),
		ProviderType: typeKubernetes,
		Region:       p.Object().ObjectMeta.Namespace,
		Type:         typeKubernetes,
		UID:          string(p.Object().ObjectMeta.UID),
		Zone:         p.Object().ObjectMeta.Namespace,
	}

	return instance
}

type Job struct {
	Account           string                   `json:"account"`
	CompletionDetails JobCompletionDetails     `json:"completionDetails"`
	CreatedTime       int64                    `json:"createdTime"`
	JobState          string                   `json:"jobState"`
	Location          string                   `json:"location"`
	Name              string                   `json:"name"`
	Pods              []map[string]interface{} `json:"pods"`
	Provider          string                   `json:"provider"`
}

type JobCompletionDetails struct {
	ExitCode string `json:"exitCode"`
	Message  string `json:"message"`
	Reason   string `json:"reason"`
	Signal   string `json:"signal"`
}

// GetJob retrieves a given Kubernetes job from a given cluster
// given a namespace and name.
func (cc *Controller) GetJob(c *gin.Context) {
	account := c.Param("account")
	// application := c.Param("application")
	location := c.Param("location")
	nameArray := strings.Split(c.Param("name"), " ")
	kind := nameArray[0]
	name := nameArray[1]

	client, err := cc.kubeConfigClient(c.Copy(), account)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	result, err := client.Get(kind, name, location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	j := kubernetes.NewJob(result.Object)

	// TODO fill in pod definitions.
	job := Job{
		Account: account,
		CompletionDetails: JobCompletionDetails{
			ExitCode: "",
			Message:  "",
			Reason:   "",
			Signal:   "",
		},
		CreatedTime: result.GetCreationTimestamp().Unix() * 1000,
		JobState:    j.State(),
		Location:    location,
		Name:        name,
		Pods:        []map[string]interface{}{},
		Provider:    typeKubernetes,
	}

	c.JSON(http.StatusOK, job)
}

// DeleteJob is not implemented for the Kubernetes provider V2.
// See https://github.com/spinnaker/spinnaker/issues/4644#issuecomment-627287782.
func DeleteJob(c *gin.Context) {
	clouddriver.Error(c, http.StatusInternalServerError, errCancelJobNotImplemented)
}

// listApplicationsResources lists all accounts for a given app, then concurrently lists
// all requested resources for the given app concurrently.
func (cc *Controller) listApplicationResources(c *gin.Context, rs []string, application string) ([]resource, error) {
	wg := &sync.WaitGroup{}
	// Create channel of resouces to send to.
	rc := make(chan resource, defaultChanSize)
	// List all accounts associated with the given Spinnaker app.
	accounts, err := cc.SQLClient.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		return nil, err
	}
	// Add the number of accounts to the wait group.
	wg.Add(len(accounts))
	// List all requested resources across accounts concurrently.
	for _, account := range accounts {
		go cc.listResources(c.Copy(), wg, rs, rc, account, application)
	}
	// Wait for all concurrent calls to finish.
	wg.Wait()
	// Close the channel.
	close(rc)
	// Receive all resources from the channel.
	resources := []resource{}
	for r := range rc {
		resources = append(resources, r)
	}
	// Return the slice of unstructured resources.
	return resources, nil
}

// listResources initializes discovery for a given client then lists
// the requested resources concurrently.
func (cc *Controller) listResources(c *gin.Context, wg *sync.WaitGroup, rs []string, rc chan resource,
	account, application string) {
	// Increment the wait group counter when we're done here.
	defer wg.Done()
	// Grab the kube client for the given account.
	client, err := cc.kubeConfigClient(c, account)
	if err != nil {
		clouddriver.Log(err)
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
	if err = client.Discover(); err != nil {
		clouddriver.Log(err)
		return
	}
	// Declare a new waitgroup to wait on concurrent resource listing.
	_wg := &sync.WaitGroup{}
	// Add the number of resources we will be listing concurrently.
	_wg.Add(len(rs))
	// List all required resources concurrently.
	for _, r := range rs {
		go list(c.Copy(), _wg, rc, client, r, account, application)
	}
	// Wait for the calls to finish.
	_wg.Wait()
}

// list lists a given resource and send to a channel of unstructured.Unstructured.
// It uses a context with a timeout of 10 seconds.
func list(c *gin.Context, wg *sync.WaitGroup, rc chan resource,
	client kubernetes.Client, r, account, application string) {
	// Finish the wait group when we're done here.
	defer wg.Done()
	// Declare server side filtering options.
	lo := metav1.ListOptions{
		LabelSelector: kubernetes.LabelKubernetesName + "=" + application,
	}
	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defaultListTimeoutSeconds)
	defer cancel()
	// List resources with the context.
	ul, err := client.ListResourceWithContext(ctx, r, lo)
	if err != nil {
		// If there was an error, log and return.
		clouddriver.Log(err)
		return
	}
	// Send all unstructured objects to the channel.
	for _, u := range ul.Items {
		res := resource{
			u:       u,
			account: account,
		}
		rc <- res
	}
}

// kubeConfigClient returns a new Kubernetes client
// for a given account.
func (cc *Controller) kubeConfigClient(c *gin.Context, account string) (kubernetes.Client, error) {
	// Get the provider info for the account.
	provider, err := cc.SQLClient.GetKubernetesProvider(account)
	if err != nil {
		return nil, err
	}
	// Decode the provider's CA data.
	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return nil, err
	}
	// Grab the auth token from arcade.
	token, err := cc.ArcadeClient.Token(provider.TokenProvider)
	if err != nil {
		return nil, err
	}
	// Generate a new rest config using this information.
	// Set the timeout to be the list timeout.
	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
		Timeout: time.Second * defaultListTimeoutSeconds,
	}
	// Create a new dynamic client for this config.
	client, err := cc.KubernetesController.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
