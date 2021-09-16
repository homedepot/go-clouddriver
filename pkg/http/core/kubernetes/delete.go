package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/homedepot/go-clouddriver/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

func Delete(c *gin.Context, dm DeleteManifestRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := sc.GetKubernetesProvider(dm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	token, err := ac.Token(provider.TokenProvider)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
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

		gvr, err := client.GVRForKind(kind)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = client.DeleteResourceByKindAndNameAndNamespace(kind, name, dm.Location, do)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			TaskType:     clouddriver.TaskTypeDelete,
			Timestamp:    util.CurrentTimeUTC(),
			APIGroup:     gvr.Group,
			Name:         name,
			Namespace:    dm.Location,
			Resource:     gvr.Resource,
			Version:      gvr.Version,
			Kind:         kind,
			SpinnakerApp: dm.App,
			Cluster:      cluster(kind, name),
		}

		err = sc.CreateKubernetesResource(kr)
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
			FieldSelector: "metadata.namespace=" + dm.Location,
		}

		for _, kind := range dm.Kinds {
			gvr, err := client.GVRForKind(kind)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
			// Find list of resources with label selectors
			list, err := client.ListByGVR(gvr, lo)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			// Spinnaker OSS does not fail when no resources are found,
			// so create an 'no op' resource.
			if list == nil || len(list.Items) == 0 {
				kr := kubernetes.Resource{
					AccountName: dm.Account,
					ID:          uuid.New().String(),
					TaskID:      taskID,
					TaskType:    clouddriver.TaskTypeNoOp,
					Timestamp:   util.CurrentTimeUTC(),
				}

				err = sc.CreateKubernetesResource(kr)
				if err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}

				continue
			}

			for _, item := range list.Items {
				// Delete each resource and record it in the database.
				err = client.DeleteResourceByKindAndNameAndNamespace(kind, item.GetName(), dm.Location, do)
				if err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}

				kr := kubernetes.Resource{
					AccountName:  dm.Account,
					ID:           uuid.New().String(),
					TaskID:       taskID,
					TaskType:     clouddriver.TaskTypeDelete,
					Timestamp:    util.CurrentTimeUTC(),
					APIGroup:     gvr.Group,
					Name:         item.GetName(),
					Namespace:    dm.Location,
					Resource:     gvr.Resource,
					Version:      gvr.Version,
					Kind:         kind,
					SpinnakerApp: dm.App,
					Cluster:      cluster(kind, item.GetName()),
				}

				err = sc.CreateKubernetesResource(kr)
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
