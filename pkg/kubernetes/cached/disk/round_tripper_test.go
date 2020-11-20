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

package disk_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	. "github.com/homedepot/go-clouddriver/pkg/kubernetes/cached/disk"
	. "github.com/onsi/ginkgo"
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
		rt       *testRoundTripper
		cacheDir string
		err      error
		cache    http.RoundTripper
		req      *http.Request
		resp     *http.Response
		content  []byte
	)

	BeforeEach(func() {
		rt = &testRoundTripper{}
		cacheDir, err = ioutil.TempDir("", "cache-rt")
		Expect(err).To(BeNil())
	})

	JustBeforeEach(func() {
		cache = NewCacheRoundTripper(cacheDir, rt)
	})

	AfterEach(func() {
		os.RemoveAll(cacheDir)
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

		When("the cache directory is removed", func() {
			BeforeEach(func() {
				os.RemoveAll(cacheDir)
				rt.Response = &http.Response{
					Header:     http.Header{"ETag": []string{`"123456"`}},
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("Content"))),
					StatusCode: http.StatusOK,
				}
			})

			It("creates the cache directories and files with the correct permissions", func() {
				err = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						if info.Mode().Perm() != os.FileMode(0750) {
							return fmt.Errorf("directory perm incorrect expected %d got %d", os.FileMode(0750), info.Mode().Perm())
						}
					} else {
						if info.Mode().Perm() != os.FileMode(0660) {
							return fmt.Errorf("file perm incorrect expected %d got %d", os.FileMode(0660), info.Mode().Perm())
						}
					}
					return nil
				})
				Expect(err).To(BeNil())
			})
		})
	})
})
