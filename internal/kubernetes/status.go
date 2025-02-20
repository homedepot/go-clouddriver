package kubernetes

import (
	"strings"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/manifest"
)

// Status definitions of kinds can be found at
// https://github.com/spinnaker/clouddriver/tree/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler
func GetStatus(kind string, m map[string]interface{}) manifest.Status {
	var status manifest.Status

	switch strings.ToLower(kind) {
	case "daemonset":
		status = NewDaemonSet(m).Status()
	case "deployment":
		status = NewDeployment(m).Status()
	case "horizontalpodautoscaler":
		status = NewHorizontalPodAutoscaler(m).Status()
	case "job":
		status = NewJob(m).Status()
	case "pod":
		status = NewPod(m).Status()
	case "replicaset":
		status = NewReplicaSet(m).Status()
	case "statefulset":
		status = NewStatefulSet(m).Status()
	default:
		status = NewCustomKind(kind, m).Status()
	}

	return status
}
