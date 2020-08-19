package v0

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/pod"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/replicaset"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ApplicationsResponse []Application

type Application struct {
	Attributes   ApplicationAttributes `json:"attributes"`
	ClusterNames map[string][]string   `json:"clusterNames"`
	Name         string                `json:"name"`
}

type ApplicationAttributes struct {
	Name string `json:"name"`
}

func ListApplications(c *gin.Context) {
	sc := sql.Instance(c)

	rs, err := sc.ListKubernetesResourcesByFields("account_name", "kind", "name", "spinnaker_app")
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	response := ApplicationsResponse{}
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

	c.JSON(http.StatusOK, response)
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

type ServerGroupManagersResponse []ServerGroupManager

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
	kc := kubernetes.Instance(c)
	spinnakerApp := c.Param("application")
	response := ServerGroupManagersResponse{}

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(spinnakerApp)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Don't actually return while attempting to create a list of server group managers.
	// We want to avoid the situation where a user cannot perform operations when any
	// cluster is not available.
	for _, account := range accounts {
		provider, err := sc.GetKubernetesProvider(account)
		if err != nil {
			log.Println("unable to get kubernetes provider for account", account)
			continue
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			log.Println("error decoding ca data for account", account)
			continue
		}

		config := &rest.Config{
			Host:        provider.Host,
			BearerToken: os.Getenv("BEARER_TOKEN"),
			TLSClientConfig: rest.TLSClientConfig{
				CAData: cd,
			},
		}

		if err = kc.WithConfig(config); err != nil {
			log.Println("error creating dynamic client for account", account)
			continue
		}

		deployments := &unstructured.UnstructuredList{}
		replicaSets := &unstructured.UnstructuredList{}

		lo := metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=" + spinnakerApp,
		}

		deploymentGVK := schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    "deployment",
		}
		deploymentGVR := schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "deployments",
		}
		replicaSetGVR := schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "replicasets",
		}

		deployments, err = kc.List(deploymentGVR, lo)
		if err != nil {
			log.Println("error listing deployments:", err.Error())
			continue
		}

		replicaSets, err = kc.List(replicaSetGVR, lo)
		if err != nil {
			log.Println("error listing replicaSets:", err.Error())
			continue
		}

		for _, deployment := range deployments.Items {
			sgs := []ServerGroupManagerServerGroup{}
			// Deployments manage replicasets, so build a list of managed replicasets for each deployment.
			for _, replicaSet := range replicaSets.Items {
				annotations := replicaSet.GetAnnotations()
				if annotations != nil {
					var name, t string
					if _, ok := annotations["artifact.spinnaker.io/name"]; ok {
						name = annotations["artifact.spinnaker.io/name"]
					}
					if _, ok := annotations["artifact.spinnaker.io/type"]; ok {
						t = annotations["artifact.spinnaker.io/type"]
					}
					if name != "" && t != "" {
						if strings.EqualFold(name, deployment.GetName()) &&
							strings.EqualFold(t, "kubernetes/deployment") {
							sequence := 0
							deploymentAnnotations := deployment.GetAnnotations()
							if deploymentAnnotations != nil {
								if _, ok := deploymentAnnotations["deployment.kubernetes.io/revision"]; ok {
									sequence, _ = strconv.Atoi(deploymentAnnotations["deployment.kubernetes.io/revision"])
								}
							}
							s := ServerGroupManagerServerGroup{
								Account: account,
								Moniker: ServerGroupManagerServerGroupMoniker{
									App:      spinnakerApp,
									Cluster:  fmt.Sprintf("%s %s", deploymentGVK.Kind, deployment.GetName()),
									Sequence: sequence,
								},
								Name:      fmt.Sprintf("%s %s", "replicaset", replicaSet.GetName()),
								Namespace: replicaSet.GetNamespace(),
								Region:    replicaSet.GetNamespace(),
							}
							sgs = append(sgs, s)
						}
					}
				}
			}

			sgr := ServerGroupManager{
				Account:       account,
				AccountName:   account,
				CloudProvider: "kubernetes",
				CreatedTime:   deployment.GetCreationTimestamp().Unix() * 1000,
				Key: Key{
					Account:        account,
					Group:          deploymentGVK.Group,
					KubernetesKind: deploymentGVK.Kind,
					Name:           deployment.GetName(),
					Namespace:      deployment.GetNamespace(),
					Provider:       "kubernetes",
				},
				Kind:     deploymentGVK.Kind,
				Labels:   deployment.GetLabels(),
				Manifest: deployment.Object,
				Moniker: Moniker{
					App:     spinnakerApp,
					Cluster: fmt.Sprintf("%s %s", deploymentGVK.Kind, deployment.GetName()),
				},
				Name:         fmt.Sprintf("%s %s", deploymentGVK.Kind, deployment.GetName()),
				ProviderType: "kubernetes",
				Region:       deployment.GetNamespace(),
				ServerGroups: sgs,
				Type:         "kubernetes",
				UID:          string(deployment.GetUID()),
				Zone:         spinnakerApp,
			}
			response = append(response, sgr)
		}
	}

	c.JSON(http.StatusOK, response)
}

type LoadBalancersResponse []LoadBalancer

type LoadBalancer struct {
	Account       string                    `json:"account"`
	CloudProvider string                    `json:"cloudProvider"`
	DispatchRules []interface{}             `json:"dispatchRules,omitempty"`
	HTTPURL       string                    `json:"httpUrl,omitempty"`
	HTTPSURL      string                    `json:"httpsUrl,omitempty"`
	Labels        map[string]string         `json:"labels,omitempty"`
	Moniker       Moniker                   `json:"moniker"`
	Name          string                    `json:"name"`
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
	kc := kubernetes.Instance(c)
	spinnakerApp := c.Param("application")
	response := LoadBalancersResponse{}

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(spinnakerApp)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Don't actually return while attempting to create a list of server group managers.
	// We want to avoid the situation where a user cannot perform operations when any
	// cluster is not available.
	for _, account := range accounts {
		provider, err := sc.GetKubernetesProvider(account)
		if err != nil {
			log.Println("unable to get kubernetes provider for account", account)
			continue
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			log.Println("error decoding ca data for account", account)
			continue
		}

		config := &rest.Config{
			Host:        provider.Host,
			BearerToken: os.Getenv("BEARER_TOKEN"),
			TLSClientConfig: rest.TLSClientConfig{
				CAData: cd,
			},
		}

		if err = kc.WithConfig(config); err != nil {
			log.Println("error creating dynamic client for account", account)
			continue
		}

		// Label selector for all that we are listing in the cluster. We
		// only want to list resources that have a label referencing the requested application.
		lo := metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=" + spinnakerApp,
		}

		// Create a GVR for ingresses.
		ingressGVR := schema.GroupVersionResource{
			Group:    "networking.k8s.io",
			Version:  "v1beta1",
			Resource: "ingresses",
		}
		ingressGVK := schema.GroupVersionKind{
			Group:   "networking.k8s.io",
			Version: "v1beta1",
			Kind:    "ingress",
		}

		ingresses, err := kc.List(ingressGVR, lo)
		if err != nil {
			log.Println("error listing ingresses:", err.Error())
			continue
		}

		for _, ingress := range ingresses.Items {
			lb := LoadBalancer{
				Account:       account,
				AccountName:   account,
				CloudProvider: "kubernetes",
				Labels:        ingress.GetLabels(),
				Moniker: Moniker{
					App:     spinnakerApp,
					Cluster: fmt.Sprintf("%s %s", ingressGVK.Kind, ingress.GetName()),
				},
				Name:        fmt.Sprintf("%s %s", "ingress", ingress.GetName()),
				Region:      ingress.GetNamespace(),
				Type:        "kubernetes",
				CreatedTime: ingress.GetCreationTimestamp().Unix() * 1000,
				Key: Key{
					Account:        account,
					Group:          ingressGVK.Group,
					KubernetesKind: ingressGVK.Kind,
					Name:           ingress.GetName(),
					Namespace:      ingress.GetNamespace(),
					Provider:       "kubernetes",
				},
				Kind:         ingressGVK.Kind,
				Manifest:     ingress.Object,
				ProviderType: "kubernetes",
				UID:          string(ingress.GetUID()),
				Zone:         spinnakerApp,
			}
			response = append(response, lb)
		}

		// Create a GVR for services.
		serviceGVR := schema.GroupVersionResource{
			Version:  "v1",
			Resource: "services",
		}
		serviceGVK := schema.GroupVersionKind{
			Version: "v1",
			Kind:    "service",
		}

		services, err := kc.List(serviceGVR, lo)
		if err != nil {
			log.Println("error listing services:", err.Error())
			continue
		}

		for _, service := range services.Items {
			lb := LoadBalancer{
				Account:       account,
				AccountName:   account,
				CloudProvider: "kubernetes",
				Labels:        service.GetLabels(),
				Moniker: Moniker{
					App:     spinnakerApp,
					Cluster: fmt.Sprintf("%s %s", serviceGVK.Kind, service.GetName()),
				},
				Name:        fmt.Sprintf("%s %s", "service", service.GetName()),
				Region:      service.GetNamespace(),
				Type:        "kubernetes",
				CreatedTime: service.GetCreationTimestamp().Unix() * 1000,
				Key: Key{
					Account:        account,
					Group:          serviceGVK.Group,
					KubernetesKind: serviceGVK.Kind,
					Name:           service.GetName(),
					Namespace:      service.GetNamespace(),
					Provider:       "kubernetes",
				},
				Kind:         serviceGVK.Kind,
				Manifest:     service.Object,
				ProviderType: "kubernetes",
				UID:          string(service.GetUID()),
				Zone:         spinnakerApp,
			}
			response = append(response, lb)
		}
	}

	c.JSON(http.StatusOK, response)
}

type ClustersResponse struct {
	SpinClusterAccount []string `json:"spin-cluster-account"`
}

func ListClusters(c *gin.Context) {
	cr := ClustersResponse{
		SpinClusterAccount: []string{
			"deployment cleanup-operator",
			"deployment nginx-deployment",
			"replicaSet cleanup-operator-6f5df67cf9",
			"replicaSet demo-deployment-5fc8ffdb68",
			"replicaSet frontend",
			"replicaSet hello-app-red-black",
			"replicaSet hello-app-red-black-v006",
			"replicaSet hello-app-red-black-v007",
			"service hello-app-red-black",
		},
	}

	c.JSON(http.StatusOK, cr)
}

type ServerGroupsResponse []ServerGroup

type ServerGroup struct {
	Account             string                          `json:"account"`
	BuildInfo           BuildInfo                       `json:"buildInfo"`
	Capacity            Capacity                        `json:"capacity"`
	CloudProvider       string                          `json:"cloudProvider"`
	Cluster             string                          `json:"cluster"`
	CreatedTime         int64                           `json:"createdTime"`
	InstanceCounts      InstanceCounts                  `json:"instanceCounts"`
	Instances           []Instance                      `json:"instances"`
	IsDisabled          bool                            `json:"isDisabled"`
	LoadBalancers       []interface{}                   `json:"loadBalancers"`
	Moniker             ServerGroupMoniker              `json:"moniker"`
	Name                string                          `json:"name"`
	Region              string                          `json:"region"`
	SecurityGroups      []interface{}                   `json:"securityGroups"`
	ServerGroupManagers []ServerGroupServerGroupManager `json:"serverGroupManagers"`
	Type                string                          `json:"type"`
	Labels              map[string]string               `json:"labels,omitempty"`
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
	Desired int  `json:"desired"`
	Pinned  bool `json:"pinned"`
}

type InstanceCounts struct {
	Down         int `json:"down"`
	OutOfService int `json:"outOfService"`
	Starting     int `json:"starting"`
	Total        int `json:"total"`
	Unknown      int `json:"unknown"`
	Up           int `json:"up"`
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
	kc := kubernetes.Instance(c)
	spinnakerApp := c.Param("application")
	response := ServerGroupsResponse{}

	accounts, err := sc.ListKubernetesAccountsBySpinnakerApp(spinnakerApp)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// Don't actually return while attempting to create a list of server groups.
	// We want to avoid the situation where a user cannot perform operations when any
	// cluster is not available.
	for _, account := range accounts {
		provider, err := sc.GetKubernetesProvider(account)
		if err != nil {
			log.Println("unable to get kubernetes provider for account", account)
			continue
		}

		cd, err := base64.StdEncoding.DecodeString(provider.CAData)
		if err != nil {
			log.Println("error decoding ca data for account", account)
			continue
		}

		config := &rest.Config{
			Host:        provider.Host,
			BearerToken: os.Getenv("BEARER_TOKEN"),
			TLSClientConfig: rest.TLSClientConfig{
				CAData: cd,
			},
		}

		if err = kc.WithConfig(config); err != nil {
			log.Println("error creating dynamic client for account", account)
			continue
		}

		lo := metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=" + spinnakerApp,
		}

		// Create a GVR for replicasets.
		replicaSetGVR := schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "replicasets",
		}
		podsGVR := schema.GroupVersionResource{
			Version:  "v1",
			Resource: "pods",
		}
		replicaSetGVK := schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    "replicaset",
		}

		replicaSets, err := kc.List(replicaSetGVR, lo)
		if err != nil {
			log.Println("error listing replicaSets:", err.Error())
			continue
		}

		pods, err := kc.List(podsGVR, lo)
		if err != nil {
			log.Println("error listing pods:", err.Error())
			continue
		}

		for _, replicaSet := range replicaSets.Items {
			b, _ := json.Marshal(replicaSet.Object)
			rs := v1.ReplicaSet{}
			json.Unmarshal(b, &rs)
			images := []string{}

			for _, container := range rs.Spec.Template.Spec.Containers {
				images = append(images, container.Image)
			}
			desired := 0
			if rs.Spec.Replicas != nil {
				desired = int(*rs.Spec.Replicas)
			}

			instances := []Instance{}
			for _, pod := range pods.Items {
				b, _ = json.Marshal(pod.Object)
				p := &corev1.Pod{}
				json.Unmarshal(b, &p)
				for _, ownerReference := range p.ObjectMeta.OwnerReferences {
					if strings.EqualFold(ownerReference.Name, replicaSet.GetName()) {
						state := "Up"
						if p.Status.Phase != "Running" {
							state = "Down"
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
						instances = append(instances, instance)
					}
				}
			}

			serverGroupManagers := []ServerGroupServerGroupManager{}

			// Build server group manager
			{
				annotations := replicaSet.GetAnnotations()
				if annotations != nil {
					var managerName, managerLocation, managerType string
					if _, ok := annotations["artifact.spinnaker.io/name"]; ok {
						managerName = annotations["artifact.spinnaker.io/name"]
					}
					if _, ok := annotations["artifact.spinnaker.io/location"]; ok {
						managerLocation = annotations["artifact.spinnaker.io/location"]
					}
					if _, ok := annotations["artifact.spinnaker.io/type"]; ok {
						managerType = annotations["artifact.spinnaker.io/type"]
					}
					if managerType == "kubernetes/deployment" {
						sgm := ServerGroupServerGroupManager{
							Account:  account,
							Location: managerLocation,
							Name:     managerName,
						}
						serverGroupManagers = append(serverGroupManagers, sgm)
					}
				}
			}

			sgs := ServerGroup{
				Account: account,
				BuildInfo: BuildInfo{
					Images: images,
				},
				Capacity: Capacity{
					Desired: desired,
					Pinned:  false,
				},
				CloudProvider: "kubernetes",
				Cluster:       fmt.Sprintf("%s %s", replicaSetGVK.Kind, replicaSet.GetName()),
				CreatedTime:   replicaSet.GetCreationTimestamp().Unix() * 1000,
				InstanceCounts: InstanceCounts{
					Down:         0,
					OutOfService: 0,
					Starting:     0,
					Total:        int(rs.Status.Replicas),
					Unknown:      0,
					Up:           int(rs.Status.ReadyReplicas),
				},
				Instances:     instances,
				IsDisabled:    false,
				LoadBalancers: nil,
				Moniker: ServerGroupMoniker{
					App:      spinnakerApp,
					Cluster:  fmt.Sprintf("%s %s", replicaSetGVK.Kind, replicaSet.GetName()),
					Sequence: 0,
				},
				Name:                fmt.Sprintf("%s %s", "replicaset", replicaSet.GetName()),
				Region:              replicaSet.GetNamespace(),
				SecurityGroups:      nil,
				ServerGroupManagers: serverGroupManagers,
				Type:                "kubernetes",
				Labels:              replicaSet.GetLabels(),
			}
			response = append(response, sgs)
		}
	}

	c.JSON(http.StatusOK, response)
}

type GetServerGroupResponse struct {
	Account        string            `json:"account"`
	AccountName    string            `json:"accountName"`
	BuildInfo      BuildInfo         `json:"buildInfo"`
	Capacity       Capacity          `json:"capacity"`
	CloudProvider  string            `json:"cloudProvider"`
	CreatedTime    int64             `json:"createdTime"`
	Disabled       bool              `json:"disabled"`
	InstanceCounts InstanceCounts    `json:"instanceCounts"`
	Instances      []Instance        `json:"instances"`
	Key            Key               `json:"key"`
	Kind           string            `json:"kind"`
	Labels         map[string]string `json:"labels"`
	// LaunchConfig struct {} `json:"launchConfig"`
	LoadBalancers       []interface{}                   `json:"loadBalancers"`
	Manifest            map[string]interface{}          `json:"manifest"`
	Moniker             ServerGroupMoniker              `json:"moniker"`
	Name                string                          `json:"name"`
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

// /applications/:application/serverGroups/:account/:location/:name
func GetServerGroup(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)
	account := c.Param("account")
	application := c.Param("application")
	location := c.Param("location")
	n := c.Param("name")
	a := strings.Split(n, " ")
	kind := a[0]
	name := a[1]

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: os.Getenv("BEARER_TOKEN"),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	if err = kc.WithConfig(config); err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	lo := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=" + application,
	}

	podsGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}

	result, err := kc.Get(kind, name, location)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// "Instances" in kubernetes are pods.
	pods, err := kc.List(podsGVR, lo)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	images := []string{}
	desired := 0
	instanceCounts := InstanceCounts{}
	if strings.EqualFold(kind, "replicaset") {
		rs := replicaset.New(result.Object)
		for _, container := range rs.Spec.Template.Spec.Containers {
			images = append(images, container.Image)
		}
		if rs.Spec.Replicas != nil {
			desired = int(*rs.Spec.Replicas)
		}
		instanceCounts.Total = int(rs.Status.Replicas)
		instanceCounts.Up = int(rs.Status.ReadyReplicas)
	}

	instances := []Instance{}
	for _, v := range pods.Items {
		p := pod.New(v.Object)
		for _, ownerReference := range p.ObjectMeta.OwnerReferences {
			if strings.EqualFold(ownerReference.Name, result.GetName()) {
				state := "Up"
				if p.Status.Phase != "Running" {
					state = "Down"
				}
				labels := p.ObjectMeta.Labels
				cluster := ""
				app := application
				if _, ok := labels["moniker.spinnaker.io/cluster"]; ok {
					cluster = labels["moniker.spinnaker.io/cluster"]
				}
				if _, ok := labels["moniker.spinnaker.io/application"]; ok {
					app = labels["moniker.spinnaker.io/application"]
				}
				instance := Instance{
					Account:          account,
					AccountName:      account,
					AvailabilityZone: p.GetNamespace(),
					CloudProvider:    "kubernetes",
					CreatedTime:      p.ObjectMeta.CreationTimestamp.Unix() * 1000,
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

	labels := result.GetLabels()
	cluster := ""
	app := application
	sequence := 0
	if _, ok := labels["moniker.spinnaker.io/cluster"]; ok {
		cluster = labels["moniker.spinnaker.io/cluster"]
	}
	if _, ok := labels["moniker.spinnaker.io/application"]; ok {
		app = labels["moniker.spinnaker.io/application"]
	}
	if _, ok := labels["deployment.kubernetes.io/revision"]; ok {
		sequence, _ = strconv.Atoi(labels["deployment.kubernetes.io/revision"])
	}

	response := GetServerGroupResponse{
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
