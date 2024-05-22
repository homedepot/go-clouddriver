/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package memory

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	openapi_v2 "github.com/google/gnostic-models/openapiv2"
	"github.com/gregjones/httpcache"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/openapi"
	restclient "k8s.io/client-go/rest"
)

type CachedDiscoveryClient interface {
	ServerResourcesForGroupVersion(string) (*metav1.APIResourceList, error)
	ServerResources() ([]*metav1.APIResourceList, error)
	ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error)
	ServerGroups() (*metav1.APIGroupList, error)
	RESTClient() restclient.Interface
	ServerPreferredResources() ([]*metav1.APIResourceList, error)
	ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error)
	ServerVersion() (*version.Info, error)
	OpenAPISchema() (*openapi_v2.Document, error)
	OpenAPIV3() openapi.Client
	WithLegacy() discovery.DiscoveryInterface
	Fresh() bool
	Invalidate()
}

// memCacheClient can Invalidate() to stay up-to-date with discovery
// information. It is modeled after the disk cache implementation.
type memCacheClient struct {
	delegate discovery.DiscoveryInterface
	// ourEntries holds entries created during this process.
	ourEntries  map[string]struct{}
	invalidated bool
	fresh       bool

	c *Cache
}

// entry represents an in-memory cache of an API discovery resource.
type entry struct {
	Content   interface{}
	CreatedAt time.Time
}

var _ discovery.CachedDiscoveryInterface = &memCacheClient{}

// ServerResourcesForGroupVersion returns the supported resources for a group and version.
func (m *memCacheClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	cachedEntry, err := m.getCachedEntry(groupVersion)
	if err == nil && cachedEntry.Content != nil {
		b, err := json.Marshal(cachedEntry.Content)
		if err != nil {
			return nil, err
		}

		cachedResources := &metav1.APIResourceList{}

		err = json.Unmarshal(b, cachedResources)
		if err != nil {
			return nil, err
		}

		return cachedResources, nil
	}

	liveResources, err := m.delegate.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return liveResources, err
	}

	m.createCachedEntry(groupVersion, liveResources)

	return liveResources, nil
}

// ServerResources returns the supported resources for all groups and versions.
// Deprecated: use ServerGroupsAndResources instead.
func (m *memCacheClient) ServerResources() ([]*metav1.APIResourceList, error) {
	_, rs, err := discovery.ServerGroupsAndResources(m)

	return rs, err
}

// ServerGroupsAndResources returns the groups and supported resources for all groups and versions.
func (m *memCacheClient) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return discovery.ServerGroupsAndResources(m)
}

// ServerGroups returns the supported groups, with information like supported versions and the
// preferred version.
func (m *memCacheClient) ServerGroups() (*metav1.APIGroupList, error) {
	cachedEntry, err := m.getCachedEntry("servergroups")
	if err == nil && cachedEntry.Content != nil {
		b, err := json.Marshal(cachedEntry.Content)
		if err != nil {
			return nil, err
		}

		cachedResources := &metav1.APIGroupList{}

		err = json.Unmarshal(b, cachedResources)
		if err != nil {
			return nil, err
		}

		return cachedResources, nil
	}

	liveGroups, err := m.delegate.ServerGroups()
	if err != nil {
		return nil, err
	}

	if liveGroups == nil || len(liveGroups.Groups) == 0 {
		return liveGroups, err
	}

	m.createCachedEntry("servergroups", liveGroups)

	return liveGroups, nil
}

func (m *memCacheClient) getCachedEntry(key string) (entry, error) {
	m.c.mutex.Lock()
	defer m.c.mutex.Unlock()

	_, ourEntry := m.ourEntries[key]
	if m.invalidated && !ourEntry {
		return entry{}, errors.New("cache invalidated")
	}

	cachedEntry, exists := m.c.entries[key]
	if !exists {
		return entry{}, errors.New("cache entry does not exist")
	}

	if time.Now().After(cachedEntry.CreatedAt.Add(m.c.ttl)) {
		return entry{}, errors.New("cache expired")
	}

	m.fresh = m.fresh && ourEntry

	return cachedEntry, nil
}

func (m *memCacheClient) createCachedEntry(key string, content interface{}) {
	m.c.mutex.Lock()
	defer m.c.mutex.Unlock()

	m.c.entries[key] = newEntry(content)
	m.ourEntries[key] = struct{}{}
}

func (m *memCacheClient) CachedDiscoveryInterface(key string, content interface{}) {
	m.c.mutex.Lock()
	defer m.c.mutex.Unlock()

	m.c.entries[key] = newEntry(content)
	m.ourEntries[key] = struct{}{}
}

// newEntry creates a cached entry and generates its created at timestamp.
func newEntry(content interface{}) entry {
	return entry{
		Content:   content,
		CreatedAt: time.Now(),
	}
}

// RESTClient returns a RESTClient that is used to communicate with API server
// by this client implementation.
func (m *memCacheClient) RESTClient() restclient.Interface {
	return m.delegate.RESTClient()
}

// ServerPreferredResources returns the supported resources with the version preferred by the
// server.
func (m *memCacheClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerPreferredResources(m)
}

// ServerPreferredNamespacedResources returns the supported namespaced resources with the
// version preferred by the server.
func (m *memCacheClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerPreferredNamespacedResources(m)
}

// ServerVersion retrieves and parses the server's version (git version).
func (m *memCacheClient) ServerVersion() (*version.Info, error) {
	return m.delegate.ServerVersion()
}

// OpenAPISchema retrieves and parses the swagger API schema the server supports.
func (m *memCacheClient) OpenAPISchema() (*openapi_v2.Document, error) {
	return m.delegate.OpenAPISchema()
}

// OpenAPIV3 implements discovery.CachedDiscoveryInterface.
func (m *memCacheClient) OpenAPIV3() openapi.Client {
	return m.delegate.OpenAPIV3()
}

// WithLegacy implements discovery.CachedDiscoveryInterface.
func (m *memCacheClient) WithLegacy() discovery.DiscoveryInterface {
	return m.delegate.WithLegacy()
}

// Fresh is supposed to tell the caller whether or not to retry if the cache
// fails to find something (false = retry, true = no need to retry).
func (m *memCacheClient) Fresh() bool {
	m.c.mutex.Lock()
	defer m.c.mutex.Unlock()

	return m.fresh
}

// Invalidate enforces that no cached data that is older than the current time
// is used.
func (m *memCacheClient) Invalidate() {
	m.c.mutex.Lock()
	defer m.c.mutex.Unlock()

	m.ourEntries = map[string]struct{}{}
	m.fresh = true
	m.invalidated = true
}

// Cache is an in-memory store for API discovery objects.
type Cache struct {
	// ttl is how long the cache should be considered valid
	ttl time.Duration

	mutex *sync.Mutex

	// entries is a respresentation of everything that has been requested from the cache.
	// Think of it like files on a filesystem - each entry holds content and a created ts
	// of an API discovery resource. Like the disk cache, it should not be emptied.
	entries map[string]entry

	// Used to cache API discovery responses.
	httpMemCache *httpcache.MemoryCache
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl:          ttl,
		entries:      map[string]entry{},
		httpMemCache: httpcache.NewMemoryCache(),
		mutex:        &sync.Mutex{},
	}
}

// NewClientForConfig returns a cached discovery client backed by the calling
// in-memory cache store.
func (c *Cache) NewClientForConfig(config *restclient.Config) (CachedDiscoveryClient, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// update the given restconfig with a custom roundtripper that
	// understands how to handle cache responses.
	config = restclient.CopyConfig(config)
	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return newMemCacheRoundTripper(rt, c.httpMemCache)
	})

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return &memCacheClient{
		delegate:    discoveryClient,
		c:           c,
		fresh:       true,
		invalidated: false,
		ourEntries:  map[string]struct{}{},
	}, nil
}
