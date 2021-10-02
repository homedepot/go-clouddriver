package core

import (
	"github.com/homedepot/go-clouddriver/internal"
)

// Controller holds all non request-scoped objects.
type Controller struct {
	*internal.Controller
}
