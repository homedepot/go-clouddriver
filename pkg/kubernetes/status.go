package kubernetes

import (
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/deployment"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/manifest"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/pod"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/replicaset"
)

// Status definitions of kinds can be found at
// https://github.com/spinnaker/clouddriver/tree/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler
func GetStatus(kind string, m map[string]interface{}) manifest.Status {
	status := manifest.DefaultStatus

	switch strings.ToLower(kind) {
	case "deployment":
		status = deployment.Status(m)
	case "pod":
		status = pod.Status(m)
	case "replicaset":
		status = replicaset.Status(m)
	}

	return status
}
