package kubernetes

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewRollingRestartAction(ac ActionConfig) Action {
	return &rollingRestart{
		sc: ac.SQLClient,
		kc: ac.KubeClient,
		id: ac.ID,
		rr: ac.Operation.RollingRestartManifest,
	}
}

type rollingRestart struct {
	sc sql.Client
	kc kubernetes.Client
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

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: provider.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	if err = rr.kc.SetDynamicClientForConfig(config); err != nil {
		return err
	}

	a := strings.Split(rr.rr.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := rr.kc.Get(kind, name, rr.rr.Location)
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

		_, err = rr.kc.Apply(&annotatedObject)
		if err != nil {
			return err
		}
	}

	return nil
}
