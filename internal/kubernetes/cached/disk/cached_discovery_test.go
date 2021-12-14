/*
Copyright 2016 The Kubernetes Authors.

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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes/cached/disk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	openapi_v2 "github.com/googleapis/gnostic/openapiv2"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/rest/fake"
)

var _ = Describe("CachedDiscovery", func() {
	var (
		d   string
		err error
		c   *fakeDiscoveryClient
		cdc CachedDiscoveryClient
	)

	Describe("#Fresh", func() {
		BeforeEach(func() {
			d, err = ioutil.TempDir("", "")
			Expect(err).To(BeNil())
			c = &fakeDiscoveryClient{}
			cdc = NewCachedDiscoveryClient(c, d, 60*time.Second)
		})

		JustBeforeEach(func() {
		})

		AfterEach(func() {
			os.RemoveAll(d)
		})

		When("the client is created", func() {
			It("should be fresh", func() {
				Expect(cdc.Fresh()).To(BeTrue())
			})
		})

		When("server groups is called", func() {
			BeforeEach(func() {
				_, err = cdc.ServerGroups()
				Expect(err).To(BeNil())
			})

			It("should be fresh", func() {
				Expect(cdc.Fresh()).To(BeTrue())
				Expect(c.groupCalls).To(Equal(1))
			})
		})

		When("server groups is called twice", func() {
			BeforeEach(func() {
				cdc.ServerGroups()
				cdc.ServerGroups()
			})

			It("should be fresh", func() {
				Expect(cdc.Fresh()).To(BeTrue())
				Expect(c.groupCalls).To(Equal(1))
			})
		})

		When("resources is called", func() {
			BeforeEach(func() {
				cdc.ServerResources()
			})

			It("should be fresh", func() {
				Expect(cdc.Fresh()).To(BeTrue())
				Expect(c.resourceCalls).To(Equal(1))
			})
		})

		When("resources is called twice", func() {
			BeforeEach(func() {
				cdc.ServerResources()
				cdc.ServerResources()
			})

			It("should be fresh", func() {
				Expect(cdc.Fresh()).To(BeTrue())
				Expect(c.resourceCalls).To(Equal(1))
			})
		})

		Context("client is recreated", func() {
			BeforeEach(func() {
				cdc.ServerGroups()
				cdc.ServerResources()
				cdc = NewCachedDiscoveryClient(c, d, 60*time.Second)
			})

			When("server groups is called", func() {
				BeforeEach(func() {
					cdc.ServerGroups()
				})

				It("should not be fresh", func() {
					Expect(cdc.Fresh()).To(BeFalse())
					Expect(c.groupCalls).To(Equal(1))
				})
			})

			When("resources is called", func() {
				BeforeEach(func() {
					cdc.ServerResources()
				})

				It("should not be fresh", func() {
					Expect(cdc.Fresh()).To(BeFalse())
					Expect(c.resourceCalls).To(Equal(1))
				})
			})

			When("invalidate is called", func() {
				BeforeEach(func() {
					cdc.Invalidate()
				})

				It("should be fresh", func() {
					Expect(cdc.Fresh()).To(BeTrue())
				})

				It("should ignore existing resources cached after validation", func() {
					Expect(cdc.Fresh()).To(BeTrue())
					Expect(c.resourceCalls).To(Equal(1))
				})

				It("should ignore existing resources cache after invalidation", func() {
					_, err = cdc.ServerResources()
					Expect(err).To(BeNil())
					Expect(cdc.Fresh()).To(BeTrue())
					Expect(c.resourceCalls).To(Equal(2))
				})
			})
		})
	})

	Describe("#TTL", func() {
		BeforeEach(func() {
			d, err = ioutil.TempDir("", "")
			Expect(err).To(BeNil())
			c = &fakeDiscoveryClient{}
			cdc = NewCachedDiscoveryClient(c, d, 1*time.Nanosecond)
		})

		JustBeforeEach(func() {
		})

		AfterEach(func() {
			os.RemoveAll(d)
		})

		It("respects the ttl", func() {
			_, err = cdc.ServerGroups()
			Expect(err).To(BeNil())
			Expect(c.groupCalls).To(Equal(1))
			time.Sleep(1 * time.Second)
			_, err = cdc.ServerGroups()
			Expect(err).To(BeNil())
			Expect(c.groupCalls).To(Equal(2))
		})
	})

	Describe("#PathPerm", func() {
		BeforeEach(func() {
			d, err = ioutil.TempDir("", "")
			Expect(err).To(BeNil())
			os.RemoveAll(d)
			c = &fakeDiscoveryClient{}
			cdc = NewCachedDiscoveryClient(c, d, 1*time.Nanosecond)
			cdc.ServerGroups()
		})

		JustBeforeEach(func() {
		})

		AfterEach(func() {
			os.RemoveAll(d)
		})

		When("it succeeds", func() {
			It("creates the cache directories and files with the correct permissions", func() {
				err = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
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

type fakeDiscoveryClient struct {
	groupCalls    int
	resourceCalls int
	versionCalls  int
	openAPICalls  int

	serverResourcesHandler func() ([]*metav1.APIResourceList, error)
}

var _ discovery.DiscoveryInterface = &fakeDiscoveryClient{}

func (c *fakeDiscoveryClient) RESTClient() restclient.Interface {
	return &fake.RESTClient{}
}

func (c *fakeDiscoveryClient) ServerGroups() (*metav1.APIGroupList, error) {
	c.groupCalls = c.groupCalls + 1
	return c.serverGroups()
}

func (c *fakeDiscoveryClient) serverGroups() (*metav1.APIGroupList, error) {
	return &metav1.APIGroupList{
		Groups: []metav1.APIGroup{
			{
				Name: "a",
				Versions: []metav1.GroupVersionForDiscovery{
					{
						GroupVersion: "a/v1",
						Version:      "v1",
					},
				},
				PreferredVersion: metav1.GroupVersionForDiscovery{
					GroupVersion: "a/v1",
					Version:      "v1",
				},
			},
		},
	}, nil
}

func (c *fakeDiscoveryClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	c.resourceCalls = c.resourceCalls + 1

	if groupVersion == "a/v1" {
		return &metav1.APIResourceList{APIResources: []metav1.APIResource{{Name: "widgets", Kind: "Widget"}}}, nil
	}

	return nil, errors.NewNotFound(schema.GroupResource{}, "")
}

// Deprecated: use ServerGroupsAndResources instead.
func (c *fakeDiscoveryClient) ServerResources() ([]*metav1.APIResourceList, error) {
	_, rs, err := c.ServerGroupsAndResources()
	return rs, err
}

func (c *fakeDiscoveryClient) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	c.resourceCalls = c.resourceCalls + 1

	gs, _ := c.serverGroups()
	resultGroups := []*metav1.APIGroup{}

	for i := range gs.Groups {
		resultGroups = append(resultGroups, &gs.Groups[i])
	}

	if c.serverResourcesHandler != nil {
		rs, err := c.serverResourcesHandler()
		return resultGroups, rs, err
	}

	return resultGroups, []*metav1.APIResourceList{}, nil
}

func (c *fakeDiscoveryClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	c.resourceCalls = c.resourceCalls + 1
	return nil, nil
}

func (c *fakeDiscoveryClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	c.resourceCalls = c.resourceCalls + 1
	return nil, nil
}

func (c *fakeDiscoveryClient) ServerVersion() (*version.Info, error) {
	c.versionCalls = c.versionCalls + 1
	return &version.Info{}, nil
}

func (c *fakeDiscoveryClient) OpenAPISchema() (*openapi_v2.Document, error) {
	c.openAPICalls = c.openAPICalls + 1
	return &openapi_v2.Document{}, nil
}
