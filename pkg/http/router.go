package http

import (
	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/http/core"
	v1 "github.com/homedepot/go-clouddriver/pkg/http/v1"
	"github.com/homedepot/go-clouddriver/pkg/middleware"
)

// Define the API.
func Initialize(r *gin.Engine) {
	// API endpoints without a version will go under "core".
	{
		api := r.Group("")
		api.GET("/health", core.OK)

		// Force cache refresh.
		api.POST("/cache/kubernetes/manifest", core.OK)

		// Credentials API controller.
		api.GET("/credentials", core.ListCredentials)
		api.GET("/credentials/:account", core.GetAccountCredentials)

		// Applications API controller.
		//
		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ApplicationsController.groovy#L38
		// @PreAuthorize("#restricted ? @fiatPermissionEvaluator.storeWholePermission() : true")
		// @PostFilter("#restricted ? hasPermission(filterObject.name, 'APPLICATION', 'READ') : true")
		api.GET("/applications", middleware.PostFilterAuthorizedApplications("READ"), core.ListApplications)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupManagerController.java#L39
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostFilter("hasPermission(filterObject.account, 'ACCOUNT', 'READ')")
		api.GET("/applications/:application/serverGroupManagers", middleware.AuthApplication("READ"), core.ListServerGroupManagers)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupController.groovy#L172
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/serverGroups", middleware.AuthApplication("READ"), core.ListServerGroups)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupController.groovy#L75
		// @PreAuthorize("hasPermission(#account, 'ACCOUNT', 'READ')")
		// @PostAuthorize("hasPermission(returnObject?.moniker?.app, 'APPLICATION', 'READ')")
		// textPayload: "Headers: map[Accept:[application/json] Accept-Encoding:[gzip] Connection:[Keep-Alive] User-Agent:[okhttp/3.14.9] X-Spinnaker-Accounts:[gke_github-replication-sandbox_us-east1_sandbox-us-east1-agent-dev,gke_github-replication-sandbox_us-east1_sandbox-us-east1-dev,gke_github-replication-sandbox_us-central1-c_prom-test] X-Spinnaker-Application:[smoketests] X-Spinnaker-User:[me@me.com]]"
		api.GET("/applications/:application/serverGroups/:account/:location/:name", middleware.AuthAccount("READ"), core.GetServerGroup)

		// https: //github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/LoadBalancerController.groovy#L42
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/loadBalancers", middleware.AuthApplication("READ"), core.ListLoadBalancers)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ClusterController.groovy#L44
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission() and hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/clusters", middleware.AuthApplication("READ"), core.ListClusters)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/JobController.groovy#L35
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ') and hasPermission(#account, 'ACCOUNT', 'READ')")
		// @ApiOperation(value = "Collect a JobStatus", notes = "Collects the output of the job.")
		api.GET("/applications/:application/jobs/:account/:location/:name", middleware.AuthApplication("READ"), middleware.AuthAccount("READ"), core.GetJob)
		// Delete job always fails, so we do not need to pass through the auth middlewares.
		api.DELETE("/applications/:application/jobs/:account/:location/:name", core.DeleteJob)

		// Create a kubernetes operation - deploy/delete/scale manifest.
		api.POST("/kubernetes/ops", middleware.TaskID(), core.CreateKubernetesOperation)

		// Manifests API controller.
		api.GET("/manifests/:account/:location/:kind", core.GetManifest)
		api.GET("/manifests/:account/:location/:kind/cluster/:application/:cluster/dynamic/:target", core.GetManifestByTarget)

		// Instances API controller.
		api.GET("/instances/:account/:location/:name", core.GetInstance)
		api.GET("/instances/:account/:location/:name/console", core.GetInstanceConsole)

		// Get results for a task triggered in CreateKubernetesOperation.
		api.GET("/task/:id", core.GetTask)

		// Generic search endpoint.
		//
		// https://github.com/spinnaker/clouddriver/blob/0524d08f6bcf775c469a0576a79b2679b5653325/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/SearchController.groovy#L55
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission()")
		api.GET("/search", core.Search)

		// Not implemented.
		//
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission()")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/securityGroups", core.ListSecurityGroups)

		// Artifacts API controller.
		api.GET("/artifacts/credentials", core.ListArtifactCredentials)
		api.GET("/artifacts/account/:accountName/names", core.ListHelmArtifactAccountNames)
		api.GET("/artifacts/account/:accountName/versions", core.ListHelmArtifactAccountVersions)
		api.PUT("/artifacts/fetch/", core.GetArtifact)

		// Features.
		api.GET("/features/stages", core.ListStages)
	}

	// New endpoint.
	{
		api := r.Group("/v1")
		// Providers endpoint for kubernetes.
		api.GET("/kubernetes/providers/:name", v1.GetKubernetesProvider)
		api.POST("/kubernetes/providers", v1.CreateKubernetesProvider)
		api.PUT("/kubernetes/providers", v1.CreateOrReplaceKubernetesProvider)
		api.DELETE("/kubernetes/providers/:name", v1.DeleteKubernetesProvider)
	}
}
