package core

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	// projectResources consist of Kubernetes kinds DaemonSets, ReplicaSets,
	// and StatefulSets.
	projectResources = []string{
		"replicaSets",
		"daemonSets",
		"statefulSets",
	}
)

type Project struct {
	Account        string               `json:"account"`
	Applications   []ProjectApplication `json:"applications"`
	Detail         string               `json:"detail"`
	InstanceCounts InstanceCounts       `json:"instanceCounts"`
	Stack          string               `json:"stack"`
}

type ProjectApplication struct {
	Application string           `json:"application"`
	Clusters    []ProjectCluster `json:"clusters"`
	LastPush    int64            `json:"lastPush"`
}

type ProjectCluster struct {
	Builds         []ProjectClusterBuild `json:"builds"`
	InstanceCounts InstanceCounts        `json:"instanceCounts"`
	LastPush       int64                 `json:"lastPush"`
	Region         string                `json:"region"`
}

type ProjectClusterBuild struct {
	BuildNumber string   `json:"buildNumber"`
	Deployed    int64    `json:"deployed"`
	Images      []string `json:"images"`
}

// ListProjectClusters retrieves the cluster details for a Spinnaker project.
//
// See https://github.com/spinnaker/clouddriver/blob/master/clouddriver-core/src/main/java/com/netflix/spinnaker/clouddriver/core/ProjectClustersService.java
func (cc *Controller) ListProjectClusters(c *gin.Context) {
	response := []Project{}
	name := c.Param("project")

	// Get the Spinnaker project configuration from front50.
	project, err := cc.Front50Client.Project(name)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}
	// If the front50 project configuration has no clusters, return an empty
	// list.
	if len(project.Config.Clusters) == 0 {
		c.JSON(http.StatusOK, response)
		return
	}

	// Get all providers.
	providers, err := cc.SQLClient.ListKubernetesProviders()
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	accounts := []string{}
	// Get list of unique accounts names.
	for _, cluster := range project.Config.Clusters {
		if !contains(accounts, cluster.Account) {
			// Only consider accounts that exist in the database.
			for _, p := range providers {
				if strings.EqualFold(p.Name, cluster.Account) {
					accounts = append(accounts, cluster.Account)
					break
				}
			}
		}
	}

	// Get the Kubenetes resources for all accounts/applications of the project.
	rs, err := cc.listApplicationResources(c, projectResources, accounts, project.Config.Applications)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Loop thru each cluster of the front50 project configuration.
	for _, cluster := range project.Config.Clusters {
		// Do not include invalid accounts.
		if !contains(accounts, cluster.Account) {
			continue
		}

		p := Project{
			Account:      cluster.Account,
			Applications: []ProjectApplication{},
			Detail:       cluster.Detail,
			Stack:        cluster.Stack,
		}

		// If the front50 project cluster's applications list is empty
		// then use 'all' the applications of the front50 project.
		applications := cluster.Applications
		if len(cluster.Applications) == 0 {
			applications = project.Config.Applications
		}

		// Create the ProjectApplication for each application.
		for _, application := range applications {
			// Create project application.
			pa := ProjectApplication{
				Application: application,
				Clusters:    listProjectClusters(rs, cluster.Account, application, cluster.Stack, cluster.Detail),
				LastPush:    0,
			}
			// Update summary information.
			for _, pc := range pa.Clusters {
				// Update project application summary info
				pa.LastPush = max(pa.LastPush, pc.LastPush)
				// Update project level summary info
				p.InstanceCounts.Down += pc.InstanceCounts.Down
				p.InstanceCounts.OutOfService += pc.InstanceCounts.OutOfService
				p.InstanceCounts.Starting += pc.InstanceCounts.Starting
				p.InstanceCounts.Total += pc.InstanceCounts.Total
				p.InstanceCounts.Unknown += pc.InstanceCounts.Unknown
				p.InstanceCounts.Up += pc.InstanceCounts.Up
			}
			// Add to list of project applications.
			p.Applications = append(p.Applications, pa)
		}

		response = append(response, p)
	}

	c.JSON(http.StatusOK, response)
}

// listProjectClusters returns a list of Project Clusters built from
// the Kubernetes resources that match the front50 project cluster
// configuration (account, application, stack, and detail).
func listProjectClusters(rs []resource, account, application, stack, detail string) []ProjectCluster {
	// Map of project clusters keyed by region (namespace).
	var pcMap = map[string]ProjectCluster{}

	for _, r := range rs {
		if strings.EqualFold(r.account, account) &&
			strings.EqualFold(r.application, application) &&
			kubernetes.AnnotationMatches(r.u, kubernetes.AnnotationSpinnakerMonikerStack, stack) &&
			kubernetes.AnnotationMatches(r.u, kubernetes.AnnotationSpinnakerMonikerDetail, detail) {
			region := r.u.GetNamespace()
			pcMap[region] = addResourceToProjectCluster(pcMap[region], r.u)
		}
	}

	// Return project clusters as a list.
	pcs := make([]ProjectCluster, 0, len(pcMap))

	for _, pc := range pcMap {
		pcs = append(pcs, pc)
	}
	// Sort project clusters by region.
	sort.Slice(pcs, func(i, j int) bool {
		return pcs[i].Region < pcs[j].Region
	})

	return pcs
}

// addResourceToProjectCluster adds this Kubernetes resource's information
// (last push, images, instance counts) to this project cluster.
func addResourceToProjectCluster(pc ProjectCluster, u unstructured.Unstructured) ProjectCluster {
	lastPush := max(pc.LastPush, u.GetCreationTimestamp().Unix()*1000)
	images := []string{}

	if pc.Builds != nil {
		images = pc.Builds[0].Images
	}
	// Add this Kubernetes resource's images to unique list for the project cluster.
	for _, image := range listImages(&u) {
		if !contains(images, image) {
			images = append(images, image)
		}
	}

	return ProjectCluster{
		Builds: []ProjectClusterBuild{
			// go-clouddriver doesn't have support for Jenkins builds,
			// so create default build "0" to accumulate list of all images.
			// See https://github.com/spinnaker/clouddriver/blob/96755fec0c04b6e361efb6d1c19a7afc3926e302/clouddriver-core/src/main/java/com/netflix/spinnaker/clouddriver/core/ProjectClustersService.java#L287
			{
				BuildNumber: "0",
				Deployed:    lastPush,
				Images:      images,
			},
		},
		InstanceCounts: InstanceCounts{
			Down:         0,
			OutOfService: 0,
			Starting:     0,
			Total:        pc.InstanceCounts.Total + getTotalReplicasCount(&u),
			Unknown:      0,
			Up:           pc.InstanceCounts.Up + getReadyReplicasCount(&u),
		},
		LastPush: lastPush,
		Region:   u.GetNamespace(),
	}
}

// Max returns the larger of x or y.
func max(x, y int64) int64 {
	if x > y {
		return x
	}

	return y
}
