package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Enable takes in manifest coordinates and grabs the list of load balancers behind which
// it needs to be fronted from the annotation `traffic.spinnaker.io/load-balancers`.
// It loops through these load balancers, adding any selectors from the load balancer's
// labels and patching the target resource using the JSON patch strategy. It then
// patches the labels of all pods that this manifest owns.
func (cc *Controller) Enable(c *gin.Context, dm EnableManifestRequest) {
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

	// ManifestName is the kind and name of the manifest, including any version, like
	// 'ReplicaSet test-rs-v001'.
	a := strings.Split(dm.ManifestName, " ")
	if len(a) != 2 {
		clouddriver.Error(c, http.StatusBadRequest, errInvalidManifestName)
		return
	}

	kind := a[0]
	name := a[1]

	err = provider.ValidateKindStatus(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	// Grab the target manifest.
	target, err := provider.Client.Get(kind, name, namespace)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			clouddriver.Error(c, http.StatusNotFound, fmt.Errorf("resource %s %s does not exist", kind, name))
			return
		}

		clouddriver.Error(c, http.StatusInternalServerError, fmt.Errorf("error getting resource (kind: %s, name: %s, namespace: %s): %v",
			kind, name, namespace, err))

		return
	}

	loadBalancers, err := kubernetes.LoadBalancers(*target)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	var pods []*unstructured.Unstructured
	// If the target manifest has load balancers and pods, list pods, grab those that has the owner UID
	// of the target manifest, and patch those pods.
	if len(loadBalancers) > 0 && hasPods(target) {
		// Declare server side filtering options.
		lo := metav1.ListOptions{
			FieldSelector: "metadata.namespace=" + namespace,
			LabelSelector: kubernetes.DefaultLabelSelector(),
		}
		// Declare a context with timeout.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
		defer cancel()
		// List resources with the context.
		ul, err := provider.Client.ListResourceWithContext(ctx, "pods", lo)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
		// Loop through all pods, finding all that are owned by the target manifest.
		for _, u := range ul.Items {
			for _, ownerReference := range u.GetOwnerReferences() {
				if ownerReference.UID == target.GetUID() {
					// Create a copy of the unstructured object since we access by reference.
					u := u
					pods = append(pods, &u)
				}
			}
		}
	}

	for _, loadBalancer := range loadBalancers {
		lb, err := getLoadBalancer(provider.Client, loadBalancer, namespace)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = attachDetach(provider.Client, lb, target, "add")
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		// Patch all pods.
		for _, pod := range pods {
			err = attachDetach(provider.Client, lb, pod, "add")
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// Just create one entry for a successful attachment to load balancers.
	kr := kubernetes.Resource{
		TaskType:     clouddriver.TaskTypeNoOp,
		AccountName:  dm.Account,
		SpinnakerApp: dm.App,
		ID:           uuid.New().String(),
		TaskID:       taskID,
		Name:         name,
		Namespace:    namespace,
		Kind:         kind,
	}

	err = cc.SQLClient.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
}
