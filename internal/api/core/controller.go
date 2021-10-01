package core

import (
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
)

// Controller holds all non request-scoped objects.
type Controller struct {
	*internal.Controller
	KC *kubernetes.Controller
}
