package api

import (
	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/api/core"
	v1 "github.com/homedepot/go-clouddriver/internal/api/v1"
	"github.com/homedepot/go-clouddriver/internal/middleware"
)

const (
	defaultCacheControlMaxAge = 30
	headerXSpinnakerAccounts  = "X-Spinnaker-Accounts"
)

// Server hold the gin engine and any clients we need for the API.
type Server struct {
	c *internal.Controller
	e *gin.Engine
	// v is verbose request logging.
	v bool
}

// NewServer returns a new instance of Server.
func NewServer(e *gin.Engine) *Server {
	return &Server{
		e: e,
	}
}

// WithController sets the internal controller instance for the server.
func (s *Server) WithController(c *internal.Controller) {
	s.c = c
}

// WithVerboseRequestLogging sets the server to use verbose request logging.
func (s *Server) WithVerboseRequestLogging() {
	s.v = true
}

// Setup sets any global middlewares then initializes the API.
func (s *Server) Setup() {
	s.e.Use(middleware.HandleError())

	// Verbose request logging.
	if s.v {
		s.e.Use(middleware.LogRequest())
	}

	// API endpoints without a version will go under "core".
	{
		// Declare all controllers.
		mc := &middleware.Controller{
			Controller: s.c,
		}
		c := &core.Controller{
			Controller: s.c,
		}
		api := s.e.Group("")
		api.GET("/health", core.OK)

		// Force cache refresh.
		api.POST("/cache/kubernetes/manifest", core.OK)

		// Credentials API controller.
		api.GET("/credentials", c.ListCredentials)
		api.GET("/credentials/:account", c.GetAccountCredentials)

		// Applications API controller.
		//
		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ApplicationsController.groovy#L38
		// @PreAuthorize("#restricted ? @fiatPermissionEvaluator.storeWholePermission() : true")
		// @PostFilter("#restricted ? hasPermission(filterObject.name, 'APPLICATION', 'READ') : true")
		api.GET("/applications", mc.PostFilterAuthorizedApplications("READ"), c.ListApplications)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupManagerController.java#L39
		api.GET("/applications/:application/serverGroupManagers", middleware.CacheControl(defaultCacheControlMaxAge), c.ListServerGroupManagers)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupController.groovy#L172
		api.GET("/applications/:application/serverGroups", middleware.CacheControl(defaultCacheControlMaxAge), c.ListServerGroups)

		// https: //github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/LoadBalancerController.groovy#L42
		api.GET("/applications/:application/loadBalancers", middleware.CacheControl(defaultCacheControlMaxAge), c.ListLoadBalancers)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupController.groovy#L75
		// @PreAuthorize("hasPermission(#account, 'ACCOUNT', 'READ')")
		// @PostAuthorize("hasPermission(returnObject?.moniker?.app, 'APPLICATION', 'READ')")
		// textPayload: "Headers: map[Accept:[application/json] Accept-Encoding:[gzip] Connection:[Keep-Alive] User-Agent:[okhttp/3.14.9] X-Spinnaker-Accounts:[gke_github-replication-sandbox_us-east1_sandbox-us-east1-agent-dev,gke_github-replication-sandbox_us-east1_sandbox-us-east1-dev,gke_github-replication-sandbox_us-central1-c_prom-test] X-Spinnaker-Application:[smoketests] X-Spinnaker-User:[me@me.com]]"
		api.GET("/applications/:application/serverGroups/:account/:location/:name", mc.AuthAccount("READ"), c.GetServerGroup)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ClusterController.groovy#L44
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission() and hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/clusters", mc.AuthApplication("READ"), c.ListClusters)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/JobController.groovy#L35
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ') and hasPermission(#account, 'ACCOUNT', 'READ')")
		// @ApiOperation(value = "Collect a JobStatus", notes = "Collects the output of the job.")
		api.GET("/applications/:application/jobs/:account/:location/:name", mc.AuthApplication("READ"), mc.AuthAccount("READ"), c.GetJob)
		// Delete job always fails, so we do not need to pass through the auth middlewares.
		api.DELETE("/applications/:application/jobs/:account/:location/:name", core.DeleteJob)

		// Create a kubernetes operation - deploy/delete/scale manifest.
		api.POST("/kubernetes/ops", mc.AuthOps(), middleware.TaskID(), c.CreateKubernetesOperation)

		// Manifests API controller.
		api.GET("/manifests/:account/:location/:kind", c.GetManifest)
		api.GET("/manifests/:account/:location/:kind/cluster/:application/:cluster", c.ListManifestsByCluster)
		api.GET("/manifests/:account/:location/:kind/cluster/:application/:cluster/dynamic/:criteria", c.GetManifestByCriteria)

		// Instances API controller.
		api.GET("/instances/:account/:location/:name", c.GetInstance)
		api.GET("/instances/:account/:location/:name/console", c.GetInstanceConsole)

		// Get results for a task triggered in CreateKubernetesOperation.
		api.GET("/task/:id", c.GetTask)

		// Generic search endpoint.
		// https://github.com/spinnaker/clouddriver/blob/0524d08f6bcf775c469a0576a79b2679b5653325/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/SearchController.groovy#L55
		api.GET("/search", middleware.CacheControl(defaultCacheControlMaxAge), middleware.Vary(headerXSpinnakerAccounts), c.Search)

		// Not implemented.
		//
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission()")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/securityGroups", core.ListSecurityGroups)

		// Artifacts API controller.
		api.GET("/artifacts/credentials", c.ListArtifactCredentials)
		api.GET("/artifacts/account/:accountName/names", c.ListHelmArtifactAccountNames)
		api.GET("/artifacts/account/:accountName/versions", c.ListHelmArtifactAccountVersions)
		api.PUT("/artifacts/fetch/", c.GetArtifact)

		// Features.
		api.GET("/features/stages", core.ListStages)

		// Projects API controller.
		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ProjectController.groovy
		api.GET("/projects/:project/clusters", middleware.CacheControl(defaultCacheControlMaxAge), c.ListProjectClusters)
	}

	// V1 endpoint.
	{
		c := &v1.Controller{
			Controller: s.c,
		}
		api := s.e.Group("/v1")
		// Providers endpoint for kubernetes.
		api.GET("/kubernetes/providers", c.ListKubernetesProvider)
		api.GET("/kubernetes/providers/:name", c.GetKubernetesProvider)
		api.POST("/kubernetes/providers", c.CreateKubernetesProvider)
		api.PUT("/kubernetes/providers", c.CreateOrReplaceKubernetesProvider)
		api.DELETE("/kubernetes/providers/:name", c.DeleteKubernetesProvider)
		// Resources endpoint for kubernetes.
		api.PUT("/kubernetes/providers/:name/resources", c.LoadKubernetesResources)
		api.DELETE("/kubernetes/providers/:name/resources", c.DeleteKubernetesResources)
	}
}
