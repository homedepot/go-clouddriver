package kubernetes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	kube "github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (cc *Controller) Delete(c *gin.Context, dm DeleteManifestRequest) {
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := dm.Location

	provider, err := cc.KubernetesProvider(dm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	if provider.Namespace != nil {
		namespace = *provider.Namespace
	}

	do := metav1.DeleteOptions{}

	if dm.Options.GracePeriodSeconds != nil {
		do.GracePeriodSeconds = dm.Options.GracePeriodSeconds
	}
	// There are two options that will set the deletion to cascade. Originally this was set when
	// `options.cascading = true`, however in later APIs they have this set in
	// `options.orphanDependants = false`.
	//
	// Using a pointer should prevent the older API versions from falsly cascading
	// when they don't intend to.
	propagationPolicy := v1.DeletePropagationOrphan
	if dm.Options.Cascading ||
		(dm.Options.OrphanDependants != nil && !*dm.Options.OrphanDependants) {
		propagationPolicy = v1.DeletePropagationForeground
	}

	do.PropagationPolicy = &propagationPolicy

	// Default to the static mode.
	mode := "static"
	if dm.Mode != "" {
		mode = dm.Mode
	}

	switch strings.ToLower(mode) {
	// Both dynamic and static use the same logic. For 'dynamic' the manifest has already been resolved and passed in.
	case "dynamic", "static":
		a := strings.Split(dm.ManifestName, " ")
		kind := a[0]
		name := a[1]

		err = provider.ValidateKindStatus(kind)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		gvr, err := provider.Client.GVRForKind(kind)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = provider.Client.DeleteResourceByKindAndNameAndNamespace(kind, name, namespace, do)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			TaskType:     clouddriver.TaskTypeDelete,
			Timestamp:    internal.CurrentTimeUTC(),
			APIGroup:     gvr.Group,
			Name:         name,
			Namespace:    namespace,
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         kind,
			SpinnakerApp: dm.App,
			Cluster:      kubernetes.Cluster(kind, name),
		}

		err = cc.SQLClient.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	case "label":
		if len(dm.LabelSelectors.Selectors) == 0 {
			clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("requested to delete manifests by label, but no label selectors were provided"))
			return
		}

		ls := labels.NewSelector()
		// The set of label selectors will be used across all Kinds.
		for _, selector := range dm.LabelSelectors.Selectors {
			req, err := kube.NewRequirement(selector.Kind, selector.Key, selector.Values)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			ls = ls.Add(*req)
		}

		lo := metav1.ListOptions{
			LabelSelector: ls.String(),
			FieldSelector: "metadata.namespace=" + namespace,
		}

		for _, kind := range dm.Kinds {
			err = provider.ValidateKindStatus(kind)
			if err != nil {
				clouddriver.Error(c, http.StatusBadRequest, err)
				return
			}

			gvr, err := provider.Client.GVRForKind(kind)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
			// Find list of resources with label selectors
			list, err := provider.Client.ListByGVR(gvr, lo)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			// Spinnaker OSS does not fail when no resources are found,
			// so create a 'no op' resource.
			if list == nil || len(list.Items) == 0 {
				kr := kubernetes.Resource{
					AccountName: dm.Account,
					ID:          uuid.New().String(),
					TaskID:      taskID,
					TaskType:    clouddriver.TaskTypeNoOp,
					Timestamp:   internal.CurrentTimeUTC(),
				}

				err = cc.SQLClient.CreateKubernetesResource(kr)
				if err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}

				continue
			}

			for _, item := range list.Items {
				// Delete each resource and record it in the database.
				err = provider.Client.DeleteResourceByKindAndNameAndNamespace(kind, item.GetName(), namespace, do)
				if err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}

				kr := kubernetes.Resource{
					AccountName:  dm.Account,
					ID:           uuid.New().String(),
					TaskID:       taskID,
					TaskType:     clouddriver.TaskTypeDelete,
					Timestamp:    internal.CurrentTimeUTC(),
					APIGroup:     gvr.Group,
					Name:         item.GetName(),
					Namespace:    namespace,
					Resource:     gvr.Resource,
					Version:      gvr.Version,
					Kind:         kind,
					SpinnakerApp: dm.App,
					Cluster:      kubernetes.Cluster(kind, item.GetName()),
				}

				err = cc.SQLClient.CreateKubernetesResource(kr)
				if err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}
			}
		}

		return
	default:
		clouddriver.Error(c, http.StatusNotImplemented,
			fmt.Errorf("requested to delete manifest %s using mode %s which is not implemented", dm.ManifestName, mode))
		return
	}
}
