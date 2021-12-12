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

type MemCachedDiscoveryClient interface {
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
}

// memCacheClient can Invalidate() to stay up-to-date with discovery
// information.
type memCacheClient struct {
	delegate discovery.DiscoveryInterface

	// ttl is how long the cache should be considered valid
	ttl time.Duration

	mutex           sync.Mutex
	ourEntries      map[string]*metav1.APIResourceList
	ourServerGroups *metav1.APIGroupList
	ourTTLs         map[string]time.Time
	// invalidated is true if all cache files should be ignored that are not ours (e.g. after Invalidate() was called)
	invalidated bool
	// fresh is true if all used cache files were ours
	fresh bool
}

var _ discovery.CachedDiscoveryInterface = &memCacheClient{}

// ServerResourcesForGroupVersion returns the supported resources for a group and version.
func (d *memCacheClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cachedEntry, err := d.getCachedEntry(groupVersion)
	if err == nil && cachedEntry != nil {
		return cachedEntry, nil
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

func (d *memCacheClient) ServerGroups() (*metav1.APIGroupList, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cachedServerGroups, err := d.getCachedServerGroups()
	// Don't fail on errors, we either don't have an entry or won't be able to run the cached check.
	// Either way we can fallback.
	if err == nil && cachedServerGroups != nil {
		return cachedServerGroups, nil
	}

	liveGroups, err := d.delegate.ServerGroups()
	if err != nil {
		return nil, err
	}

	if liveGroups == nil || len(liveGroups.Groups) == 0 {
		return liveGroups, err
	}

	d.createCachedServerGroups(liveGroups)

	return liveGroups, nil
}

func (d *memCacheClient) getCachedServerGroups() (*metav1.APIGroupList, error) {
	if d.invalidated {
		return nil, errors.New("cache invalidated")
	}

	if d.ourServerGroups == nil {
		return nil, errors.New("server groups not defined")
	}

	t, ok := d.ourTTLs["servergroups"]
	if ok && time.Now().After(t) {
		return nil, errors.New("cache expired")
	}

	return d.ourServerGroups, nil
}

func (d *memCacheClient) createCachedServerGroups(serverGroups *metav1.APIGroupList) {
	d.ourServerGroups = serverGroups
	d.ourTTLs["servergroups"] = time.Now().Add(d.ttl)
}

func (d *memCacheClient) getCachedEntry(key string) (*metav1.APIResourceList, error) {
	cachedEntry, ourEntry := d.ourEntries[key]
	if d.invalidated && !ourEntry {
		return nil, errors.New("cache invalidated")
	}

	if len(d.ourEntries) == 0 && !ourEntry {
		return nil, errors.New("entry not found")
	}

	t, ok := d.ourTTLs[key]
	if ok && time.Now().After(t) {
		return nil, errors.New("cache expired")
	}

	d.fresh = d.fresh && ourEntry

	return cachedEntry, nil
}

func (d *memCacheClient) createCachedEntry(key string, entry *metav1.APIResourceList) {
	d.ourEntries[key] = entry
	d.ourTTLs[key] = time.Now().Add(d.ttl)
}

func (d *memCacheClient) RESTClient() restclient.Interface {
	return d.delegate.RESTClient()
}

func (d *memCacheClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerPreferredResources(d)
}

func (d *memCacheClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerPreferredNamespacedResources(d)
}

func (d *memCacheClient) ServerVersion() (*version.Info, error) {
	return d.delegate.ServerVersion()
}

func (d *memCacheClient) OpenAPISchema() (*openapi_v2.Document, error) {
	return d.delegate.OpenAPISchema()
}

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

	d.ourEntries = map[string]*metav1.APIResourceList{}
	d.ourTTLs = map[string]time.Time{}
	d.ourServerGroups = nil
	d.fresh = true
	d.invalidated = true
}

func NewMemCachedDiscoveryClientForConfig(config *restclient.Config, ttl time.Duration) (MemCachedDiscoveryClient, error) {
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

	return newMemCachedDiscoveryClient(discoveryClient, ttl), nil
}

func NewMemCachedDiscoveryClient(delegate discovery.DiscoveryInterface, ttl time.Duration) MemCachedDiscoveryClient {
	return newMemCachedDiscoveryClient(delegate, ttl)
}

// NewMemCacheClient creates a new CachedDiscoveryInterface which caches
// discovery information in memory and will stay up-to-date if Invalidate is
// called with regularity.
func newMemCachedDiscoveryClient(delegate discovery.DiscoveryInterface, ttl time.Duration) MemCachedDiscoveryClient {
	return &memCacheClient{
		delegate:   delegate,
		ttl:        ttl,
		ourEntries: map[string]*metav1.APIResourceList{},
		ourTTLs:    map[string]time.Time{},
		fresh:      true,
	}
}
