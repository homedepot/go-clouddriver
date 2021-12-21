package core

import (
	"context"
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
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

const (
	stateUp        = "Up"
	stateDown      = "Down"
	statusRunning  = "Running"
	typeKubernetes = "kubernetes"
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

	m := makeMapOfSpinnakerAppsToClusters(rs)
	for app, clusterNames := range m {
		application := Application{
			Attributes: ApplicationAttributes{
				Name: app,
			},
			ClusterNames: clusterNames,
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

// contains returns true if slice s contains element e.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

// containsIgnoreCase returns true if slice s contains element e, ignoring case.
func containsIgnoreCase(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}

	return false
}

func makeMapOfSpinnakerAppsToClusters(rs []kubernetes.Resource) map[string]map[string][]string {
	m := map[string]map[string][]string{}

	for _, r := range rs {
		if _, ok := m[r.SpinnakerApp]; !ok {
			m[r.SpinnakerApp] = map[string][]string{}
		}

		clusterNames := m[r.SpinnakerApp]
		resources := clusterNames[r.AccountName]
		resources = append(resources, fmt.Sprintf("%s %s", r.Kind, r.Name))
		clusterNames[r.AccountName] = resources
		m[r.SpinnakerApp] = clusterNames
	}

	return m
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
	// List all accounts associated with the given Spinnaker app.
	accounts, err := cc.SQLClient.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)

		return
	}
	// List all request resources for the given application.
	rs, err := cc.listApplicationResources(c, serverGroupManagerResources, accounts, []string{application})
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
	account     string
	application string
	u           unstructured.Unstructured
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
	// List all accounts associated with the given Spinnaker app.
	accounts, err := cc.SQLClient.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)

		return
	}
	// List all request resources for the given application.
	rs, err := cc.listApplicationResources(c, loadBalancerResources, accounts, []string{application})
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

type ClusterServerGroups struct {
	AccountName   string                           `json:"accountName"`
	Application   string                           `json:"application"`
	LoadBalancers []string                         `json:"loadBalancers"`
	Moniker       Moniker                          `json:"moniker"`
	Name          string                           `json:"name"`
	ServerGroups  []ClusterServerGroupsServerGroup `json:"serverGroups"`
	Type          string                           `json:"type"`
}

type ClusterServerGroupsServerGroup struct {
	Account    string `json:"account"`
	APIVersion string `json:"apiVersion"`
	// Buildinfo  struct {
	// 	Images []string `json:"images"`
	// } `json:"buildInfo"`
	// Capacity struct {
	// 	Desired int  `json:"desired"`
	// 	Pinned  bool `json:"pinned"`
	// } `json:"capacity"`
	CloudProvider string `json:"cloudProvider"`
	// CreatedTime    int64  `json:"createdTime"`
	// Disabled       bool   `json:"disabled"`
	DisplayName string `json:"displayName"`
	// Instancecounts struct {
	// 	Down         int `json:"down"`
	// 	Outofservice int `json:"outOfService"`
	// 	Starting     int `json:"starting"`
	// 	Total        int `json:"total"`
	// 	Unknown      int `json:"unknown"`
	// 	Up           int `json:"up"`
	// } `json:"instanceCounts"`
	// Instances []struct {
	// 	Account       string `json:"account"`
	// 	Apiversion    string `json:"apiVersion"`
	// 	Cloudprovider string `json:"cloudProvider"`
	// 	Createdtime   int64  `json:"createdTime"`
	// 	Displayname   string `json:"displayName"`
	// 	Health        []struct {
	// 		Platform string `json:"platform"`
	// 		Source   string `json:"source"`
	// 		State    string `json:"state"`
	// 		Type     string `json:"type"`
	// 	} `json:"health"`
	// 	Healthstate       string `json:"healthState"`
	// 	Humanreadablename string `json:"humanReadableName"`
	// 	Kind              string `json:"kind"`
	// 	Labels            struct {
	// 		AppKubernetesIoManagedBy   string `json:"app.kubernetes.io/managed-by"`
	// 		AppKubernetesIoName        string `json:"app.kubernetes.io/name"`
	// 		MonikerSpinnakerIoSequence string `json:"moniker.spinnaker.io/sequence"`
	// 		PodTemplateHash            string `json:"pod-template-hash"`
	// 		Run                        string `json:"run"`
	// 	} `json:"labels"`
	// 	Moniker struct {
	// 		App      string `json:"app"`
	// 		Cluster  string `json:"cluster"`
	// 		Sequence int    `json:"sequence"`
	// 	} `json:"moniker"`
	// 	Name         string `json:"name"`
	// 	Namespace    string `json:"namespace"`
	// 	Providertype string `json:"providerType"`
	// 	Zone         string `json:"zone"`
	// } `json:"instances"`
	Kind   string            `json:"kind"`
	Labels map[string]string `json:"labels"`
	// Launchconfig struct {
	// } `json:"launchConfig"`
	// Loadbalancers []string{} `json:"loadBalancers"`
	Moniker   ServerGroupMoniker `json:"moniker"`
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Region    string             `json:"region"`
	// Securitygroups      []interface{} `json:"securityGroups"`
	ServerGroupManagers []ClusterServerGroupsServerGroupManager `json:"serverGroupManagers"`
	Type                string                                  `json:"type"`
	Zones               []interface{}                           `json:"zones"`
}

type ClusterServerGroupsServerGroupManager struct {
	Account  string `json:"account"`
	Location string `json:"location"`
	Name     string `json:"name"`
}

// ListClustersByName returns a list of clusters for a given application,
// account, and name, where name is something like "deployment my-deployment".
func (cc *Controller) ListClustersByName(c *gin.Context) {
	application := c.Param("application")
	account := c.Param("account")
	clusterName := c.Param("clusterName")

	a := strings.Split(clusterName, " ")
	if len(a) != 2 {
		clouddriver.Error(c, http.StatusBadRequest,
			fmt.Errorf("clusterName parameter must be in the format of 'kind name', got: %s", clusterName))
		return
	}

	// The Kubernetes kind to list is defined by the cluster kind. If
	// the cluster kind is "deployment" we still want to list ReplicaSets.
	kind := a[0]
	if strings.EqualFold(kind, "deployment") {
		kind = "replicaSet"
	}

	provider, err := cc.KubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
	defer cancel()

	lo := metav1.ListOptions{
		LabelSelector: kubernetes.DefaultLabelSelector(),
	}
	// If namespace-scoped account, then only get resources in the namespace.
	if provider.Namespace != nil {
		lo.FieldSelector = "metadata.namespace=" + *provider.Namespace
	}

	ul, err := provider.Client.ListResourceWithContext(ctx, kind, lo)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Filter out all unassociated objects based on the 'moniker.spinnaker.io/cluster' annotation.
	items := kubernetes.FilterOnAnnotation(ul.Items,
		kubernetes.AnnotationSpinnakerMonikerCluster, clusterName)
	// Filter out all unassociated objects based on the 'moniker.spinnaker.io/application' annotation.
	items = kubernetes.FilterOnAnnotation(items,
		kubernetes.AnnotationSpinnakerMonikerApplication, application)

	serverGroups := []ClusterServerGroupsServerGroup{}

	for _, item := range items {
		serverGroup := ClusterServerGroupsServerGroup{
			Account:       account,
			APIVersion:    item.GetAPIVersion(),
			CloudProvider: typeKubernetes,
			DisplayName:   item.GetName(),
			Kind:          item.GetKind(),
			Labels:        item.GetLabels(),
			Moniker: ServerGroupMoniker{
				App:      application,
				Cluster:  clusterName,
				Sequence: sequence(item.GetAnnotations()),
			},
			Name:      fmt.Sprintf("%s %s", lowercaseFirst(item.GetKind()), item.GetName()),
			Namespace: item.GetNamespace(),
			Region:    item.GetNamespace(),
			Type:      typeKubernetes,
			Zones:     []interface{}{},
		}
		serverGroupManagers := []ClusterServerGroupsServerGroupManager{}

		ownerReferences := item.GetOwnerReferences()
		for _, ownerReference := range ownerReferences {
			serverGroupManager := ClusterServerGroupsServerGroupManager{
				Account:  account,
				Location: item.GetNamespace(),
				Name:     ownerReference.Name,
			}
			serverGroupManagers = append(serverGroupManagers, serverGroupManager)
		}

		serverGroup.ServerGroupManagers = serverGroupManagers
		serverGroups = append(serverGroups, serverGroup)
	}

	cg := ClusterServerGroups{
		AccountName:   account,
		Application:   application,
		LoadBalancers: []string{},
		Moniker: Moniker{
			App:     application,
			Cluster: clusterName,
		},
		Name:         clusterName,
		ServerGroups: serverGroups,
		Type:         typeKubernetes,
	}

	c.JSON(http.StatusOK, cg)
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
	// List all accounts associated with the given Spinnaker app.
	accounts, err := cc.SQLClient.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)

		return
	}
	// List all request resources for the given application.
	rs, err := cc.listApplicationResources(c, serverGroupResources, accounts, []string{application})
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

		managedLoadBalancers, err := kubernetes.LoadBalancers(sg.u)
		if err == nil {
			for _, managedLoadBalancer := range managedLoadBalancers {
				a := strings.Split(managedLoadBalancer, " ")
				if len(a) != 2 {
					continue
				}

				// For now, limit the kind of load balancer available to attach to Services.
				kind := a[0]
				if !strings.EqualFold(kind, "service") {
					continue
				}

				if !containsIgnoreCase(_sg.LoadBalancers, managedLoadBalancer) {
					_sg.LoadBalancers = append(_sg.LoadBalancers, managedLoadBalancer)
					_sg.IsDisabled = true
				}
			}
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
			// Add pod to the server group, if it doesn't already contain it.
			if !containsPod(serverGroupMap[uid], pod) {
				serverGroupMap[uid] = append(serverGroupMap[uid], newInstance(pod.u))
			}
		}
	}

	return serverGroupMap
}

// containsPod returns true if the given slice of pod instances
// contains an element with the same UID as the given pod.
func containsPod(instances []Instance, pod resource) bool {
	for _, instance := range instances {
		if instance.ID == string(pod.u.GetUID()) {
			return true
		}
	}

	return false
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

	cluster := annotations[kubernetes.AnnotationSpinnakerMonikerCluster]
	app := annotations[kubernetes.AnnotationSpinnakerMonikerApplication]
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

	if _, ok := annotations[kubernetes.AnnotationSpinnakerMonikerSequence]; ok {
		sequence, _ := strconv.Atoi(annotations[kubernetes.AnnotationSpinnakerMonikerSequence])
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

	if len(nameArray) != 2 {
		clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("name parameter must be in the format of 'kind name', got: %s", c.Param("name")))
		return
	}

	kind := nameArray[0]
	name := nameArray[1]

	provider, err := cc.KubernetesProviderWithTimeout(account, time.Second*internal.DefaultListTimeoutSeconds)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	result, err := provider.Client.Get(kind, name, location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
	defer cancel()
	// Declare a label selector.
	lo := metav1.ListOptions{
		LabelSelector: kubernetes.DefaultLabelSelector(),
		FieldSelector: "metadata.namespace=" + location,
	}
	// "Instances" in kubernetes are pods.
	pods, err := provider.Client.ListResourceWithContext(ctx, "pods", lo)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Filter the results to only the application annotation requested.
	pods.Items = kubernetes.FilterOnAnnotation(pods.Items,
		kubernetes.AnnotationSpinnakerMonikerApplication, application)
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
	cluster := annotations[kubernetes.AnnotationSpinnakerMonikerCluster]
	sequence := sequence(annotations)

	app := annotations[kubernetes.AnnotationSpinnakerMonikerApplication]
	if app == "" {
		app = application
	}

	loadBalancers, disabled := isDisabled(provider, *result)

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
		Disabled:       disabled,
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
		LoadBalancers: loadBalancers,
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

// isDisabled returns the list of load balancers in its annotation
// `traffic.spinnaker.io/load-balancers` and if the server group is fronted by
// all of these load balancers - false (not disabled) if it is, true (disabled)
// if it is not.
func isDisabled(provider *kubernetes.Provider, serverGroup unstructured.Unstructured) ([]string, bool) {
	loadBalancers := []string{}

	managedLoadBalancers, err := kubernetes.LoadBalancers(serverGroup)
	if err != nil || managedLoadBalancers == nil {
		return loadBalancers, false
	}

	labels, found, err := unstructured.NestedStringMap(serverGroup.Object,
		"spec", "template", "metadata", "labels")
	if err != nil || !found {
		return loadBalancers, false
	}

	disabled := false

	for _, managedLoadBalancer := range managedLoadBalancers {
		a := strings.Split(managedLoadBalancer, " ")
		if len(a) != 2 {
			continue
		}

		kind := a[0]
		name := a[1]

		if !strings.EqualFold(kind, "service") {
			continue
		}

		service, err := provider.Client.Get(kind, name, serverGroup.GetNamespace())
		if err != nil {
			continue
		}

		// Only append to the server groups load balancers if we're able to grab
		// it from the cluster.
		loadBalancers = append(loadBalancers, managedLoadBalancer)

		matching := 0
		// Define the Service and get the selector.
		selector, _, _ := unstructured.NestedStringMap(service.Object, "spec", "selector")
		for k, v := range selector {
			// If the selector key is not a label key then this Pod Template
			// is not fronted by the service.
			if _, ok := labels[k]; !ok {
				break
			} else if labels[k] == v {
				matching++
			}
		}

		if len(selector) != matching {
			disabled = true
		}
	}

	return loadBalancers, disabled
}

// newPodInstance returns a new instance that represents a kind Kubernetes
// pod.
func newPodInstance(p *kubernetes.Pod, application, account string) Instance {
	state := stateUp
	if p.Object().Status.Phase != statusRunning {
		state = stateDown
	}

	annotations := p.Object().ObjectMeta.Annotations
	cluster := annotations[kubernetes.AnnotationSpinnakerMonikerCluster]
	app := annotations[kubernetes.AnnotationSpinnakerMonikerApplication]

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

	if len(nameArray) != 2 {
		clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("name parameter must be in the format of 'kind name', got: %s", c.Param("name")))
		return
	}

	kind := nameArray[0]
	name := nameArray[1]

	provider, err := cc.KubernetesProviderWithTimeout(account, time.Second*internal.DefaultListTimeoutSeconds)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	result, err := provider.Client.Get(kind, name, location)
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
func (cc *Controller) listApplicationResources(c *gin.Context, rs, accounts, applications []string) ([]resource, error) {
	wg := &sync.WaitGroup{}
	// Create channel of resouces to send to.
	rc := make(chan resource, internal.DefaultChanSize)
	// Add the number of accounts to the wait group.
	wg.Add(len(accounts))
	// List all requested resources across accounts concurrently.
	for _, account := range accounts {
		go cc.listResources(wg, rs, rc, account, applications)
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
func (cc *Controller) listResources(wg *sync.WaitGroup, rs []string, rc chan resource,
	account string, applications []string) {
	// Increment the wait group counter when we're done here.
	defer wg.Done()
	// Grab the kube provider for the given account.
	provider, err := cc.KubernetesProviderWithTimeout(account, time.Second*internal.DefaultListTimeoutSeconds)
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
	if err = provider.Client.Discover(); err != nil {
		clouddriver.Log(err)
		return
	}
	// Declare a new waitgroup to wait on concurrent resource listing.
	_wg := &sync.WaitGroup{}
	// Add the number of resources we will be listing concurrently.
	_wg.Add(len(rs))
	// List all required resources concurrently.
	for _, r := range rs {
		go list(_wg, rc, provider, r, account, applications)
	}
	// Wait for the calls to finish.
	_wg.Wait()
}

// list lists a given resource and send to a channel of unstructured.Unstructured.
// It uses a context with a timeout of 10 seconds.
func list(wg *sync.WaitGroup, rc chan resource,
	provider *kubernetes.Provider, r, account string, applications []string) {
	// Finish the wait group when we're done here.
	defer wg.Done()
	// Declare server side filtering options.
	lo := metav1.ListOptions{
		LabelSelector: kubernetes.DefaultLabelSelector(),
	}
	// If namespace-scoped account, then only get resources in the namespace.
	if provider.Namespace != nil {
		lo.FieldSelector = "metadata.namespace=" + *provider.Namespace
	}
	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
	defer cancel()
	// List resources with the context.
	ul, err := provider.Client.ListResourceWithContext(ctx, r, lo)
	if err != nil {
		// If there was an error, log and return.
		clouddriver.Log(err)
		return
	}
	// Sometimes the application annotation has a double quote ('"')
	// character prefix and suffix, remove those to make sure we associate
	// the resources correctly. This happens with the Spinnaker Operator, for example.
	for _, item := range ul.Items {
		annotations := item.GetAnnotations()
		if annotations != nil {
			if _, ok := annotations[kubernetes.AnnotationSpinnakerMonikerApplication]; ok {
				a := annotations[kubernetes.AnnotationSpinnakerMonikerApplication]
				a = strings.TrimPrefix(a, "\"")
				a = strings.TrimSuffix(a, "\"")
				annotations[kubernetes.AnnotationSpinnakerMonikerApplication] = a
			}

			item.SetAnnotations(annotations)
		}
	}
	// Filter the results to only the application annotation requested.
	for _, application := range applications {
		items := kubernetes.FilterOnAnnotation(ul.Items,
			kubernetes.AnnotationSpinnakerMonikerApplication, application)
		// Send all unstructured objects to the channel.
		for _, u := range items {
			res := resource{
				u:           u,
				account:     account,
				application: application,
			}
			rc <- res
		}
	}
}
