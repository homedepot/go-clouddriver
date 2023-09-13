package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
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
	"k8s.io/apimachinery/pkg/types"
)

var (
	errInvalidManifestName = errors.New("manifest name must be in format '{kind} {name}'")
)

// Disable takes in manifest coordinates and grabs the list of load balancers
// fronting it from the annotation `traffic.spinnaker.io/load-balancers`.
// It loops through these load balancers, removing any selectors from the target manifest's
// labels and patching the target resource using the JSON patch strategy. It then
// patches the labels of all pods that this manifest owns.
func (cc *Controller) Disable(c *gin.Context, dm DisableManifestRequest) {
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := dm.Location

	provider, err := cc.KubernetesProvider(dm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	// Preserve backwards compatibility
	if len(provider.Namespaces) == 1 {
		namespace = provider.Namespaces[0]
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

	err = provider.ValidateNamespaceAccess(namespace)
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

		err = attachDetach(provider.Client, lb, target, "remove")
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		// Patch all pods.
		for _, pod := range pods {
			err = attachDetach(provider.Client, lb, pod, "remove")
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// Just create one entry for a successful detachment of load balancers.
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

// getLoadBalancer gets a given load balancer from a specified namespace.
func getLoadBalancer(client kubernetes.Client, loadBalancer, namespace string) (*unstructured.Unstructured, error) {
	a := strings.Split(loadBalancer, " ")
	if len(a) != 2 {
		return nil, fmt.Errorf("failed to attach/detach to/from load balancer '%s'. "+
			"load balancers must be specified in the form '{kind} {name}', e.g. 'service my-service'", loadBalancer)
	}

	kind := a[0]
	name := a[1]
	// For now, limit the kind of load balancer available to detach to Services.
	if !strings.EqualFold(kind, "service") {
		return nil, fmt.Errorf("no support for load balancing via %s exists in Spinnaker", kind)
	}

	// Grab the load balancer from the cluster.
	lb, err := client.Get(kind, name, namespace)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("load balancer %s %s does not exist", kind, name)
		}

		return nil, fmt.Errorf("error getting service %s: %v", name, err)
	}

	return lb, nil
}

type jsonPatch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value,omitempty"`
}

// attachDetach attaches or detaches a given load balancer based on the op passed in ("add" or "remove")
// from a target and then patches the target's labels. It uses the JSON patch strategy to do so.
//
// For a pod, this looks like:
// [
//
//	{"op": "add", "path": "/metadata/labels/key1", "value": "value1"}
//
// ]
//
// For other kinds, this looks like:
// [
//
//	{"op": "remove", "path": "/spec/template/metadata/labels/key1"}
//
// ]
func attachDetach(client kubernetes.Client, lb, target *unstructured.Unstructured, op string) error {
	labelsPath := "spec.template.metadata.labels"

	labels, found, err := unstructured.NestedStringMap(target.Object, strings.Split(labelsPath, ".")...)
	if err != nil {
		return err
	}

	if !found {
		labelsPath = "metadata.labels"

		labels, _, err = unstructured.NestedStringMap(target.Object, strings.Split(labelsPath, ".")...)
		if err != nil {
			return err
		}
	}

	selector, found, _ := unstructured.NestedStringMap(lb.Object, "spec", "selector")
	if !found || len(selector) == 0 {
		// No selectors here so just return.
		return nil
	}

	if op == "add" && !disjoint(labels, selector) {
		return fmt.Errorf("service selector must have no label keys in common with target workload")
	}

	patches := []jsonPatch{}

	for k, v := range selector {
		// If the op is "remove" and the label does not exist, skip!
		if op == "remove" {
			if _, ok := labels[k]; !ok {
				continue
			}
		}

		// See https://datatracker.ietf.org/doc/html/rfc6901#section-3.
		k = strings.ReplaceAll(k, "~", "~0")
		k = strings.ReplaceAll(k, "/", "~1")

		patch := jsonPatch{
			Op:   op,
			Path: fmt.Sprintf("/%s/%s", strings.ReplaceAll(labelsPath, ".", "/"), k),
		}
		// If the op is "add" we set the value to be what we're adding.
		if op == "add" {
			patch.Value = v
		}

		patches = append(patches, patch)
	}

	// If there's nothing to patch, just return.
	if len(patches) == 0 {
		return nil
	}

	sort.Slice(patches, func(i, j int) bool { return patches[i].Path < patches[j].Path })

	b, err := json.Marshal(patches)
	if err != nil {
		return err
	}

	// Source code for Clouddriver always uses the JSON patch type for enable/disable manifest operations.
	// See https://github.com/spinnaker/clouddriver/blob/c52df8fb055de77ac800b41fd843761f506e7e08/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/manifest/AbstractKubernetesEnableDisableManifestOperation.java#L112.
	_, _, err = client.PatchUsingStrategy(target.GetKind(), target.GetName(),
		target.GetNamespace(), b, types.JSONPatchType)
	if err != nil {
		return err
	}

	return nil
}

// hasPods returns true if the kind of a Kubernetes object is
// - CronJob
// - DaemonSet
// - Job
// - ReplicaSet
// - StatefulSet
//
// This list is taken from the Disable (Manifest) stage in the Spinnaker UI.
func hasPods(u *unstructured.Unstructured) bool {
	return strings.EqualFold(u.GetKind(), "CronJob") ||
		strings.EqualFold(u.GetKind(), "DaemonSet") ||
		strings.EqualFold(u.GetKind(), "Job") ||
		strings.EqualFold(u.GetKind(), "ReplicaSet") ||
		strings.EqualFold(u.GetKind(), "StatefulSet")
}
