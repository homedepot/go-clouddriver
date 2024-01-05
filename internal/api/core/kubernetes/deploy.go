package kubernetes

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	listTimeout = int64(30)
)

// Deploy performs a "Deploy (Manifest)" Spinnaker operation.
// It takes in a list of manifest and the Kubernetes provider
// to "apply" them to. It adds Spinnaker annotations/labels,
// and handles any Spinnaker versioning, then applies each manifest
// one by one.
func (cc *Controller) Deploy(c *gin.Context, dm DeployManifestRequest) {
	taskID := clouddriver.TaskIDFromContext(c)
	namespace := strings.TrimSpace(dm.NamespaceOverride)

	provider, err := cc.KubernetesProvider(dm.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	// Preserve backwards compatibility
	if len(provider.Namespaces) == 1 {
		namespace = provider.Namespaces[0]
	}

	// First, convert all manifests to unstructured objects.
	manifests, err := toUnstructured(dm.Manifests)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}
	// Merge all list element items into the manifest list.
	manifests, err = mergeManifests(manifests)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}
	// Sort the manifests by their kind's priority.
	manifests = kubernetes.SortManifests(manifests)

	// Set the namespace on all manifests.
	for _, manifest := range manifests {
		kubernetes.SetNamespaceOnManifest(&manifest, namespace)
	}

	application := dm.Moniker.App
	// Consolidate all deploy manifest request artifacts.
	artifacts := []clouddriver.Artifact{}
	artifacts = append(artifacts, dm.RequiredArtifacts...)
	artifacts = append(artifacts, dm.OptionalArtifacts...)

	for _, manifest := range manifests {
		// Create a copy of the unstructured object since we access by reference.
		manifest := manifest

		log.Println("before processing ", manifest, manifest.GetAnnotations())

		err = provider.ValidateKindStatus(manifest.GetKind())
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		if namespace == "" {
			err = provider.ValidateNamespaceAccess(manifest.GetNamespace()) // pass in the current manifest's namespace
		} else {
			err = provider.ValidateNamespaceAccess(namespace)
		}

		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		nameWithoutVersion := manifest.GetName()
		// If the kind is a job, its name is not set, and generateName is set,
		// generate a name for the job as `apply` will throw the error
		// `resource name may not be empty`.
		if strings.EqualFold(manifest.GetKind(), "job") && nameWithoutVersion == "" {
			generateName := manifest.GetGenerateName()
			manifest.SetName(generateName + rand.String(randNameNumber))
		}

		err = kubernetes.AddSpinnakerAnnotations(&manifest, application)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		err = kubernetes.AddSpinnakerLabels(&manifest, application)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		log.Println("after spinnaker annotations ", manifest, manifest.GetAnnotations())

		kubernetes.BindArtifacts(&manifest, artifacts, dm.Account)

		if kubernetes.IsVersioned(manifest) {
			err := handleVersionedManifest(provider.Client, &manifest, application)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		if kubernetes.UseSourceCapacity(manifest) {
			err = handleUseSourceCapacity(provider.Client, &manifest)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		if kubernetes.Recreate(manifest) {
			err = handleRecreate(provider.Client, &manifest)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		// Set the `traffic.spinnaker.io/load-balancers` annotation if the user
		// has requested for Spinnaker to manage a resources traffic.
		if dm.TrafficManagement.Enabled {
			err = handleTrafficManagement(&manifest, dm.TrafficManagement)
			if err != nil {
				clouddriver.Error(c, http.StatusBadRequest, err)
				return
			}
		}

		// Only handle attaching load balancers if not using Spinnaker traffic management
		// (i.e. user is manually setting the `traffic.spinnaker.io/load-balancers` annotation)
		// or if using Spinnaker traffic management and the user has requested to route traffic
		// to pods.
		if !dm.TrafficManagement.Enabled || (dm.TrafficManagement.Enabled && dm.TrafficManagement.Options.EnableTraffic) {
			err = handleAttachingLoadBalancers(provider.Client, &manifest, manifests)
			if err != nil {
				clouddriver.Error(c, http.StatusBadRequest, err)
				return
			}
		}

		meta := kubernetes.Metadata{}
		if kubernetes.Replace(manifest) {
			meta, err = provider.Client.Replace(&manifest)
			if err != nil {
				e := fmt.Errorf("error replacing manifest (kind: %s, apiVersion: %s, name: %s): %s",
					manifest.GetKind(), manifest.GroupVersionKind().Version, manifest.GetName(), err.Error())
				clouddriver.Error(c, http.StatusInternalServerError, e)

				return
			}
		} else {
			meta, err = provider.Client.Apply(&manifest)
			if err != nil {
				e := fmt.Errorf("error applying manifest (kind: %s, apiVersion: %s, name: %s): %s",
					manifest.GetKind(), manifest.GroupVersionKind().Version, manifest.GetName(), err.Error())
				clouddriver.Error(c, http.StatusInternalServerError, e)

				return
			}
		}

		kr := kubernetes.Resource{
			AccountName:  dm.Account,
			ID:           uuid.New().String(),
			TaskID:       taskID,
			Timestamp:    internal.CurrentTimeUTC(),
			APIGroup:     meta.Group,
			Name:         meta.Name,
			ArtifactName: nameWithoutVersion,
			Namespace:    meta.Namespace,
			Resource:     meta.Resource,
			Version:      meta.Version,
			Kind:         meta.Kind,
			SpinnakerApp: dm.Moniker.App,
			Cluster:      kubernetes.Cluster(meta.Kind, nameWithoutVersion),
		}

		annotations := manifest.GetAnnotations()
		artifactType := annotations[kubernetes.AnnotationSpinnakerArtifactType]
		artifact := clouddriver.Artifact{
			Name:      nameWithoutVersion,
			Reference: meta.Name,
			Type:      artifact.Type(artifactType),
		}
		artifacts = append(artifacts, artifact)

		err = cc.SQLClient.CreateKubernetesResource(kr)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
	}
}

// toUnstructured converts a slice of map[string]interface{} to unstructured.Unstructured.
func toUnstructured(manifests []map[string]interface{}) ([]unstructured.Unstructured, error) {
	m := []unstructured.Unstructured{}

	for _, manifest := range manifests {
		u, err := kubernetes.ToUnstructured(manifest)
		if err != nil {
			return nil, fmt.Errorf("kubernetes: unable to convert manifest to unstructured: %w", err)
		}

		m = append(m, u)
	}

	return m, nil
}

func getListOptions(app string) (metav1.ListOptions, error) {
	labelSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      kubernetes.LabelSpinnakerMonikerSequence,
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}

	ls, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		return metav1.ListOptions{}, err
	}

	lo := metav1.ListOptions{
		LabelSelector:  ls.String(),
		TimeoutSeconds: &listTimeout,
	}

	return lo, err
}

// mergeManifests merges manifests of kind List into the parent list of manifests.
func mergeManifests(manifests []unstructured.Unstructured) ([]unstructured.Unstructured, error) {
	mergedManifests := []unstructured.Unstructured{}

	for _, manifest := range manifests {
		if manifest.IsList() {
			ul, err := manifest.ToList()
			if err != nil {
				return nil, fmt.Errorf("error converting manifest to list: %w", err)
			}

			mergedManifests = append(mergedManifests, ul.Items...)
		} else {
			mergedManifests = append(mergedManifests, manifest)
		}
	}

	return mergedManifests, nil
}

func handleVersionedManifest(client kubernetes.Client, u *unstructured.Unstructured, application string) error {
	lo, err := getListOptions(application)
	if err != nil {
		return err
	}

	kind := strings.ToLower(u.GetKind())
	namespace := u.GetNamespace()

	results, err := client.ListResourcesByKindAndNamespace(kind, namespace, lo)
	if err != nil {
		return err
	}

	// Filter results to only those associated with this application.
	results.Items = kubernetes.FilterOnAnnotation(results.Items,
		kubernetes.AnnotationSpinnakerMonikerApplication, application)
	nameWithoutVersion := u.GetName()
	currentVersion := kubernetes.GetCurrentVersion(results, kind, nameWithoutVersion)
	latestVersion := kubernetes.IncrementVersion(currentVersion)
	u.SetName(nameWithoutVersion + "-" + latestVersion.Long)

	err = kubernetes.AddSpinnakerVersionAnnotations(u, latestVersion)
	if err != nil {
		return err
	}

	err = kubernetes.AddSpinnakerVersionLabels(u, latestVersion)
	if err != nil {
		return err
	}

	return nil
}

func handleUseSourceCapacity(client kubernetes.Client, u *unstructured.Unstructured) error {
	current, err := client.Get(u.GetKind(), u.GetName(), u.GetNamespace())
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}
	// If the resource is currently deployed then replace the replicas value
	// with the current value, if it has one
	if current != nil {
		r, found, err := unstructured.NestedInt64(current.Object, "spec", "replicas")
		if err != nil {
			return err
		}

		if found {
			err = unstructured.SetNestedField(u.Object, r, "spec", "replicas")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func handleRecreate(kubeClient kubernetes.Client, u *unstructured.Unstructured) error {
	current, err := kubeClient.Get(u.GetKind(), u.GetName(), u.GetNamespace())
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}

	// If the resource is currently deployed, then delete the resource prior to deploying.
	if current != nil {
		err := kubeClient.DeleteResourceByKindAndNameAndNamespace(u.GetKind(), u.GetName(), u.GetNamespace(), metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// handleTrafficManagement sets the `traffic.spinnaker.io/load-balancers`
// annotation accordingly if not set and errors if it is set.
func handleTrafficManagement(target *unstructured.Unstructured, tm TrafficManagement) error {
	annotations := target.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	if value, ok := annotations[kubernetes.AnnotationSpinnakerTrafficLoadBalancers]; ok && value != "" {
		return fmt.Errorf("manifest already has traffic.spinnaker.io/load-balancers annotation set to %s. "+
			"Failed attempting to set it to [%s]", value, strings.Join(tm.Options.Services, ", "))
	}

	loadBalancers := "["
	for i, service := range tm.Options.Services {
		loadBalancers += `"` + service + `"`
		if i < len(tm.Options.Services)-1 {
			loadBalancers += `, `
		}
	}

	loadBalancers += "]"
	annotations[kubernetes.AnnotationSpinnakerTrafficLoadBalancers] = loadBalancers
	target.SetAnnotations(annotations)

	return nil
}

// handleAttachingLoadBalancers grabs load balancers from a target manifests
// `traffic.spinnaker.io/load-balancers` annotation and attaches that load
// balancers selectors to the pod template labels of the target.
//
// See https://github.com/spinnaker/clouddriver/blob/62325f922533d9e96b35d88698959def4ad517b5/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/manifest/KubernetesDeployManifestOperation.java#L281
func handleAttachingLoadBalancers(client kubernetes.Client, target *unstructured.Unstructured,
	manifests []unstructured.Unstructured) error {
	lbs, err := kubernetes.LoadBalancers(*target)
	if err != nil {
		return err
	}

	for _, lb := range lbs {
		err = attachLoadBalancer(client, lb, target, manifests)
		if err != nil {
			return err
		}
	}

	return nil
}

// attachLoadBalancer modifies the labels of the target manifest to include the selectors of the load balancer.
func attachLoadBalancer(client kubernetes.Client, loadBalancer string,
	target *unstructured.Unstructured, manifests []unstructured.Unstructured) error {
	a := strings.Split(loadBalancer, " ")
	if len(a) != 2 {
		return fmt.Errorf("Failed to attach load balancer '%s'. "+
			"Load balancers must be specified in the form '{kind} {name}', e.g. 'service my-service'.", loadBalancer)
	}

	kind := a[0]
	name := a[1]
	// For now, limit the kind of load balancer available to attach to Services.
	if !strings.EqualFold(kind, "service") {
		// https://github.com/spinnaker/clouddriver/blob/8c377ef6be07278cd8a54448980f2b2065069a34/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/handler/CanLoadBalance.java#L39
		return fmt.Errorf("No support for load balancing via %s exists in Spinnaker.", kind)
	}

	var (
		lb    unstructured.Unstructured
		found bool
	)

	// First, see if the load balancer exists in the current request's manifests.
	for _, manifest := range manifests {
		if strings.EqualFold(manifest.GetKind(), "service") &&
			strings.EqualFold(manifest.GetName(), name) &&
			strings.EqualFold(target.GetNamespace(), manifest.GetNamespace()) {
			lb = manifest
			found = true
		}

		if found {
			break
		}
	}
	// If the manifest does not exist in the current request, get it from the
	// cluster.
	if !found {
		result, err := client.Get(kind, name, target.GetNamespace())
		if err != nil {
			if errors.IsNotFound(err) {
				// https://github.com/spinnaker/clouddriver/blob/62325f922533d9e96b35d88698959def4ad517b5/clouddriver-kubernetes/src/main/java/com/netflix/spinnaker/clouddriver/kubernetes/op/manifest/KubernetesDeployManifestOperation.java#L329
				return fmt.Errorf("Load balancer %s %s does not exist", kind, name)
			}

			return fmt.Errorf("error getting service %s: %v", name, err)
		}

		lb = *result
	}

	if err := attach(lb, target); err != nil {
		return err
	}

	return nil
}

// attach grabs the labels from the target manifest and appends the selectors
// from the passed in load balancer.
func attach(lb unstructured.Unstructured, target *unstructured.Unstructured) error {
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
		return fmt.Errorf("Service must have a non-empty selector in order to be attached to a workload")
	}

	if !disjoint(labels, selector) {
		return fmt.Errorf("Service selector must have no label keys in common with target workload")
	}

	for k, v := range selector {
		labels[k] = v
	}

	err = unstructured.SetNestedStringMap(target.Object, labels, strings.Split(labelsPath, ".")...)
	if err != nil {
		return fmt.Errorf("error attaching load balancer labels for manifest (kind: %s, name: %s, namespace: %s): %v",
			target.GetKind(),
			target.GetName(),
			target.GetNamespace(),
			err)
	}

	return nil
}

// disjoint returns true if the two specified maps have no keys in common.
func disjoint(m1, m2 map[string]string) bool {
	disjoint := true

	for k := range m1 {
		if _, ok := m2[k]; ok {
			disjoint = false
			break
		}
	}

	return disjoint
}
