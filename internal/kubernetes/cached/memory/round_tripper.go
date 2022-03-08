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
	"net/http"

	"github.com/gregjones/httpcache"
	"k8s.io/klog/v2"
)

//go:generate counterfeiter . CacheRoundTripper
type MemCacheRoundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
	CancelRequest(req *http.Request)
}

type memCacheRoundTripper struct {
	rt *httpcache.Transport
}

func NewMemCacheRoundTripper(rt http.RoundTripper, c *httpcache.MemoryCache) http.RoundTripper {
	return newMemCacheRoundTripper(rt, c)
}

// newMemCacheRoundTripper creates a roundtripper that reads the ETag on
// response headers and send the If-None-Match header on subsequent
// corresponding requests.
func newMemCacheRoundTripper(rt http.RoundTripper, c *httpcache.MemoryCache) http.RoundTripper {
	t := httpcache.NewTransport(c)
	t.Transport = rt

	return &memCacheRoundTripper{rt: t}
}

func (rt *memCacheRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.rt.RoundTrip(req)
}

func (rt *memCacheRoundTripper) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}

	if cr, ok := rt.rt.Transport.(canceler); ok {
		cr.CancelRequest(req)
	} else {
		klog.Errorf("CancelRequest not implemented by %T", rt.rt.Transport)
	}
}

func (rt *memCacheRoundTripper) WrappedRoundTripper() http.RoundTripper { return rt.rt.Transport }
