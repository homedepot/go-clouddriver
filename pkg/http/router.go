package http

import (
	v0 "github.com/billiford/go-clouddriver/pkg/http/v0"
	v1 "github.com/billiford/go-clouddriver/pkg/http/v1"
	"github.com/gin-gonic/gin"
)

// Define the API.
func Initialize(r *gin.Engine) {
	// API endpoints without a version will go under "v0".
	api := r.Group("")
	{
		api.GET("/health", v0.OK)
		// Force cache refresh.
		api.POST("/cache/kubernetes/manifest", v0.OK)
		api.GET("/credentials", v0.ListCredentials)
		api.GET("/credentials/:account", v0.GetAccountCredentials)

		// Applications API controller.
		api.GET("/applications", v0.ListApplications)
		api.GET("/applications/:application/serverGroupManagers", v0.ListServerGroupManagers)
		api.GET("/applications/:application/serverGroups", v0.ListServerGroups)
		api.GET("/applications/:application/serverGroups/:account/:location/:name", v0.GetServerGroup)
		api.GET("/applications/:application/loadBalancers", v0.ListLoadBalancers)
		api.GET("/applications/:application/clusters", v0.ListClusters)

		// Create a kubernetes operation - deploy/delete/scale manifest.
		r.POST("/kubernetes/ops", v0.CreateKubernetesOperation)

		// Monitor deploy.
		r.GET("/manifests/:account/:location/:name", v0.GetManifest)

		r.GET("/task/:id", v0.GetTask)

		r.GET("/securityGroups", v0.ListSecurityGroups)
		r.GET("/search", v0.Search)
	}

	api = r.Group("/v1")
	{
		api.POST("/kubernetes/providers", v1.CreateKubernetesProvider)
	}
}
