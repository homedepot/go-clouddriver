package internal

import (
	"github.com/homedepot/go-clouddriver/internal/arcade"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/fiat"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/sql"
)

// Controller holds all non request-scoped objects.
type Controller struct {
	ArcadeClient                  arcade.Client
	ArtifactCredentialsController artifact.CredentialsController
	FiatClient                    fiat.Client
	KubernetesController          kubernetes.Controller
	SQLClient                     sql.Client
}
