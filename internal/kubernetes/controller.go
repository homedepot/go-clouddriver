package kubernetes

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/cached/disk"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

var (
	useDiskCache bool
)

//go:generate counterfeiter . Controller
// Controller holds the ability to generate a new
// dynamic kubernetes client.
type Controller interface {
	NewClient(*rest.Config) (Client, error)
	NewClientset(*rest.Config) (Clientset, error)
}

// NewController returns an instance of Controller.
func NewController() Controller {
	return &controller{}
}

type controller struct{}

// UseDiskCache sets the controller to generate clients that use
// disk cache instead of memory cache for discovery and HTTP responses.
func UseDiskCache() {
	useDiskCache = true
}

// NewClient returns a new dynamic Kubernetes client. By default it returns
// a client that uses in-memory cache store, unless `useDiskCache` is set
// to true. This is where the client stores and references its discovery of
// the Kubernetes API server.
func (c *controller) NewClient(config *rest.Config) (Client, error) {
	var (
		client Client
		err    error
	)

	if useDiskCache {
		client, err = newClientWithDefaultDiskCache(config)
	} else {
		client, err = newClientWithMemoryCache(config)
	}

	return client, err
}

// NewClientset returns a new kubernetes Clientset wrapper.
func (c *controller) NewClientset(config *rest.Config) (Clientset, error) {
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &clientset{
		clientset: cs,
	}, nil
}

const (
	// Default cache directory.
	cacheDir       = "/var/kube/cache"
	defaultTimeout = 180 * time.Second
	ttl            = 10 * time.Minute
)

var (
	mux                   sync.Mutex
	cachedMemCacheClients = map[string]memory.CachedDiscoveryClient{}
)

func newClientWithMemoryCache(config *rest.Config) (Client, error) {
	// If the timeout is not set, set it to the default timeout.
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	mc, err := memCacheClientForConfig(config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(mc)
	kubeClient := &client{
		c:      dynamicClient,
		config: config,
		mapper: mapper,
	}

	return kubeClient, nil
}

func memCacheClientForConfig(inConfig *rest.Config) (memory.CachedDiscoveryClient, error) {
	config := inConfig

	mux.Lock()
	defer mux.Unlock()

	if _, ok := cachedMemCacheClients[config.Host]; !ok {
		mc, err := memory.NewCachedDiscoveryClientForConfig(config, ttl)
		if err != nil {
			return nil, err
		}

		cachedMemCacheClients[config.Host] = mc
	}

	memCacheClient := cachedMemCacheClients[config.Host]

	return memCacheClient.CopyForConfig(config)
}

func newClientWithDefaultDiskCache(config *rest.Config) (Client, error) {
	// If the timeout is not set, set it to the default timeout.
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

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
