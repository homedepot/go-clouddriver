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

package disk

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
	"k8s.io/klog/v2"
)

//go:generate counterfeiter . CacheRoundTripper
type CacheRoundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
	CancelRequest(req *http.Request)
}

type cacheRoundTripper struct {
	rt *httpcache.Transport
}

func NewCacheRoundTripper(cachDir string, rt http.RoundTripper) http.RoundTripper {
	return newCacheRoundTripper(cachDir, rt)
}

// newCacheRoundTripper creates a roundtripper that reads the ETag on
// response headers and send the If-None-Match header on subsequent
// corresponding requests.
func newCacheRoundTripper(cacheDir string, rt http.RoundTripper) http.RoundTripper {
	d := diskv.New(diskv.Options{
		PathPerm: os.FileMode(0750),
		FilePerm: os.FileMode(0660),
		BasePath: cacheDir,
		TempDir:  filepath.Join(cacheDir, ".diskv-temp"),
	})
	t := httpcache.NewTransport(diskcache.NewWithDiskv(d))
	t.Transport = rt

	return &cacheRoundTripper{rt: t}
}

func (rt *cacheRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.rt.RoundTrip(req)
}

func (rt *cacheRoundTripper) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}

	if cr, ok := rt.rt.Transport.(canceler); ok {
		cr.CancelRequest(req)
	} else {
		klog.Errorf("CancelRequest not implemented by %T", rt.rt.Transport)
	}
}

func (rt *cacheRoundTripper) WrappedRoundTripper() http.RoundTripper { return rt.rt.Transport }
