package kubernetes

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes/cached/disk"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

const (
	ControllerInstanceKey = `KubeController`
)

//go:generate counterfeiter . Controller
type Controller interface {
	NewClient(*rest.Config) (Client, error)
	ToUnstructured(map[string]interface{}) (*unstructured.Unstructured, error)
	AddSpinnakerAnnotations(u *unstructured.Unstructured, application string) error
	AddSpinnakerLabels(u *unstructured.Unstructured, application string) error
}

func NewController() Controller {
	return &controller{}
}

type controller struct{}

func (c *controller) NewClient(config *rest.Config) (Client, error) {
	return newClientWithDefaultDiskCache(config)
}

const (
	// Default cache directory.
	cacheDir       = "/var/kube/cache"
	defaultTimeout = 180 * time.Second
)

var (
	ttl = 10 * time.Minute
)

func newClientWithDefaultDiskCache(config *rest.Config) (Client, error) {
	config.Timeout = defaultTimeout

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Some code to define this take from
	// https://github.com/kubernetes/cli-runtime/blob/master/pkg/genericclioptions/config_flags.go#L215
	httpCacheDir := filepath.Join(cacheDir, "http")
	discoveryCacheDir := computeDiscoverCacheDir(filepath.Join(cacheDir, "discovery"), config.Host)

	// DiscoveryClient queries API server about the resources
	cdc, err := disk.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, ttl)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cdc)

	kubeClient := &client{
		c:      dynamicClient,
		config: config,
		mapper: mapper,
	}

	return kubeClient, nil
}

// overlyCautiousIllegalFileCharacters matches characters that *might* not be supported.
// Windows is really restrictive, so this is really restrictive.
var overlyCautiousIllegalFileCharacters = regexp.MustCompile(`[^(\w/\.)]`)

// computeDiscoverCacheDir takes the parentDir and the host and comes up with a "usually non-colliding" name.
func computeDiscoverCacheDir(parentDir, host string) string {
	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.
	// Even if we do collide the problem is short lived
	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
	return filepath.Join(parentDir, safeHost)
}

func ControllerInstance(c *gin.Context) Controller {
	return c.MustGet(ControllerInstanceKey).(Controller)
}
