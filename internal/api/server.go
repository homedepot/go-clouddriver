package api

import (
	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/api/core"
	"github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	v1 "github.com/homedepot/go-clouddriver/internal/api/v1"
	"github.com/homedepot/go-clouddriver/internal/middleware"
)

// Server hold the gin engine and any clients we need for the API.
type Server struct {
	verboseRequestLogging bool
	e                     *gin.Engine
	c                     *internal.Controller
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
	s.verboseRequestLogging = true
}

// Setup sets any global middlewares then initializes the API.
func (s *Server) Setup() {
	s.e.Use(middleware.HandleError())

	if s.verboseRequestLogging {
		s.e.Use(middleware.LogRequest())
	}

	// API endpoints without a version will go under "core".
	{
		// Declare all controllers.
		mc := &middleware.Controller{s.c}
		kc := &kubernetes.Controller{s.c}
		ctlr := &core.Controller{s.c, kc}
		api := s.e.Group("")
		api.GET("/health", core.OK)

		// Force cache refresh.
		api.POST("/cache/kubernetes/manifest", core.OK)

		// Credentials API controller.
		api.GET("/credentials", ctlr.ListCredentials)
		api.GET("/credentials/:account", ctlr.GetAccountCredentials)

		// Applications API controller.
		//
		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ApplicationsController.groovy#L38
		// @PreAuthorize("#restricted ? @fiatPermissionEvaluator.storeWholePermission() : true")
		// @PostFilter("#restricted ? hasPermission(filterObject.name, 'APPLICATION', 'READ') : true")
		api.GET("/applications", mc.PostFilterAuthorizedApplications("READ"), ctlr.ListApplications)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupManagerController.java#L39
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostFilter("hasPermission(filterObject.account, 'ACCOUNT', 'READ')")
		api.GET("/applications/:application/serverGroupManagers", mc.AuthApplication("READ"), ctlr.ListServerGroupManagers)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupController.groovy#L172
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/serverGroups", mc.AuthApplication("READ"), ctlr.ListServerGroups)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ServerGroupController.groovy#L75
		// @PreAuthorize("hasPermission(#account, 'ACCOUNT', 'READ')")
		// @PostAuthorize("hasPermission(returnObject?.moniker?.app, 'APPLICATION', 'READ')")
		// textPayload: "Headers: map[Accept:[application/json] Accept-Encoding:[gzip] Connection:[Keep-Alive] User-Agent:[okhttp/3.14.9] X-Spinnaker-Accounts:[gke_github-replication-sandbox_us-east1_sandbox-us-east1-agent-dev,gke_github-replication-sandbox_us-east1_sandbox-us-east1-dev,gke_github-replication-sandbox_us-central1-c_prom-test] X-Spinnaker-Application:[smoketests] X-Spinnaker-User:[me@me.com]]"
		api.GET("/applications/:application/serverGroups/:account/:location/:name", mc.AuthAccount("READ"), ctlr.GetServerGroup)

		// https: //github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/LoadBalancerController.groovy#L42
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/loadBalancers", mc.AuthApplication("READ"), ctlr.ListLoadBalancers)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/ClusterController.groovy#L44
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission() and hasPermission(#application, 'APPLICATION', 'READ')")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/applications/:application/clusters", mc.AuthApplication("READ"), ctlr.ListClusters)

		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/JobController.groovy#L35
		// @PreAuthorize("hasPermission(#application, 'APPLICATION', 'READ') and hasPermission(#account, 'ACCOUNT', 'READ')")
		// @ApiOperation(value = "Collect a JobStatus", notes = "Collects the output of the job.")
		api.GET("/applications/:application/jobs/:account/:location/:name", mc.AuthApplication("READ"), mc.AuthAccount("READ"), ctlr.GetJob)
		// Delete job always fails, so we do not need to pass through the auth middlewares.
		api.DELETE("/applications/:application/jobs/:account/:location/:name", core.DeleteJob)

		// Create a kubernetes operation - deploy/delete/scale manifest.
		api.POST("/kubernetes/ops", middleware.TaskID(), ctlr.CreateKubernetesOperation)

		// Manifests API controller.
		api.GET("/manifests/:account/:location/:kind", ctlr.GetManifest)
		api.GET("/manifests/:account/:location/:kind/cluster/:application/:cluster/dynamic/:target", ctlr.GetManifestByTarget)

		// Instances API controller.
		api.GET("/instances/:account/:location/:name", ctlr.GetInstance)
		api.GET("/instances/:account/:location/:name/console", ctlr.GetInstanceConsole)

		// Get results for a task triggered in CreateKubernetesOperation.
		api.GET("/task/:id", ctlr.GetTask)

		// Generic search endpoint.
		//
		// https://github.com/spinnaker/clouddriver/blob/0524d08f6bcf775c469a0576a79b2679b5653325/clouddriver-web/src/main/groovy/com/netflix/spinnaker/clouddriver/controllers/SearchController.groovy#L55
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission()")
		api.GET("/search", ctlr.Search)

		// Not implemented.
		//
		// @PreAuthorize("@fiatPermissionEvaluator.storeWholePermission()")
		// @PostAuthorize("@authorizationSupport.filterForAccounts(returnObject)")
		api.GET("/securityGroups", core.ListSecurityGroups)

		// Artifacts API controller.
		api.GET("/artifacts/credentials", ctlr.ListArtifactCredentials)
		api.GET("/artifacts/account/:accountName/names", ctlr.ListHelmArtifactAccountNames)
		api.GET("/artifacts/account/:accountName/versions", ctlr.ListHelmArtifactAccountVersions)
		api.PUT("/artifacts/fetch/", ctlr.GetArtifact)

		// Features.
		api.GET("/features/stages", core.ListStages)
	}

	// V1 endpoint.
	{
		ctlr := &v1.Controller{s.c}
		api := s.e.Group("/v1")
		// Providers endpoint for kubernetes.
		api.GET("/kubernetes/providers", ctlr.ListKubernetesProvider)
		api.GET("/kubernetes/providers/:name", ctlr.GetKubernetesProvider)
		api.POST("/kubernetes/providers", ctlr.CreateKubernetesProvider)
		api.PUT("/kubernetes/providers", ctlr.CreateOrReplaceKubernetesProvider)
		api.DELETE("/kubernetes/providers/:name", ctlr.DeleteKubernetesProvider)
	}
}
