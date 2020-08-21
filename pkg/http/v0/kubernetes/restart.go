package kubernetes

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/deployment"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
)

// Perform a `kubectl rollout restart` by setting an annotation on a pod template
// to the current time in RFC3339.
func RollingRestartManifest(c *gin.Context, rrm RollingRestartManifestRequest) error {
	sc := sql.Instance(c)
	kc := kubernetes.Instance(c)

	provider, err := sc.GetKubernetesProvider(rrm.Account)
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

	if err = kc.SetDynamicClientForConfig(config); err != nil {
		return err
	}

	a := strings.Split(rrm.ManifestName, " ")
	kind := a[0]
	name := a[1]

	u, err := kc.Get(kind, name, rrm.Location)
	if err != nil {
		return err
	}

	switch strings.ToLower(kind) {
	case "deployment":
		d := deployment.New(u.Object)

		// add annotations to pod spec
		// kubectl.kubernetes.io/restartedAt: "2020-08-21T03:56:27Z"
		d.AnnotateTemplate("clouddriver.spinnaker.io/restartedAt",
			time.Now().In(time.UTC).Format(time.RFC3339))

		annotatedObject, err := d.ToUnstructured()
		if err != nil {
			return err
		}

		_, err = kc.Apply(&annotatedObject)
		if err != nil {
			return err
		}
	}

	return nil
}
