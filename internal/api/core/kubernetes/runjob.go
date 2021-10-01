package kubernetes

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	kube "github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
)

const randNameNumber = 5

func (cc *Controller) RunJob(c *gin.Context, rj RunJobRequest) {
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := cc.SQLClient.GetKubernetesProvider(rj.Account)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	token, err := cc.ArcadeClient.Token(provider.TokenProvider)
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

	client, err := cc.KubernetesController.NewClient(config)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	u, err := kube.ToUnstructured(rj.Manifest)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
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

	namespace := u.GetNamespace()
	if provider.Namespace != nil {
		namespace = *provider.Namespace
	}

	meta, err := client.ApplyWithNamespaceOverride(&u, namespace)
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
