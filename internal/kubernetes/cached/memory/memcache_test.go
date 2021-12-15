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

package memory_test

import (
	"time"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes/cached/memory"
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
		err error
		c   *fakeDiscoveryClient
		mc  CachedDiscoveryClient
	)

	Describe("#Fresh", func() {
		BeforeEach(func() {
			c = &fakeDiscoveryClient{}
			mc = NewCachedDiscoveryClient(c, 60*time.Second)
		})

		JustBeforeEach(func() {
		})

		When("the client is created", func() {
			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
			})
		})

		When("server groups is called", func() {
			BeforeEach(func() {
				_, err = mc.ServerGroups()
				Expect(err).To(BeNil())
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(c.groupCalls).To(Equal(1))
			})
		})

		When("server groups is called twice", func() {
			BeforeEach(func() {
				mc.ServerGroups()
				mc.ServerGroups()
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(c.groupCalls).To(Equal(1))
			})
		})

		When("resources is called", func() {
			BeforeEach(func() {
				mc.ServerResources()
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(c.resourceCalls).To(Equal(2))
			})
		})

		When("resources is called twice", func() {
			BeforeEach(func() {
				mc.ServerResources()
				mc.ServerResources()
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(c.resourceCalls).To(Equal(2))
			})
		})

		When("the cache is valid but a request is made for a non-existing resource", func() {
			BeforeEach(func() {
				mc.ServerResources()
				mc.ServerResourcesForGroupVersion("c/v1")
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(c.resourceCalls).To(Equal(3))
			})
		})

		Context("client is reset", func() {
			BeforeEach(func() {
				mc.ServerGroups()
				mc.ServerResources()
				mc.Reset()
			})

			When("server groups is called", func() {
				BeforeEach(func() {
					mc.ServerGroups()
				})

				It("should not be fresh", func() {
					Expect(mc.Fresh()).To(BeFalse())
					Expect(c.groupCalls).To(Equal(1))
				})
			})

			When("resources is called", func() {
				BeforeEach(func() {
					mc.ServerResources()
				})

				It("should not be fresh", func() {
					Expect(mc.Fresh()).To(BeFalse())
					Expect(c.resourceCalls).To(Equal(2))
				})
			})

			When("invalidate is called", func() {
				BeforeEach(func() {
					mc.Invalidate()
				})

				It("should be fresh", func() {
					Expect(mc.Fresh()).To(BeTrue())
				})

				It("should ignore existing resources cached after validation", func() {
					Expect(mc.Fresh()).To(BeTrue())
					Expect(c.resourceCalls).To(Equal(2))
				})

				It("should ignore existing resources cache after invalidation", func() {
					_, err = mc.ServerResources()
					Expect(err).To(BeNil())
					Expect(mc.Fresh()).To(BeTrue())
					Expect(c.resourceCalls).To(Equal(4))
				})
			})
		})
	})

	Describe("#TTL", func() {
		BeforeEach(func() {
			c = &fakeDiscoveryClient{}
			mc = NewCachedDiscoveryClient(c, 1*time.Second)
		})

		JustBeforeEach(func() {
		})

		It("respects the ttl", func() {
			_, err = mc.ServerGroups()
			Expect(err).To(BeNil())
			Expect(c.groupCalls).To(Equal(1))
			time.Sleep(1 * time.Second)
			_, err = mc.ServerGroups()
			Expect(err).To(BeNil())
			Expect(c.groupCalls).To(Equal(2))
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
			{
				Name: "b",
				Versions: []metav1.GroupVersionForDiscovery{
					{
						GroupVersion: "b/v1",
						Version:      "v1",
					},
				},
				PreferredVersion: metav1.GroupVersionForDiscovery{
					GroupVersion: "b/v1",
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
	if groupVersion == "b/v1" {
		return &metav1.APIResourceList{APIResources: []metav1.APIResource{{Name: "bidgets", Kind: "Bidget"}}}, nil
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
