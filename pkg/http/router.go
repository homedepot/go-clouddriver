package http

import (
	"github.com/billiford/go-clouddriver/pkg/http/core"
	v1 "github.com/billiford/go-clouddriver/pkg/http/v1"
	"github.com/gin-gonic/gin"
)

// Define the API.
func Initialize(r *gin.Engine) {
	// API endpoints without a version will go under "core".
	api := r.Group("")
	{
		api.GET("/health", core.OK)

		// Force cache refresh.
		api.POST("/cache/kubernetes/manifest", core.OK)

		// Credentials API controller.
		api.GET("/credentials", core.ListCredentials)
		api.GET("/credentials/:account", core.GetAccountCredentials)

		// Applications API controller.
		api.GET("/applications", core.ListApplications)
		api.GET("/applications/:application/serverGroupManagers", core.ListServerGroupManagers)
		api.GET("/applications/:application/serverGroups", core.ListServerGroups)
		api.GET("/applications/:application/serverGroups/:account/:location/:name", core.GetServerGroup)
		api.GET("/applications/:application/loadBalancers", core.ListLoadBalancers)
		api.GET("/applications/:application/clusters", core.ListClusters)

		// Create a kubernetes operation - deploy/delete/scale manifest.
		r.POST("/kubernetes/ops", core.CreateKubernetesOperation)

		// Monitor deploy.
		r.GET("/manifests/:account/:location/:name", core.GetManifest)

		// Get results for a task triggered in CreateKubernetesOperation.
		r.GET("/task/:id", core.GetTask)

		// Not implemented.
		r.GET("/securityGroups", core.ListSecurityGroups)
		r.GET("/search", core.Search)
	}

	// New endpoint.
	api = r.Group("/v1")
	{
		// Providers endpoint for kubernetes.
		api.POST("/kubernetes/providers", v1.CreateKubernetesProvider)
	}
}
