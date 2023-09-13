package kubernetes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	kube "github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/util/rand"
)

const randNameNumber = 5

func (cc *Controller) RunJob(c *gin.Context, rj RunJobRequest) {
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := cc.KubernetesProvider(rj.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	u, err := kube.ToUnstructured(rj.Manifest)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	namespace := ""

	// Preserve backwards compatibility
	if len(provider.Namespaces) == 1 {
		namespace = provider.Namespaces[0]
	}

	kubernetes.SetNamespaceOnManifest(&u, namespace)

	err = provider.ValidateNamespaceAccess(u.GetNamespace()) // pass in the current manifest's namespace
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	err = kube.AddSpinnakerAnnotations(&u, rj.Application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	err = kube.AddSpinnakerLabels(&u, rj.Application)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	name := u.GetName()
	generateName := u.GetGenerateName()

	if name == "" && generateName != "" {
		u.SetName(generateName + rand.String(randNameNumber))
	}

	kubernetes.BindArtifacts(&u, append(rj.RequiredArtifacts, rj.OptionalArtifacts...), rj.Account)

	meta := kubernetes.Metadata{}
	if kubernetes.Replace(u) {
		meta, err = provider.Client.Replace(&u)
	} else {
		meta, err = provider.Client.Apply(&u)
	}

	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	kr := kubernetes.Resource{
		AccountName:  rj.Account,
		ID:           uuid.New().String(),
		TaskID:       taskID,
		APIGroup:     meta.Group,
		Name:         meta.Name,
		Namespace:    meta.Namespace,
		Resource:     meta.Resource,
		Version:      meta.Version,
		Kind:         "job",
		SpinnakerApp: rj.Application,
	}

	err = cc.SQLClient.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}
}
