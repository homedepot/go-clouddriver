package kubernetes

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	// RevisionAnnotation is the revision annotation of a deployment's replica sets which records its rollout sequence
	RevisionAnnotation = "deployment.kubernetes.io/revision"
	// RevisionHistoryAnnotation maintains the history of all old revisions that a replica set has served for a deployment.
	RevisionHistoryAnnotation = "deployment.kubernetes.io/revision-history"
	// DesiredReplicasAnnotation is the desired replicas for a deployment recorded as an annotation
	// in its replica sets. Helps in separating scaling events from the rollout process and for
	// determining if the new replica set for a deployment is really saturated.
	DesiredReplicasAnnotation = "deployment.kubernetes.io/desired-replicas"
	// MaxReplicasAnnotation is the maximum replicas a deployment can have at a given point, which
	// is deployment.spec.replicas + maxSurge. Used by the underlying replica sets to estimate their
	// proportions in case the deployment has surge replicas.
	MaxReplicasAnnotation = "deployment.kubernetes.io/max-replicas"
)

var (
	errNoApplicationProvided = errors.New("no application provided")
	errRevisionNotFound      = errors.New("revision not found")
)

func UndoRolloutManifest(c *gin.Context, urm UndoRolloutManifestRequest) error {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)
	application := c.GetHeader("X-Spinnaker-Application")
	a := strings.Split(urm.ManifestName, " ")
	manifestKind := a[0]
	manifestName := a[1]

	if application == "" {
		return errNoApplicationProvided
	}

	provider, err := sc.GetKubernetesProvider(urm.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: provider.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	kc.SetDynamicClientForConfig(config)

	d, err := kc.Get(manifestKind, manifestName, urm.Location)
	if err != nil {
		return err
	}

	replicaSetGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "replicasets",
	}
	lo := metav1.ListOptions{
		LabelSelector: kubernetes.LabelKubernetesSpinnakerApp + "=" + application,
	}

	replicaSets, err := kc.List(replicaSetGVR, lo)
	if err != nil {
		return err
	}

	var targetRS *unstructured.Unstructured

	// Deployments manage replicasets, so build a list of managed replicasets for each deployment.
	for _, replicaSet := range replicaSets.Items {
		annotations := replicaSet.GetAnnotations()
		if annotations != nil {
			name := annotations[kubernetes.AnnotationSpinnakerArtifactName]
			t := annotations[kubernetes.AnnotationSpinnakerArtifactType]
			if name != "" && t != "" &&
				strings.EqualFold(name, manifestName) &&
				strings.EqualFold(t, "kubernetes/"+manifestKind) {
				replicaSetAnnotations := replicaSet.GetAnnotations()
				sequence := replicaSetAnnotations["deployment.kubernetes.io/revision"]
				if sequence != "" && sequence == urm.Revision {
					targetRS = &replicaSet
					break
				}
			}
		}
	}

	deployment := kubernetes.NewDeployment(d.Object)
	rs := kubernetes.NewReplicaSet(targetRS.Object).Object()

	if rs == nil {
		return errRevisionNotFound
	}

	SetFromReplicaSetTemplate(deployment.Object(), rs.Spec.Template)
	// set RS (the old RS we'll rolling back to) annotations back to the deployment;
	// otherwise, the deployment's current annotations (should be the same as current new RS)
	// will be copied to the RS after the rollback.
	//
	// For example,
	// A Deployment has old RS1 with annotation {change-cause:create}, and new RS2 {change-cause:edit}.
	// Note that both annotations are copied from Deployment, and the Deployment
	// should be annotated {change-cause:edit} as well.
	// Now, rollback Deployment to RS1, we should update Deployment's pod-template and also copy annotation from RS1.
	// Deployment is now annotated {change-cause:create}, and we have new
	// RS1 {change-cause:create}, old RS2 {change-cause:edit}.
	//
	// If we don't copy the annotations back from RS to deployment on rollback,
	// the Deployment will stay as {change-cause:edit},
	// and new RS1 becomes {change-cause:edit} (copied from deployment after rollback),
	// old RS2 {change-cause:edit}, which is not correct.
	SetDeploymentAnnotationsTo(deployment.Object(), rs)

	u, err := deployment.ToUnstructured()
	if err != nil {
		return err
	}

	_, err = kc.Apply(&u)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/util/deployment_util.go#L679
//
// SetFromReplicaSetTemplate sets the desired PodTemplateSpec from a replica set template to the given deployment.
func SetFromReplicaSetTemplate(deployment *apps.Deployment, template v1.PodTemplateSpec) *apps.Deployment {
	deployment.Spec.Template.ObjectMeta = template.ObjectMeta
	deployment.Spec.Template.Spec = template.Spec
	deployment.Spec.Template.ObjectMeta.Labels = CloneAndRemoveLabel(
		deployment.Spec.Template.ObjectMeta.Labels,
		apps.DefaultDeploymentUniqueLabelKey)

	return deployment
}

// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/util/deployment_util.go#L331
//
// SetDeploymentAnnotationsTo sets deployment's annotations as given RS's annotations.
// This action should be done if and only if the deployment is rolling back to this rs.
// Note that apply and revision annotations are not changed.
func SetDeploymentAnnotationsTo(deployment *apps.Deployment, rollbackToRS *apps.ReplicaSet) {
	deployment.Annotations = getSkippedAnnotations(deployment.Annotations)

	for k, v := range rollbackToRS.Annotations {
		if !skipCopyAnnotation(k) {
			deployment.Annotations[k] = v
		}
	}
}

func getSkippedAnnotations(annotations map[string]string) map[string]string {
	skippedAnnotations := make(map[string]string)

	for k, v := range annotations {
		if skipCopyAnnotation(k) {
			skippedAnnotations[k] = v
		}
	}

	return skippedAnnotations
}

// skipCopyAnnotation returns true if we should skip copying the annotation with the given annotation key.
func skipCopyAnnotation(key string) bool {
	return annotationsToSkip[key]
}

var annotationsToSkip = map[string]bool{
	v1.LastAppliedConfigAnnotation: true,
	RevisionAnnotation:             true,
	RevisionHistoryAnnotation:      true,
	DesiredReplicasAnnotation:      true,
	MaxReplicasAnnotation:          true,
	apps.DeprecatedRollbackTo:      true,
}

// Taken from https://github.com/kubernetes/kubernetes/blob/master/pkg/util/labels/labels.go
//
// CloneAndRemoveLabel clones the given map and returns a new map with the given key removed.
// Returns the given map, if labelKey is empty.
func CloneAndRemoveLabel(labels map[string]string, labelKey string) map[string]string {
	if labelKey == "" {
		// Don't need to add a label.
		return labels
	}
	// Clone.
	newLabels := map[string]string{}
	for key, value := range labels {
		newLabels[key] = value
	}

	delete(newLabels, labelKey)

	return newLabels
}
