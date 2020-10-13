package kubernetes

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewRollingRestartAction(ac ActionConfig) Action {
	return &rollingRestart{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		rr: ac.Operation.RollingRestartManifest,
	}
}

type rollingRestart struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	rr *RollingRestartManifestRequest
}

// Perform a `kubectl rollout restart` by setting an annotation on a pod template
// to the current time in RFC3339.
func (rr *rollingRestart) Run() error {
	provider, err := rr.sc.GetKubernetesProvider(rr.rr.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	token, err := rr.ac.Token()
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

	client, err := rr.kc.NewClient(config)
	if err != nil {
		return err
	}

	a := strings.Split(rr.rr.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := client.Get(kind, name, rr.rr.Location)
	if err != nil {
		return err
	}

	switch strings.ToLower(kind) {
	case "deployment":
		d := kubernetes.NewDeployment(u.Object)

		// add annotations to pod spec
		// kubectl.kubernetes.io/restartedAt: "2020-08-21T03:56:27Z"
		d.AnnotateTemplate("clouddriver.spinnaker.io/restartedAt",
			time.Now().In(time.UTC).Format(time.RFC3339))

		annotatedObject, err := d.ToUnstructured()
		if err != nil {
			return err
		}

		_, err = client.Apply(&annotatedObject)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("restarting kind %s not currently supported", kind)
	}

	return nil
}
