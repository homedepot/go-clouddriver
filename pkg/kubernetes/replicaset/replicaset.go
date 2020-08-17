package replicaset

import (
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/mitchellh/mapstructure"
)

type replicaset struct {
	// APIVersion string `json:"apiVersion"`
	// Kind       string `json:"kind"`
	Metadata struct {
		// Annotations struct {
		// 	ArtifactSpinnakerIoLocation           string `json:"artifact.spinnaker.io/location"`
		// 	ArtifactSpinnakerIoName               string `json:"artifact.spinnaker.io/name"`
		// 	ArtifactSpinnakerIoType               string `json:"artifact.spinnaker.io/type"`
		// 	DeploymentKubernetesIoDesiredReplicas string `json:"deployment.kubernetes.io/desired-replicas"`
		// 	DeploymentKubernetesIoMaxReplicas     string `json:"deployment.kubernetes.io/max-replicas"`
		// 	DeploymentKubernetesIoRevision        string `json:"deployment.kubernetes.io/revision"`
		// 	MonikerSpinnakerIoApplication         string `json:"moniker.spinnaker.io/application"`
		// 	MonikerSpinnakerIoCluster             string `json:"moniker.spinnaker.io/cluster"`
		// } `json:"annotations"`
		// CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation int `json:"generation"`
		// Labels            struct {
		// 	App                      string `json:"app"`
		// 	AppKubernetesIoManagedBy string `json:"app.kubernetes.io/managed-by"`
		// 	AppKubernetesIoName      string `json:"app.kubernetes.io/name"`
		// 	PodTemplateHash          string `json:"pod-template-hash"`
		// } `json:"labels"`
		// Name            string `json:"name"`
		// Namespace       string `json:"namespace"`
		// OwnerReferences []struct {
		// 	APIVersion         string `json:"apiVersion"`
		// 	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
		// 	Controller         bool   `json:"controller"`
		// 	Kind               string `json:"kind"`
		// 	Name               string `json:"name"`
		// 	UID                string `json:"uid"`
		// } `json:"ownerReferences"`
		// ResourceVersion string `json:"resourceVersion"`
		// SelfLink        string `json:"selfLink"`
		// UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		// 	Selector struct {
		// 		MatchLabels struct {
		// 			App             string `json:"app"`
		// 			PodTemplateHash string `json:"pod-template-hash"`
		// 		} `json:"matchLabels"`
		// 	} `json:"selector"`
		// 	Template struct {
		// 		Metadata struct {
		// 			Annotations struct {
		// 				ArtifactSpinnakerIoLocation   string `json:"artifact.spinnaker.io/location"`
		// 				ArtifactSpinnakerIoName       string `json:"artifact.spinnaker.io/name"`
		// 				ArtifactSpinnakerIoType       string `json:"artifact.spinnaker.io/type"`
		// 				MonikerSpinnakerIoApplication string `json:"moniker.spinnaker.io/application"`
		// 				MonikerSpinnakerIoCluster     string `json:"moniker.spinnaker.io/cluster"`
		// 			} `json:"annotations"`
		// 			CreationTimestamp interface{} `json:"creationTimestamp"`
		// 			Labels            struct {
		// 				App                      string `json:"app"`
		// 				AppKubernetesIoManagedBy string `json:"app.kubernetes.io/managed-by"`
		// 				AppKubernetesIoName      string `json:"app.kubernetes.io/name"`
		// 				PodTemplateHash          string `json:"pod-template-hash"`
		// 			} `json:"labels"`
		// 		} `json:"metadata"`
		// 		Spec struct {
		// 			Containers []struct {
		// 				Image           string `json:"image"`
		// 				ImagePullPolicy string `json:"imagePullPolicy"`
		// 				Name            string `json:"name"`
		// 				Ports           []struct {
		// 					ContainerPort int    `json:"containerPort"`
		// 					Protocol      string `json:"protocol"`
		// 				} `json:"ports"`
		// 				Resources struct {
		// 				} `json:"resources"`
		// 				TerminationMessagePath   string `json:"terminationMessagePath"`
		// 				TerminationMessagePolicy string `json:"terminationMessagePolicy"`
		// 			} `json:"containers"`
		// 			DNSPolicy       string `json:"dnsPolicy"`
		// 			RestartPolicy   string `json:"restartPolicy"`
		// 			SchedulerName   string `json:"schedulerName"`
		// 			SecurityContext struct {
		// 			} `json:"securityContext"`
		// 			TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds"`
		// 		} `json:"spec"`
		// 	} `json:"template"`
	} `json:"spec"`
	Status struct {
		AvailableReplicas    int `json:"availableReplicas"`
		FullyLabeledReplicas int `json:"fullyLabeledReplicas"`
		ObservedGeneration   int `json:"observedGeneration"`
		ReadyReplicas        int `json:"readyReplicas"`
		Replicas             int `json:"replicas"`
	} `json:"status"`
}

func Status(m map[string]interface{}) kubernetes.ManifestStatus {
	s := kubernetes.DefaultStatus

	r := &replicaset{}
	if err := mapstructure.Decode(m, &r); err != nil {
		return kubernetes.NoneReported
	}

	desired := r.Spec.Replicas
	fullyLabeled := r.Status.FullyLabeledReplicas
	available := r.Status.AvailableReplicas
	ready := r.Status.ReadyReplicas

	if desired > fullyLabeled {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be fully-labeled"

		return s
	}

	if desired > ready {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be ready"

		return s
	}

	if desired > available {
		s.Stable.State = false
		s.Stable.Message = "Waiting for all replicas to be available"

		return s
	}

	if r.Metadata.Generation != r.Status.ObservedGeneration {
		s.Stable.State = false
		s.Stable.Message = "Waiting for replicaset spec update to be observed"

		return s
	}

	return s
}
