package kubernetes

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (cc *Controller) CleanupArtifacts(c *gin.Context, ca CleanupArtifactsRequest) {
	app := c.GetHeader("X-Spinnaker-Application")
	taskID := clouddriver.TaskIDFromContext(c)

	for _, manifest := range ca.Manifests {
		u, err := kubernetes.ToUnstructured(manifest)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		provider, err := cc.KubernetesProvider(ca.Account)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		gvr, err := provider.Client.GVRForKind(u.GetKind())
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		namespace := u.GetNamespace()

		// Preserve backwards compatibility
		if len(provider.Namespaces) == 1 {
			namespace = provider.Namespaces[0]
		}

		err = provider.ValidateNamespaceAccess(namespace)
		if err != nil {
			clouddriver.Log(err)
			continue
		}

		// Grab the cluster of this resource from its annotations.
		cluster := clusterAnnotation(u)
		// Handle max version history. Source code here:
		// https://github.com/spinnaker/clouddriver/blob/master/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/artifact/KubernetesCleanupArtifactsOperation.java#L102
		maxVersionHistory, err := kubernetes.MaxVersionHistory(u)
		if err == nil && maxVersionHistory > 0 && cluster != "" {
			// Only list resources that are managed by Spinnaker.
			lo := metav1.ListOptions{
				LabelSelector: kubernetes.DefaultLabelSelector(),
			}

			ul, err := provider.Client.ListResourcesByKindAndNamespace(u.GetKind(), namespace, lo)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError,
					fmt.Errorf("error listing resources to cleanup for max version history (kind: %s, name: %s, namespace: %s): %v",
						u.GetKind(), u.GetName(), namespace, err))

				return
			}

			artifacts := kubernetes.FilterOnAnnotation(ul.Items,
				kubernetes.AnnotationSpinnakerMonikerCluster, cluster)
			if maxVersionHistory < len(artifacts) {
				// Sort on creation timestamp oldest to newest.
				sort.Slice(artifacts, func(i, j int) bool {
					return artifacts[i].GetCreationTimestamp().String() < artifacts[j].GetCreationTimestamp().String()
				})

				artifactsToDelete := artifacts[0 : len(artifacts)-maxVersionHistory]
				for _, a := range artifactsToDelete {
					// Delete the resource and any dependants in the foreground.
					pp := v1.DeletePropagationForeground
					do := metav1.DeleteOptions{
						PropagationPolicy: &pp,
					}

					err = provider.Client.DeleteResourceByKindAndNameAndNamespace(a.GetKind(), a.GetName(), namespace, do)
					if err != nil {
						clouddriver.Error(c, http.StatusInternalServerError,
							fmt.Errorf("error deleting resource to cleanup for max version history (kind: %s, name: %s, namespace: %s): %v",
								a.GetKind(), a.GetName(), namespace, err))

						return
					}
				}
			}
		}

		kr := kubernetes.Resource{
			AccountName:  ca.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			TaskType:     clouddriver.TaskTypeCleanup,
			Timestamp:    internal.CurrentTimeUTC(),
			APIGroup:     gvr.Group,
			Name:         u.GetName(),
			Namespace:    namespace,
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         u.GetKind(),
			SpinnakerApp: app,
			Cluster:      cluster,
		}

		err = cc.SQLClient.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	}
}

// clusterAnnotation returns the value of the annotation
// 'moniker.spinnaker.io/cluster' if it exists, otherwise
// it returns an empty string.
func clusterAnnotation(u unstructured.Unstructured) string {
	cluster := ""

	annotations := u.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[kubernetes.AnnotationSpinnakerMonikerCluster]; ok {
			cluster = value
		}
	}

	return cluster
}
