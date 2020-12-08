package kubernetes

import (
	"encoding/json"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/manifest"
	v1 "k8s.io/api/batch/v1"
)

type Job interface {
	Status() manifest.Status
	State() string
	Object() *v1.Job
}

func NewJob(m map[string]interface{}) Job {
	j := &v1.Job{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &j)

	return &job{j: j}
}

type job struct {
	j *v1.Job
}

func (j *job) Object() *v1.Job {
	return j.j
}

// Calculated at https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/model/KubernetesJobStatus.java#L71
func (j *job) State() string {
	obj := j.Object()
	status := obj.Status

	if status.CompletionTime == nil {
		return "Running"
	}

	completions := int32(1)
	if obj.Spec.Completions != nil {
		completions = *obj.Spec.Completions
	}

	succeeded := status.Succeeded

	if succeeded < completions {
		conditions := status.Conditions
		failed := false

		for _, condition := range conditions {
			if condition.Type == v1.JobFailed {
				failed = true
				break
			}
		}

		if failed {
			return "Failed"
		}

		return "Running"
	}

	return "Succeeded"
}

func (j *job) Status() manifest.Status {
	s := manifest.DefaultStatus

	completions := int32(1)
	spec := j.j.Spec
	status := j.j.Status

	if spec.Completions != nil {
		completions = *spec.Completions
	}

	succeeded := status.Succeeded
	if succeeded < completions {
		conditions := status.Conditions
		for _, condition := range conditions {
			if condition.Type == v1.JobFailed {
				s.Failed.State = true
				return s
			}
		}

		s.Stable.State = false
		s.Stable.Message = "Waiting for jobs to finish"
	}

	return s
}
