package internal

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/homedepot/go-clouddriver/internal/arcade"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/fiat"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/sql"
	"k8s.io/client-go/rest"
)

// Controller holds all non request-scoped objects.
type Controller struct {
	ArcadeClient                  arcade.Client
	ArtifactCredentialsController artifact.CredentialsController
	FiatClient                    fiat.Client
	KubernetesController          kubernetes.Controller
	SQLClient                     sql.Client
}

// KubernetesProvider returns a kubernetes provider instance
// for a given account name. It instantiates a kubernetes
// Client and Clientset and attaches these instances to the provider.
func (cc *Controller) KubernetesProvider(account string) (*kubernetes.Provider, error) {
	return cc.KubernetesProviderWithTimeout(account, 0)
}

// KubernetesProviderWithTimeout returns a kubernetes provider,
// defining its client and clientset timeouts to be the timeout
// passed in. If no timeout is passed this field is not set.
func (cc *Controller) KubernetesProviderWithTimeout(account string,
	timeout time.Duration) (*kubernetes.Provider, error) {
	// Get the provider info for the account.
	provider, err := cc.SQLClient.GetKubernetesProvider(account)
	if err != nil {
		return nil, fmt.Errorf("internal: error getting kubernetes provider %s: %v", account, err)
	}

	// Decode the provider's CA data.
	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		return nil, fmt.Errorf("internal: error decoding provider CA data: %v", err)
	}

	// Grab the auth token from arcade.
	token, err := cc.ArcadeClient.Token(provider.TokenProvider)
	if err != nil {
		return nil, fmt.Errorf("internal: error getting token from arcade for provider %s: %v",
			provider.TokenProvider, err)
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	if timeout > 0 {
		config.Timeout = timeout
	}

	client, err := cc.KubernetesController.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("internal: error creating new kubernetes client: %v", err)
	}

	clientset, err := cc.KubernetesController.NewClientset(config)
	if err != nil {
		return nil, fmt.Errorf("internal: error creating new kubernetes clientset: %v", err)
	}

	provider.WithClient(client)
	provider.WithClientset(clientset)

	return &provider, nil
}
