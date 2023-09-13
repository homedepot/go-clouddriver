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

package memory_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gregjones/httpcache"
	. "github.com/homedepot/go-clouddriver/internal/kubernetes/cached/memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// copied from k8s.io/client-go/transport/round_trippers_test.go
type testRoundTripper struct {
	Request  *http.Request
	Response *http.Response
	Err      error
}

func (rt *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.Request = req
	return rt.Response, rt.Err
}

var _ = Describe("CachedDiscovery", func() {
	var (
		rt           *testRoundTripper
		err          error
		cache        http.RoundTripper
		httpMemCache *httpcache.MemoryCache
		req          *http.Request
		resp         *http.Response
		content      []byte
	)

	BeforeEach(func() {
		rt = &testRoundTripper{}
		httpMemCache = httpcache.NewMemoryCache()
	})

	JustBeforeEach(func() {
		cache = NewMemCacheRoundTripper(rt, httpMemCache)
	})

	AfterEach(func() {
		if resp != nil {
			resp.Body.Close()
		}
	})

	Describe("#RoundTrip", func() {
		BeforeEach(func() {
			req = &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Host: "localhost"},
			}
		})

		JustBeforeEach(func() {
			resp, err = cache.RoundTrip(req)
			Expect(err).To(BeNil())
			content, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
		})

		When("the data is not cached", func() {
			BeforeEach(func() {
				rt.Response = &http.Response{
					Header:     http.Header{"ETag": []string{`"123456"`}},
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("Content"))),
					StatusCode: http.StatusOK,
				}
			})

			It("succeeds", func() {
				Expect(string(content)).To(Equal("Content"))
			})
		})

		When("the data is cached", func() {
			BeforeEach(func() {
				rt.Response = &http.Response{
					StatusCode: http.StatusNotModified,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("Other Content"))),
				}
			})

			It("succeeds", func() {
				Expect(string(content)).To(Equal("Other Content"))
			})
		})
	})
})
