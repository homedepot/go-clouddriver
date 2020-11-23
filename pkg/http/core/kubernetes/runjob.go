package kubernetes

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
)

const randNameNumber = 5

func RunJob(c *gin.Context, rj RunJobRequest) {
	ac := arcade.Instance(c)
	kc := kube.ControllerInstance(c)
	sc := sql.Instance(c)
	taskID := clouddriver.TaskIDFromContext(c)

	provider, err := sc.GetKubernetesProvider(rj.Account)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	token, err := ac.Token()
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
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
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	u, err := kc.ToUnstructured(rj.Manifest)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	err = kc.AddSpinnakerAnnotations(u, rj.Application)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	err = kc.AddSpinnakerLabels(u, rj.Application)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	name := u.GetName()
	generateName := u.GetGenerateName()

	if name == "" && generateName != "" {
		u.SetName(generateName + rand.String(randNameNumber))
	}

	meta, err := client.Apply(u)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	// TODO don't hardcode kind.
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

	err = sc.CreateKubernetesResource(kr)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}
}
