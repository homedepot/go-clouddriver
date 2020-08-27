package kubernetes

import (
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
)

// Status definitions of kinds can be found at
// https://github.com/spinnaker/clouddriver/tree/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler
func GetStatus(kind string, m map[string]interface{}) manifest.Status {
	status := manifest.DefaultStatus

	// TODO need to fill in statuses for all kinds here.
	switch strings.ToLower(kind) {
	case "deployment":
		status = NewDeployment(m).Status()
	case "pod":
		status = NewPod(m).Status()
	case "replicaset":
		status = NewReplicaSet(m).Status()
	}

	return status
}
