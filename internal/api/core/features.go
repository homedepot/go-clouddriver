package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Stages []Stage

type Stage struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

var stages = []string{
	"resizeServerGroup",
	"runJob",
	"undoRolloutManifest",
	"rollingRestartManifest",
	"pauseRolloutManifest",
	"enableManifest",
	"scaleManifest",
	"disableManifest",
	"patchManifest",
	"resumeRolloutManifest",
	"deleteManifest",
	"deployManifest",
	"cleanupArtifacts",
	"upsertLoadBalancer",
	"enableServerGroup",
	"createServerGroup",
	"deleteLoadBalancer",
	"upsertScalingPolicy",
	"terminateInstances",
	"stopServerGroup",
	"disableServerGroup",
	"startServerGroup",
	"destroyServerGroup",
}

// Expected response:
//
// [
//   {
//     "enabled": true,
//     "name": "resizeServerGroup"
//   },
//   {
//     "enabled": true,
//     "name": "runJob"
//   },
//   {
//     "enabled": true,
//     "name": "undoRolloutManifest"
//   },
//   {
//     "enabled": true,
//     "name": "rollingRestartManifest"
//   },
//   {
//     "enabled": true,
//     "name": "pauseRolloutManifest"
//   },
//   {
//     "enabled": true,
//     "name": "enableManifest"
//   },
//   {
//     "enabled": true,
//     "name": "scaleManifest"
//   },
//   {
//     "enabled": true,
//     "name": "disableManifest"
//   },
//   {
//     "enabled": true,
//     "name": "patchManifest"
//   },
//   {
//     "enabled": true,
//     "name": "resumeRolloutManifest"
//   },
//   {
//     "enabled": true,
//     "name": "deleteManifest"
//   },
//   {
//     "enabled": true,
//     "name": "deployManifest"
//   },
//   {
//     "enabled": true,
//     "name": "cleanupArtifacts"
//   },
//   {
//     "enabled": true,
//     "name": "upsertLoadBalancer"
//   },
//   {
//     "enabled": true,
//     "name": "enableServerGroup"
//   },
//   {
//     "enabled": true,
//     "name": "createServerGroup"
//   },
//   {
//     "enabled": true,
//     "name": "deleteLoadBalancer"
//   },
//   {
//     "enabled": true,
//     "name": "upsertScalingPolicy"
//   },
//   {
//     "enabled": true,
//     "name": "terminateInstances"
//   },
//   {
//     "enabled": true,
//     "name": "stopServerGroup"
//   },
//   {
//     "enabled": true,
//     "name": "disableServerGroup"
//   },
//   {
//     "enabled": true,
//     "name": "startServerGroup"
//   },
//   {
//     "enabled": true,
//     "name": "destroyServerGroup"
//   }
// ]
func ListStages(c *gin.Context) {
	response := Stages{}

	for _, stage := range stages {
		s := Stage{
			Enabled: false,
			Name:    stage,
		}
		response = append(response, s)
	}

	c.JSON(http.StatusOK, response)
}
