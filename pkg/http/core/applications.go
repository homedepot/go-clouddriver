package core

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

var (
	// Default to a timeout of 10 seconds on all lists.
	listTimeout = int64(10)
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

func ListApplications(c *gin.Context) {
	sc := sql.Instance(c)

	rs, err := sc.ListKubernetesClustersByFields("account_name", "kind", "name", "spinnaker_app")
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

	c.Set(KeyAllApplications, response)
}

func uniqueSpinnakerApps(rs []kubernetes.Resource) []string {
	apps := []string{}

	for _, r := range rs {
		if !contains(apps, r.SpinnakerApp) {
			apps = append(apps, r.SpinnakerApp)
		}
	}

	return apps
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

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

type ServerGroupManager struct {
	Account       string                          `json:"account"`
	AccountName   string                          `json:"accountName"`
	CloudProvider string                          `json:"cloudProvider"`
	CreatedTime   int64                           `json:"createdTime"`
	Key           Key                             `json:"key"`
	Kind          string                          `json:"kind"`
	Labels        map[string]string               `json:"labels"`
	Manifest      map[string]interface{}          `json:"manifest"`
	Moniker       Moniker                         `json:"moniker"`
	Name          string                          `json:"name"`
	DisplayName   string                          `json:"displayName"`
	Namespace     string                          `json:"namespace"`
	ProviderType  string                          `json:"providerType"`
	Region        string                          `json:"region"`
	ServerGroups  []ServerGroupManagerServerGroup `json:"serverGroups"`
	Type          string                          `json:"type"`
	UID           string                          `json:"uid"`
	Zone          string                          `json:"zone"`
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

// Server Group Managers for a kubernetes target are deployments.
func ListServerGroupManagers(c *gin.Context) {
	sc := sql.Instance(c)
	application := c.Param("application")
	response := ServerGroupManagers{}
	wg := &sync.WaitGroup{}
	sgms := make(chan ServerGroupManager, 100000)

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	wg.Add(len(accounts))

	// Don't actually return while attempting to create a list of server group managers.
	// We want to avoid the situation where a user cannot perform operations when any
	// cluster is not available.
	for _, account := range accounts {
		go listServerGroupManagers(c, wg, sgms, account, application)
	}

	wg.Wait()

	close(sgms)

	for sgm := range sgms {
		response = append(response, sgm)
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Name < response[j].Name
	})

	c.JSON(http.StatusOK, response)
}

func listServerGroupManagers(c *gin.Context, wg *sync.WaitGroup, sgms chan ServerGroupManager,
	account, application string) {
	defer wg.Done()

	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	token, err := ac.Token(provider.TokenProvider)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	lo := metav1.ListOptions{
		LabelSelector:  kubernetes.LabelKubernetesName + "=" + application,
		TimeoutSeconds: &listTimeout,
	}

	deployments, err := client.ListResource("deployments", lo)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	replicaSets, err := client.ListResource("replicaSets", lo)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	for _, deployment := range deployments.Items {
		sgm := newServerGroupManager(deployment, account, application)
		sgm.ServerGroups = buildServerGroups(replicaSets, deployment, account, application)
		sgms <- sgm
	}
}

func newServerGroupManager(deployment unstructured.Unstructured,
	account, application string) ServerGroupManager {
	return ServerGroupManager{
		Account:       account,
		AccountName:   account,
		CloudProvider: "kubernetes",
		CreatedTime:   deployment.GetCreationTimestamp().Unix() * 1000,
		Key: Key{
			Account:        account,
			Group:          "deployment",
			KubernetesKind: "deployment",
			Name:           deployment.GetName(),
			Namespace:      deployment.GetNamespace(),
			Provider:       "kubernetes",
		},
		Kind:     "deployment",
		Labels:   deployment.GetLabels(),
		Manifest: deployment.Object,
		Moniker: Moniker{
			App:     application,
			Cluster: fmt.Sprintf("%s %s", "deployment", deployment.GetName()),
		},
		Name:         fmt.Sprintf("%s %s", "deployment", deployment.GetName()),
		DisplayName:  deployment.GetName(),
		Namespace:    deployment.GetNamespace(),
		ProviderType: "kubernetes",
		Region:       deployment.GetNamespace(),
		Type:         "kubernetes",
		UID:          string(deployment.GetUID()),
		Zone:         application,
	}
}

func buildServerGroups(replicaSets *unstructured.UnstructuredList,
	deployment unstructured.Unstructured,
	account, application string) []ServerGroupManagerServerGroup {
	sgs := []ServerGroupManagerServerGroup{}

	// Deployments manage replicasets, so build a list of managed replicasets for each deployment.
	for _, replicaSet := range replicaSets.Items {
		annotations := replicaSet.GetAnnotations()
		if annotations != nil {
			name := annotations["artifact.spinnaker.io/name"]
			t := annotations["artifact.spinnaker.io/type"]

			if strings.EqualFold(name, deployment.GetName()) &&
				strings.EqualFold(t, "kubernetes/deployment") {
				sequence := getSequence(annotations)
				s := ServerGroupManagerServerGroup{
					Account: account,
					Moniker: ServerGroupManagerServerGroupMoniker{
						App:      application,
						Cluster:  fmt.Sprintf("%s %s", "deployment", deployment.GetName()),
						Sequence: sequence,
					},
					Name:      fmt.Sprintf("%s %s", "replicaSet", replicaSet.GetName()),
					Namespace: replicaSet.GetNamespace(),
					Region:    replicaSet.GetNamespace(),
				}
				sgs = append(sgs, s)
			}
		}
	}

	return sgs
}

type LoadBalancers []LoadBalancer

type LoadBalancer struct {
	Account       string                    `json:"account"`
	CloudProvider string                    `json:"cloudProvider"`
	DispatchRules []interface{}             `json:"dispatchRules,omitempty"`
	HTTPURL       string                    `json:"httpUrl,omitempty"`
	HTTPSURL      string                    `json:"httpsUrl,omitempty"`
	Labels        map[string]string         `json:"labels,omitempty"`
	Moniker       Moniker                   `json:"moniker"`
	Name          string                    `json:"name"`
	DisplayName   string                    `json:"displayName"`
	Project       string                    `json:"project,omitempty"`
	Region        string                    `json:"region"`
	SelfLink      string                    `json:"selfLink,omitempty"`
	ServerGroups  []LoadBalancerServerGroup `json:"serverGroups"`
	Type          string                    `json:"type"`
	AccountName   string                    `json:"accountName,omitempty"`
	CreatedTime   int64                     `json:"createdTime,omitempty"`
	Key           Key                       `json:"key,omitempty"`
	Kind          string                    `json:"kind,omitempty"`
	Manifest      map[string]interface{}    `json:"manifest,omitempty"`
	ProviderType  string                    `json:"providerType,omitempty"`
	UID           string                    `json:"uid,omitempty"`
	Zone          string                    `json:"zone,omitempty"`
}

type LoadBalancerServerGroup struct {
	AllowsGradualTrafficMigration bool          `json:"allowsGradualTrafficMigration"`
	CloudProvider                 string        `json:"cloudProvider"`
	DetachedInstances             []interface{} `json:"detachedInstances"`
	Instances                     []interface{} `json:"instances"`
	IsDisabled                    bool          `json:"isDisabled"`
	Name                          string        `json:"name"`
	Region                        string        `json:"region"`
}

// List "load balancers", which for kubernetes are kinds "ingress" and "service".
func ListLoadBalancers(c *gin.Context) {
	sc := sql.Instance(c)
	application := c.Param("application")
	response := LoadBalancers{}
	wg := &sync.WaitGroup{}
	lbs := make(chan LoadBalancer, 100000)

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	wg.Add(len(accounts))

	// Don't actually return while attempting to create a list of load balancers.
	// We want to avoid the situation where a user cannot perform operations when any
	// cluster is not available.
	for _, account := range accounts {
		go listLoadBalancers(c, wg, lbs, account, application)
	}

	wg.Wait()

	close(lbs)

	for lb := range lbs {
		response = append(response, lb)
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Name < response[j].Name
	})

	c.JSON(http.StatusOK, response)
}

func listLoadBalancers(c *gin.Context, wg *sync.WaitGroup, lbs chan LoadBalancer,
	account, application string) {
	defer wg.Done()

	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	// Load balancer resources.
	resources := []string{
		"services",
		"ingresses",
	}

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	token, err := ac.Token(provider.TokenProvider)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	// Label selector for all that we are listing in the cluster. We
	// only want to list resources that have a label referencing the requested application.
	lo := metav1.ListOptions{
		LabelSelector:  kubernetes.LabelKubernetesName + "=" + application,
		TimeoutSeconds: &listTimeout,
	}

	for _, resource := range resources {
		results, err := client.ListResource(resource, lo)
		if err != nil {
			clouddriver.Log(err)
			continue
		}

		for _, result := range results.Items {
			lb := newLoadBalancer(result, account, application)
			lbs <- lb
		}
	}
}

func newLoadBalancer(u unstructured.Unstructured, account, application string) LoadBalancer {
	kind := strings.ToLower(u.GetKind())

	return LoadBalancer{
		Account:       account,
		AccountName:   account,
		CloudProvider: "kubernetes",
		Labels:        u.GetLabels(),
		Moniker: Moniker{
			App:     application,
			Cluster: fmt.Sprintf("%s %s", kind, u.GetName()),
		},
		Name:        fmt.Sprintf("%s %s", kind, u.GetName()),
		DisplayName: u.GetName(),
		Region:      u.GetNamespace(),
		Type:        "kubernetes",
		CreatedTime: u.GetCreationTimestamp().Unix() * 1000,
		Key: Key{
			Account:        account,
			Group:          u.GroupVersionKind().Group,
			KubernetesKind: kind,
			Name:           fmt.Sprintf("%s %s", kind, u.GetName()),
			Namespace:      u.GetNamespace(),
			Provider:       "kubernetes",
		},
		Kind:         kind,
		Manifest:     u.Object,
		ProviderType: "kubernetes",
		UID:          string(u.GetUID()),
		Zone:         application,
	}
}

type Clusters map[string][]string

// List clusters for a given application, which for kubernetes is a map of provider names to kubernetes deployment
// kinds and names.
//
// Clusters are kinds deployment, statefulSet, replicaSet, ingress, service, and daemonSet.
func ListClusters(c *gin.Context) {
	sc := sql.Instance(c)
	application := c.Param("application")

	rs, err := sc.ListKubernetesClustersByApplication(application)
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
	LoadBalancers       []interface{}                   `json:"loadBalancers"`
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

func ListServerGroups(c *gin.Context) {
	sc := sql.Instance(c)
	application := c.Param("application")
	response := ServerGroups{}
	wg := &sync.WaitGroup{}
	sgs := make(chan ServerGroup, 100000)

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	wg.Add(len(accounts))

	// List server groups concurrently across accounts.
	for _, account := range accounts {
		go listServerGroups(c, wg, sgs, account, application)
	}

	wg.Wait()

	close(sgs)

	for sg := range sgs {
		response = append(response, sg)
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Name < response[j].Name
	})

	c.JSON(http.StatusOK, response)
}

func listServerGroups(c *gin.Context, wg *sync.WaitGroup, sgs chan ServerGroup,
	account, application string) {
	defer wg.Done()

	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	// Resources which are "server groups".
	resources := []string{
		"replicaSets",
		"daemonSets",
		"statefulSets",
	}

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	token, err := ac.Token(provider.TokenProvider)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	lo := metav1.ListOptions{
		LabelSelector:  kubernetes.LabelKubernetesName + "=" + application,
		TimeoutSeconds: &listTimeout,
	}

	pods, err := client.ListResource("pods", lo)
	if err != nil {
		clouddriver.Log(err)
		return
	}

	for _, resource := range resources {
		results, err := client.ListResource(resource, lo)
		if err != nil {
			clouddriver.Log(err)
			continue
		}

		for _, result := range results.Items {
			sg := newServerGroup(result, pods, account)
			sgs <- sg
		}
	}
}

func newServerGroup(result unstructured.Unstructured,
	pods *unstructured.UnstructuredList, account string) ServerGroup {
	images := listImages(&result)
	desired := getDesiredReplicasCount(&result)

	serverGroupManagers := []ServerGroupServerGroupManager{}
	instances := buildInstances(pods, result)
	annotations := result.GetAnnotations()

	// Build server group manager
	managerName := annotations["artifact.spinnaker.io/name"]
	managerLocation := annotations["artifact.spinnaker.io/location"]
	managerType := annotations["artifact.spinnaker.io/type"]

	if managerType == "kubernetes/deployment" {
		sgm := ServerGroupServerGroupManager{
			Account:  account,
			Location: managerLocation,
			Name:     managerName,
		}
		serverGroupManagers = append(serverGroupManagers, sgm)
	}

	cluster := annotations["moniker.spinnaker.io/cluster"]
	app := annotations["moniker.spinnaker.io/application"]
	sequence := getSequence(annotations)

	return ServerGroup{
		Account: account,
		BuildInfo: BuildInfo{
			Images: images,
		},
		Capacity: Capacity{
			Desired: desired,
			Pinned:  false,
		},
		CloudProvider: "kubernetes",
		Cluster:       cluster,
		CreatedTime:   result.GetCreationTimestamp().Unix() * 1000,
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
		Name:                fmt.Sprintf("%s %s", result.GetKind(), result.GetName()),
		Namespace:           result.GetNamespace(),
		Region:              result.GetNamespace(),
		SecurityGroups:      nil,
		ServerGroupManagers: serverGroupManagers,
		Type:                "kubernetes",
		Labels:              result.GetLabels(),
	}
}

// getSequence returns the sequence of a given resource.
// A versioned resource contains its sequence in the
// `moniker.spinnaker.io/sequence` annotation.
// A resource which is owned by some deployment defines its sequence in
// the `deployment.kubernetes.io/revision` annotation.
func getSequence(annotations map[string]string) int {
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

		images = rs.ListImages()
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

// Get desired replicas for replicaSets, statefulSets, and daemonSets.
func getDesiredReplicasCount(result *unstructured.Unstructured) int32 {
	desired := int32(0)

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		if rs.GetReplicaSetSpec().Replicas != nil {
			desired = *rs.GetReplicaSetSpec().Replicas
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

// Get total replicas for replicaSets, statefulSets, and daemonSets.
func getTotalReplicasCount(result *unstructured.Unstructured) int32 {
	total := int32(0)

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		total = rs.GetReplicaSetStatus().Replicas
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

// Get ready replicas for replicaSets, statefulSets, and daemonSets.
func getReadyReplicasCount(result *unstructured.Unstructured) int32 {
	ready := int32(0)

	switch strings.ToLower(result.GetKind()) {
	case "replicaset":
		rs := kubernetes.NewReplicaSet(result.Object)
		if rs.GetReplicaSetSpec().Replicas != nil {
			ready = rs.GetReplicaSetStatus().ReadyReplicas
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

func buildInstances(pods *unstructured.UnstructuredList,
	serverGroup unstructured.Unstructured) []Instance {
	instances := []Instance{}

	for _, u := range pods.Items {
		p := kubernetes.NewPod(u.Object)
		for _, ownerReference := range p.GetObjectMeta().OwnerReferences {
			if strings.EqualFold(ownerReference.Name, serverGroup.GetName()) {
				state := "Up"
				if p.GetPodStatus().Phase != "Running" {
					state = "Down"
				}

				instance := Instance{
					AvailabilityZone: u.GetNamespace(),
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
					ID:          string(u.GetUID()),
					Name:        fmt.Sprintf("%s %s", "pod", u.GetName()),
				}
				instances = append(instances, instance)
			}
		}
	}

	return instances
}

// /applications/:application/serverGroups/:account/:location/:name
func GetServerGroup(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	account := c.Param("account")
	application := c.Param("application")
	location := c.Param("location")
	nameArray := strings.Split(c.Param("name"), " ")
	kind := nameArray[0]
	name := nameArray[1]

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	token, err := ac.Token(provider.TokenProvider)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	lo := metav1.ListOptions{
		LabelSelector:  kubernetes.LabelKubernetesName + "=" + application,
		TimeoutSeconds: &listTimeout,
	}

	result, err := client.Get(kind, name, location)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// "Instances" in kubernetes are pods.
	pods, err := client.ListResource("pods", lo)
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
		for _, ownerReference := range p.GetObjectMeta().OwnerReferences {
			if strings.EqualFold(ownerReference.Name, result.GetName()) {
				state := "Up"
				if p.GetPodStatus().Phase != "Running" {
					state = "Down"
				}

				annotations := p.GetObjectMeta().Annotations
				cluster := annotations["moniker.spinnaker.io/cluster"]
				app := annotations["moniker.spinnaker.io/application"]

				if app == "" {
					app = application
				}

				instance := Instance{
					Account:          account,
					AccountName:      account,
					AvailabilityZone: p.GetNamespace(),
					CloudProvider:    "kubernetes",
					CreatedTime:      p.GetObjectMeta().CreationTimestamp.Unix() * 1000,
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
					HumanReadableName: fmt.Sprintf("%s %s", "pod", p.GetName()),
					ID:                string(p.GetUID()),
					Key: Key{
						Account:        account,
						Group:          "pod",
						KubernetesKind: "pod",
						Name:           p.GetName(),
						Namespace:      p.GetNamespace(),
						Provider:       "kubernetes",
					},
					Kind:     "pod",
					Labels:   p.GetLabels(),
					Manifest: v.Object,
					Moniker: Moniker{
						App:     app,
						Cluster: cluster,
					},
					Name:         fmt.Sprintf("%s %s", "pod", p.GetName()),
					ProviderType: "kubernetes",
					Region:       p.GetNamespace(),
					Type:         "kubernetes",
					UID:          string(p.GetUID()),
					Zone:         p.GetNamespace(),
				}
				instances = append(instances, instance)
			}
		}
	}

	annotations := result.GetAnnotations()
	cluster := annotations["moniker.spinnaker.io/cluster"]
	app := annotations["moniker.spinnaker.io/application"]
	sequence := getSequence(annotations)

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
		CloudProvider:  "kubernetes",
		CreatedTime:    result.GetCreationTimestamp().Unix() * 1000,
		Disabled:       false,
		DisplayName:    result.GetName(),
		InstanceCounts: instanceCounts,
		Instances:      instances,
		Key: Key{
			Account:        account,
			Group:          result.GetKind(),
			KubernetesKind: result.GetKind(),
			Name:           result.GetName(),
			Namespace:      result.GetNamespace(),
			Provider:       "kubernetes",
		},
		Kind:          result.GetKind(),
		Labels:        result.GetLabels(),
		LoadBalancers: []interface{}{},
		Manifest:      result.Object,
		Moniker: ServerGroupMoniker{
			App:      app,
			Cluster:  cluster,
			Sequence: sequence,
		},
		Name:                fmt.Sprintf("%s %s", result.GetKind(), result.GetName()),
		Namespace:           result.GetNamespace(),
		ProviderType:        "kubernetes",
		Region:              result.GetNamespace(),
		SecurityGroups:      []interface{}{},
		ServerGroupManagers: []ServerGroupServerGroupManager{},
		Type:                "kubernetes",
		UID:                 string(result.GetUID()),
		Zone:                result.GetNamespace(),
		Zones:               []interface{}{},
		InsightActions:      []interface{}{},
	}

	c.JSON(http.StatusOK, response)
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

func GetJob(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	account := c.Param("account")
	// application := c.Param("application")
	location := c.Param("location")
	nameArray := strings.Split(c.Param("name"), " ")
	kind := nameArray[0]
	name := nameArray[1]

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	token, err := ac.Token(provider.TokenProvider)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
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
		Provider:    "kubernetes",
	}

	c.JSON(http.StatusOK, job)
}
