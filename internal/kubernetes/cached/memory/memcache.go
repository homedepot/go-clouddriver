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
	"fmt"
	"sync"
	"syscall"
	"time"

	openapi_v2 "github.com/googleapis/gnostic/openapiv2"

	errorsutil "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	restclient "k8s.io/client-go/rest"
)

type cacheEntry struct {
	resourceList *metav1.APIResourceList
	err          error
}

// memCacheClient can Invalidate() to stay up-to-date with discovery
// information.
//
// TODO: Switch to a watch interface. Right now it will poll after each
// Invalidate() call.
type memCacheClient struct {
	delegate discovery.DiscoveryInterface

	// ttl is how long the cache should be considered valid
	ttl time.Duration

	// lock                   sync.RWMutex
	mutex sync.Mutex
	// groupToServerResources map[string]*cacheEntry
	ourEntries map[string]*cacheEntry
	ourTTLs    map[string]time.Duration
	// groupList              *metav1.APIGroupList
	// cacheValid             bool
	// invalidated is true if all cache files should be ignored that are not ours (e.g. after Invalidate() was called)
	invalidated bool
	// fresh is true if all used cache files were ours
	fresh bool
}

// Error Constants
var (
	ErrCacheNotFound = errors.New("not found")
)

var _ discovery.CachedDiscoveryInterface = &memCacheClient{}

// isTransientConnectionError checks whether given error is "Connection refused" or
// "Connection reset" error which usually means that apiserver is temporarily
// unavailable.
func isTransientConnectionError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno == syscall.ECONNREFUSED || errno == syscall.ECONNRESET
	}
	return false
}

func isTransientError(err error) bool {
	if isTransientConnectionError(err) {
		return true
	}

	if t, ok := err.(errorsutil.APIStatus); ok && t.Status().Code >= 500 {
		return true
	}

	return errorsutil.IsTooManyRequests(err)
}

// ServerResourcesForGroupVersion returns the supported resources for a group and version.
func (d *memCacheClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	// if !d.cacheValid {
	// 	if err := d.refreshLocked(); err != nil {
	// 		return nil, err
	// 	}
	// }

	cachedVal, ok := d.groupToServerResources[groupVersion]
	if ok {
		if cachedVal.err != nil && isTransientError(cachedVal.err) {
			r, err := d.serverResourcesForGroupVersion(groupVersion)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf("couldn't get resource list for %v: %v", groupVersion, err))
			}

			cachedVal = &cacheEntry{r, err}
			d.groupToServerResources[groupVersion] = cachedVal
		}

		return cachedVal.resourceList, cachedVal.err
	}

	liveResources, err := d.delegate.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return liveResources, err
	}

	cachedVal = &cacheEntry{liveResources, err}
	d.groupToServerResources[groupVersion] = cachedVal

	return cachedVal.resourceList, cachedVal.err
}

// ServerResources returns the supported resources for all groups and versions.
// Deprecated: use ServerGroupsAndResources instead.
func (d *memCacheClient) ServerResources() ([]*metav1.APIResourceList, error) {
	return discovery.ServerResources(d)
}

// ServerGroupsAndResources returns the groups and supported resources for all groups and versions.
func (d *memCacheClient) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return discovery.ServerGroupsAndResources(d)
}

func (d *memCacheClient) ServerGroups() (*metav1.APIGroupList, error) {
	cachedEntry, err := d.getCachedEntry("servergroups")
	// d.lock.Lock()
	// defer d.lock.Unlock()

	// if !d.cacheValid {
	// 	if err := d.refreshLocked(); err != nil {
	// 		return nil, err
	// 	}
	// }

	// if d.groupList != nil {
	// 	return d.groupList, nil
	// }
	//
	// liveGroups, err := d.delegate.ServerGroups()
	// if err != nil {
	// 	return liveGroups, err
	// }
	//
	// if liveGroups == nil || len(liveGroups.Groups) == 0 {
	// 	return liveGroups, err
	// }
	//
	// d.groupList = liveGroups

	// cachedVal, ok := d.groupToServerResources[groupVersion]
	// if ok {
	// 	if cachedVal.err != nil && isTransientError(cachedVal.err) {
	// 		r, err := d.serverResourcesForGroupVersion(groupVersion)
	// 		if err != nil {
	// 			utilruntime.HandleError(fmt.Errorf("couldn't get resource list for %v: %v", groupVersion, err))
	// 		}
	// 		cachedVal = &cacheEntry{r, err}
	// 		d.groupToServerResources[groupVersion] = cachedVal
	// 	}
	//
	// 	return cachedVal.resourceList, cachedVal.err
	// }
	//
	// r, err := d.serverResourcesForGroupVersion(groupVersion)
	// if err != nil {
	// 	utilruntime.HandleError(fmt.Errorf("couldn't get resource list for %v: %v", groupVersion, err))
	// }
	// cachedVal = &cacheEntry{r, err}
	// d.groupToServerResources[groupVersion] = cachedVal
	//
	// return cachedVal.resourceList, cachedVal.err

	// return d.groupList, nil
}

func (d *memCacheClient) getCachedEntry(key string) (*cacheEntry, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cachedEntry, ok := d.ourEntries[key]
	if d.invalidated && !ok {
		return nil, errors.New("cache invalidated")
	}

	t, ok := d.ourTTLs[key]
	if ok && time.Now().After(t.Add(d.ttl)) {
		return nil, errors.New("cache expired")
	}

	d.fresh = d.fresh && ok

	return cachedEntry, nil
}

func (d *memCacheClient) createCachedEntry(key string, entry *cacheEntry) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.ourEntries[key] = entry
	d.ourTTLs[key] = time.Now().UTC().Add(d.ttl)
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
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.fresh
}

// Invalidate enforces that no cached data that is older than the current time
// is used.
func (d *memCacheClient) Invalidate() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.ourEntries = map[string]*cacheEntry{}
	d.ourTTLs = map[string]time.Duration{}
	d.fresh = true
	d.invalidated = true
}

// NewMemCacheClient creates a new CachedDiscoveryInterface which caches
// discovery information in memory and will stay up-to-date if Invalidate is
// called with regularity.
func NewMemCacheClient(delegate discovery.DiscoveryInterface, ttl time.Duration) discovery.CachedDiscoveryInterface {
	return &memCacheClient{
		delegate: delegate,
		ttl:      ttl,
		// groupToServerResources: map[string]*cacheEntry{},
		ourEntries: map[string]*cacheEntry{},
		ourTTLs:    map[string]time.Duration{},
		fresh:      true,
	}
}
