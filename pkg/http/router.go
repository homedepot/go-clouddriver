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
		api.GET("/applications/:application/jobs/:account/:location/:name", core.GetJob)

		// Create a kubernetes operation - deploy/delete/scale manifest.
		api.POST("/kubernetes/ops", core.CreateKubernetesOperation)

		// Monitor deploy.
		api.GET("/manifests/:account/:location/:name", core.GetManifest)

		// Get results for a task triggered in CreateKubernetesOperation.
		api.GET("/task/:id", core.GetTask)

		// Not implemented.
		api.GET("/securityGroups", core.ListSecurityGroups)
		api.GET("/search", core.Search)

		// Artifacts controller.
		api.GET("/artifacts/credentials", core.ListArtifactCredentials)
		api.GET("/artifacts/account/:accountName/names", core.ListHelmArtifactAccountNames)
		api.GET("/artifacts/account/:accountName/versions", core.ListHelmArtifactAccountVersions)
		api.PUT("/artifacts/fetch/", core.GetArtifact)

		// Features.
		api.GET("/features/stages", core.ListStages)
	}

	// New endpoint.
	api = r.Group("/v1")
	{
		// Providers endpoint for kubernetes.
		api.POST("/kubernetes/providers", v1.CreateKubernetesProvider)
	}
}
