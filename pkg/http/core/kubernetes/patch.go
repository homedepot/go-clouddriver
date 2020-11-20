package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func (ah *actionHandler) NewPatchManifestAction(ac ActionConfig) Action {
	return &patchManfest{
		ac: ac.ArcadeClient,
		sc: ac.SQLClient,
		kc: ac.KubeController,
		id: ac.ID,
		pm: ac.Operation.PatchManifest,
	}
}

type patchManfest struct {
	ac arcade.Client
	sc sql.Client
	kc kubernetes.Controller
	id string
	pm *PatchManifestRequest
}

func (p *patchManfest) Run() error {
	provider, err := p.sc.GetKubernetesProvider(p.pm.Account)
	if err != nil {
		return err
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return err
	}

	token, err := p.ac.Token()
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

	client, err := p.kc.NewClient(config)
	if err != nil {
		return err
	}

	b, err := json.Marshal(p.pm.PatchBody)
	if err != nil {
		return err
	}

	// Manifest name is *really* the Spinnaker cluster - i.e. "deployment test-deployment", so we
	// need to split on a whitespace and get the actual name of the manifest.
	kind := ""
	name := p.pm.ManifestName

	a := strings.Split(p.pm.ManifestName, " ")
	if len(a) > 1 {
		kind = a[0]
		name = a[1]
	}

	// Merge strategy can be "strategic", "json", or "merge".
	var strategy types.PatchType

	switch p.pm.Options.MergeStrategy {
	case "strategic":
		strategy = types.StrategicMergePatchType
	case "json":
		strategy = types.JSONPatchType
	case "merge":
		strategy = types.MergePatchType
	default:
		return fmt.Errorf("invalid merge strategy %s", p.pm.Options.MergeStrategy)
	}

	meta, _, err := client.PatchUsingStrategy(kind, name, p.pm.Location, b, strategy)
	if err != nil {
		return err
	}

	// TODO Record the applied patch in the kubernetes.io/change-cause annotation. If the annotation already exists, the contents are replaced.
	if p.pm.Options.Record {
	}

	kr := kubernetes.Resource{
		AccountName:  p.pm.Account,
		ID:           uuid.New().String(),
		TaskID:       p.id,
		APIGroup:     meta.Group,
		Name:         meta.Name,
		Namespace:    meta.Namespace,
		Resource:     meta.Resource,
		Version:      meta.Version,
		Kind:         meta.Kind,
		SpinnakerApp: p.pm.App,
	}

	err = p.sc.CreateKubernetesResource(kr)
	if err != nil {
		return err
	}

	return nil
}
