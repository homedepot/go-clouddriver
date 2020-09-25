package kubernetes

import (
	"encoding/base64"

	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
)

const randNameNumber = 5

func (ah *actionHandler) NewRunJobAction(ac ActionConfig) Action {
	return &runJob{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		rj: ac.Operation.RunJob,
	}
}

type runJob struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	rj *RunJobRequest
}

func (r *runJob) Run() error {
	provider, err := r.sc.GetKubernetesProvider(r.rj.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	token, err := r.ac.Token()
	if err != nil {
		return err
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := r.kc.NewClient(config)
	if err != nil {
		return err
	}

	u, err := r.kc.ToUnstructured(r.rj.Manifest)
	if err != nil {
		return err
	}

	err = r.kc.AddSpinnakerAnnotations(u, r.rj.Application)
	if err != nil {
		return err
	}

	err = r.kc.AddSpinnakerLabels(u, r.rj.Application)
	if err != nil {
		return err
	}

	name := u.GetName()
	generateName := u.GetGenerateName()

	if name == "" && generateName != "" {
		u.SetName(generateName + rand.String(randNameNumber))
	}

	meta, err := client.Apply(u)
	if err != nil {
		return err
	}

	// TODO don't hardcode kind.
	kr := kubernetes.Resource{
		AccountName:  r.rj.Account,
		ID:           uuid.New().String(),
		TaskID:       r.id,
		APIGroup:     meta.Group,
		Name:         meta.Name,
		Namespace:    meta.Namespace,
		Resource:     meta.Resource,
		Version:      meta.Version,
		Kind:         "job",
		SpinnakerApp: r.rj.Application,
	}

	err = r.sc.CreateKubernetesResource(kr)
	if err != nil {
		return err
	}

	return nil
}
