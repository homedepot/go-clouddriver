package kubernetes

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/cached/disk"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
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

// NewClient returns a new dynamic Kubernetes client with a default
// disk cache directory of /var/kube/cache. This is where the client
// stores and references its discovery of the Kubernetes API server.
func (c *controller) NewClient(config *rest.Config) (Client, error) {
	return newClientWithMemoryCache(config)
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

func newClientWithMemoryCache(config *rest.Config) (Client, error) {
	// If the timeout is not set, set it to the default timeout.
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	mapper, err := mapperForConfig(config)
	if err != nil {
		return nil, err
	}

	kubeClient := &client{
		c:      dynamicClient,
		config: config,
		mapper: mapper,
	}

	return kubeClient, nil
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

var (
	mux           sync.Mutex
	cachedConfigs = map[string]*rest.Config{}
	cachedMappers = map[string]*restmapper.DeferredDiscoveryRESTMapper{}
)

func mapperForConfig(inConfig *rest.Config) (*restmapper.DeferredDiscoveryRESTMapper, error) {
	config := inConfig

	if _, ok := cachedConfigs[config.Host]; ok {
		cachedConfig := cachedConfigs[config.Host]
		if string(cachedConfig.TLSClientConfig.CAData) != string(config.TLSClientConfig.CAData) ||
			cachedConfig.BearerToken != config.BearerToken {
			err := setCaches(config)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err := setCaches(config)
		if err != nil {
			return nil, err
		}
	}

	return cachedMapper(config), nil
}

func setCaches(config *rest.Config) error {
	m, err := newMapperForConfig(config)
	if err != nil {
		return err
	}

	mux.Lock()
	defer mux.Unlock()

	cachedConfigs[config.Host] = config
	cachedMappers[config.Host] = m

	return nil
}

func newMapperForConfig(config *rest.Config) (*restmapper.DeferredDiscoveryRESTMapper, error) {
	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper, nil
}

func cachedMapper(config *rest.Config) *restmapper.DeferredDiscoveryRESTMapper {
	mux.Lock()
	defer mux.Unlock()

	return cachedMappers[config.Host]
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
