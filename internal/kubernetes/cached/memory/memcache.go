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
	"fmt"
	"net/http"
	"sync"
	"time"

	openapi_v2 "github.com/googleapis/gnostic/openapiv2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
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

type cacheEntry struct {
	entry interface{}
	err   error
}

// memCacheClient can Invalidate() to stay up-to-date with discovery
// information.
type memCacheClient struct {
	delegate discovery.DiscoveryInterface

	// ttl is how long the cache should be considered valid
	ttl time.Duration

	mutex       sync.Mutex
	entries     map[string]*cacheEntry
	expirations map[string]time.Time
	invalidated bool
	// valid is true if the cache is not populated or until a requested groupVersion does not exist in entries.
	valid bool
}

var _ discovery.CachedDiscoveryInterface = &memCacheClient{}

// ServerResourcesForGroupVersion returns the supported resources for a group and version.
func (d *memCacheClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.valid {
		if err := d.refreshLocked(); err != nil {
			return nil, err
		}
	}

	cachedEntry := &metav1.APIResourceList{}

	err := d.getCachedEntry(groupVersion, cachedEntry)
	if err == nil && d.valid {
		return cachedEntry, nil
	}

	if err := d.refreshLocked(); err != nil {
		return nil, err
	}

	err = d.getCachedEntry(groupVersion, cachedEntry)
	if err != nil {
		return nil, err
	}

	return cachedEntry, nil
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

	if !d.valid {
		if err := d.refreshLocked(); err != nil {
			return nil, err
		}
	}

	cachedEntry := &metav1.APIGroupList{}
	// Don't fail on errors, we either don't have an entry or won't be able to run the cached check.
	// Either way we can fallback.
	err := d.getCachedEntry("servergroups", cachedEntry)
	if err == nil && d.valid {
		return cachedEntry, nil
	}

	if err := d.refreshLocked(); err != nil {
		return nil, err
	}

	err = d.getCachedEntry("servergroups", cachedEntry)
	if err != nil {
		return nil, err
	}

	return cachedEntry, nil
}

func (d *memCacheClient) getCachedEntry(key string, into interface{}) error {
	cachedEntry, exists := d.entries[key]
	if d.invalidated && !exists {
		return errors.New("cache invalidated")
	}

	t, ok := d.expirations[key]
	if ok && time.Now().After(t) {
		return errors.New("cache expired")
	}

	if exists {
		b, err := json.Marshal(cachedEntry.entry)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, into)
		if err != nil {
			return err
		}
	}

	d.valid = d.valid && exists

	return nil
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

	return d.valid
}

// refreshLocked refreshes the state of cache. The caller must hold d.lock for
// writing.
func (d *memCacheClient) refreshLocked() error {
	// TODO: Could this multiplicative set of calls be replaced by a single call
	// to ServerResources? If it's possible for more than one resulting
	// APIResourceList to have the same GroupVersion, the lists would need merged.
	gl, err := d.delegate.ServerGroups()
	if err != nil || len(gl.Groups) == 0 {
		utilruntime.HandleError(fmt.Errorf("couldn't get current server API group list: %v", err))

		return err
	}

	wg := &sync.WaitGroup{}
	resultLock := &sync.Mutex{}
	rl := map[string]*cacheEntry{}

	for _, g := range gl.Groups {
		for _, v := range g.Versions {
			gv := v.GroupVersion

			wg.Add(1)

			go func() {
				defer wg.Done()
				defer utilruntime.HandleCrash()

				r, err := d.serverResourcesForGroupVersion(gv)
				if err != nil {
					utilruntime.HandleError(fmt.Errorf("couldn't get resource list for %v: %v", gv, err))
				}

				resultLock.Lock()
				defer resultLock.Unlock()

				rl[gv] = &cacheEntry{r, err}
			}()
		}
	}

	wg.Wait()

	d.entries = rl
	d.entries["servergroups"] = &cacheEntry{gl, nil}
	e := time.Now().Add(d.ttl)

	for k := range d.entries {
		d.expirations[k] = e
	}

	d.valid = true

	return nil
}

func (d *memCacheClient) serverResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	r, err := d.delegate.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return r, err
	}

	if len(r.APIResources) == 0 {
		return r, fmt.Errorf("Got empty response for: %v", groupVersion)
	}

	return r, nil
}

// Invalidate enforces that no cached data that is older than the current time
// is used.
func (d *memCacheClient) Invalidate() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.entries = map[string]*cacheEntry{}
	d.expirations = map[string]time.Time{}
	d.valid = false
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
		delegate:    delegate,
		ttl:         ttl,
		entries:     map[string]*cacheEntry{},
		expirations: map[string]time.Time{},
		// valid:       true,
	}
}
