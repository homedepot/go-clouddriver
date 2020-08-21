package kubernetes

import (
	"fmt"
	"log"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/kubernetes/deployment"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/replicaset"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationSpinnakerArtifactLocation   = `artifact.spinnaker.io/location`
	AnnotationSpinnakerArtifactName       = `artifact.spinnaker.io/name`
	AnnotationSpinnakerArtifactType       = `artifact.spinnaker.io/type`
	AnnotationSpinnakerMonikerApplication = `moniker.spinnaker.io/application`
	AnnotationSpinnakerMonikerCluster     = `moniker.spinnaker.io/cluster`
)

func AddSpinnakerAnnotations(u *unstructured.Unstructured, application string) error {
	var err error

	name := u.GetName()
	namespace := u.GetNamespace()
	gvk := u.GroupVersionKind()

	t := fmt.Sprintf("kubernetes/%s", strings.ToLower(gvk.Kind))
	cluster := fmt.Sprintf("%s %s", strings.ToLower(gvk.Kind), name)

	// Add reserved annotations.
	// https://spinnaker.io/reference/providers/kubernetes-v2/#reserved-annotations
	annotate(u, AnnotationSpinnakerArtifactLocation, namespace)
	annotate(u, AnnotationSpinnakerArtifactName, name)
	annotate(u, AnnotationSpinnakerArtifactType, t)
	annotate(u, AnnotationSpinnakerMonikerApplication, application)
	annotate(u, AnnotationSpinnakerMonikerCluster, cluster)

	log.Println("ANNOTATING KIND:", gvk.Kind)
	if strings.EqualFold(gvk.Kind, "deployment") {
		d := deployment.New(u.Object)

		// Add spinnaker annotations to the deployment pod template.
		d.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		d.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		d.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		d.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		*u, err = d.ToUnstructured()
		if err != nil {
			return err
		}
	}

	if strings.EqualFold(gvk.Kind, "replicaset") {
		rs := replicaset.New(u.Object)

		// Add spinnaker annotations to the replicaset pod template.
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactLocation, namespace)
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactName, name)
		rs.AnnotateTemplate(AnnotationSpinnakerArtifactType, t)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerApplication, application)
		rs.AnnotateTemplate(AnnotationSpinnakerMonikerCluster, cluster)

		*u, err = rs.ToUnstructured()
		if err != nil {
			return err
		}
	}

	return nil
}

func annotate(o *unstructured.Unstructured, key, value string) {
	annotations := o.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key] = value
	o.SetAnnotations(annotations)
}