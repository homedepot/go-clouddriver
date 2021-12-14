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

	openapi_v2 "github.com/googleapis/gnostic/openapiv2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
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
	Fresh() bool
	Invalidate()
	Reset()
}

// memCacheClient can Invalidate() to stay up-to-date with discovery
// information. It is modeled after the disk cache implementation.
type memCacheClient struct {
	delegate discovery.DiscoveryInterface

	// ttl is how long the cache should be considered valid
	ttl time.Duration

	mutex sync.Mutex

	// entries is a respresentation of everything that has been requested from the cache.
	// Think of it like files on a filesystem - it will never be emptied unless the pod
	// is restarted.
	entries map[string]interface{}

	// expirations holds the expiration time (creation time plus ttl) of given entries.
	// Think of it like a cached files mod-time on disk plus ttl.
	expirations map[string]time.Time

	// ourEntries holds entries created during this process. The caller should call Reset()
	// to empty these entries if they are caching these clients.
	ourEntries  map[string]struct{}
	invalidated bool
	fresh       bool
}

var _ discovery.CachedDiscoveryInterface = &memCacheClient{}

// ServerResourcesForGroupVersion returns the supported resources for a group and version.
func (d *memCacheClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cachedEntry, err := d.getCachedEntry(groupVersion)
	if err == nil && cachedEntry != nil {
		b, err := json.Marshal(cachedEntry)
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

	liveResources, err := d.delegate.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return liveResources, err
	}

	d.createCachedEntry(groupVersion, liveResources)

	return liveResources, nil
}

// ServerResources returns the supported resources for all groups and versions.
// Deprecated: use ServerGroupsAndResources instead.
func (d *memCacheClient) ServerResources() ([]*metav1.APIResourceList, error) {
	_, rs, err := discovery.ServerGroupsAndResources(d)

	return rs, err
}

// ServerGroupsAndResources returns the groups and supported resources for all groups and versions.
func (d *memCacheClient) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return discovery.ServerGroupsAndResources(d)
}

// ServerGroups returns the supported groups, with information like supported versions and the
// preferred version.
func (d *memCacheClient) ServerGroups() (*metav1.APIGroupList, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cachedEntry, err := d.getCachedEntry("servergroups")
	if err == nil && cachedEntry != nil {
		b, err := json.Marshal(cachedEntry)
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

	liveGroups, err := d.delegate.ServerGroups()
	if err != nil {
		return nil, err
	}

	if liveGroups == nil || len(liveGroups.Groups) == 0 {
		return liveGroups, err
	}

	d.createCachedEntry("servergroups", liveGroups)

	return liveGroups, nil
}

func (d *memCacheClient) getCachedEntry(key string) (interface{}, error) {
	_, ourEntry := d.ourEntries[key]
	if d.invalidated && !ourEntry {
		return nil, errors.New("cache invalidated")
	}

	cachedEntry, exists := d.entries[key]
	if !exists {
		return nil, errors.New("cache entry does not exist")
	}

	t, ok := d.expirations[key]
	if ok && time.Now().After(t) {
		return nil, errors.New("cache expired")
	}

	d.fresh = d.fresh && ourEntry

	return cachedEntry, nil
}

func (d *memCacheClient) createCachedEntry(key string, entry interface{}) {
	d.entries[key] = entry
	d.ourEntries[key] = struct{}{}
	d.expirations[key] = time.Now().Add(d.ttl)
}

// RESTClient returns a RESTClient that is used to communicate with API server
// by this client implementation.
func (d *memCacheClient) RESTClient() restclient.Interface {
	return d.delegate.RESTClient()
}

// ServerPreferredResources returns the supported resources with the version preferred by the
// server.
func (d *memCacheClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerPreferredResources(d)
}

// ServerPreferredNamespacedResources returns the supported namespaced resources with the
// version preferred by the server.
func (d *memCacheClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerPreferredNamespacedResources(d)
}

// ServerVersion retrieves and parses the server's version (git version).
func (d *memCacheClient) ServerVersion() (*version.Info, error) {
	return d.delegate.ServerVersion()
}

// OpenAPISchema retrieves and parses the swagger API schema the server supports.
func (d *memCacheClient) OpenAPISchema() (*openapi_v2.Document, error) {
	return d.delegate.OpenAPISchema()
}

// Fresh is supposed to tell the caller whether or not to retry if the cache
// fails to find something (false = retry, true = no need to retry).
func (d *memCacheClient) Fresh() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.fresh
}

// Invalidate enforces that no cached data that is older than the current time
// is used.
func (d *memCacheClient) Invalidate() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.ourEntries = map[string]struct{}{}
	d.fresh = true
	d.invalidated = true
}

// Reset sets the memcache to its original state without invalidating it.
func (d *memCacheClient) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.fresh = true
	d.invalidated = false
	d.ourEntries = map[string]struct{}{}
}

// NewCachedDiscoveryClientForConfig creates a new DiscoveryClient for the given config, and wraps
// the created client in a CachedDiscoveryClient. The provided configuration is updated with a
// custom transport that understands cache responses.
func NewCachedDiscoveryClientForConfig(config *restclient.Config, ttl time.Duration) (CachedDiscoveryClient, error) {
	// update the given restconfig with a custom roundtripper that
	// understands how to handle cache responses.
	config = restclient.CopyConfig(config)
	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return newMemCacheRoundTripper(rt)
	})

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return newCachedDiscoveryClient(discoveryClient, ttl), nil
}

// NewCachedDiscoveryClient creates a new CachedDiscoveryClient which caches
// discovery information in memory and will stay up-to-date if Invalidate is
// called with regularity.
func NewCachedDiscoveryClient(delegate discovery.DiscoveryInterface, ttl time.Duration) CachedDiscoveryClient {
	return newCachedDiscoveryClient(delegate, ttl)
}

func newCachedDiscoveryClient(delegate discovery.DiscoveryInterface, ttl time.Duration) CachedDiscoveryClient {
	return &memCacheClient{
		delegate:    delegate,
		ttl:         ttl,
		entries:     map[string]interface{}{},
		expirations: map[string]time.Time{},
		fresh:       true,
		ourEntries:  map[string]struct{}{},
	}
}
