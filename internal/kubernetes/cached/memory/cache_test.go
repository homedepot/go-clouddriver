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
	"fmt"
	"net/http"
	"time"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes/cached/memory"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/gomega"

	restclient "k8s.io/client-go/rest"
)

var _ = Describe("CachedDiscovery", func() {
	var (
		err         error
		config      *restclient.Config
		fakeCluster *ghttp.Server
		mc          CachedDiscoveryClient
		c           *Cache
	)

	BeforeEach(func() {
		fakeCluster = ghttp.NewServer()
		fakeCluster.RouteToHandler(http.MethodGet, "/api", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(fmt.Sprintf(`{
								"kind": "APIVersions",
								"versions": [
									"v1"
								],
								"serverAddressByClientCIDRs": [
									{
										"clientCIDR": "0.0.0.0/0",
										"serverAddress": "%s"
									}
								]
							}`, fakeCluster.URL()),
			))
		})
		fakeCluster.RouteToHandler(http.MethodGet, "/apis", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`{
							"kind": "APIGroupList",
							"apiVersion": "v1",
							"groups": [
								{
									"name": "apps",
									"versions": [
										{
											"groupVersion": "apps/v1",
											"version": "v1"
										}
									],
									"preferredVersion": {
										"groupVersion": "apps/v1",
										"version": "v1"
									}
								}
							]
						}`),
			)
		})
		fakeCluster.RouteToHandler(http.MethodGet, "/apis/apps/v1", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`
					{
						"kind": "APIResourceList",
						"apiVersion": "v1",
						"groupVersion": "apps/v1",
						"resources": [
							{
								"name": "deployments",
								"singularName": "",
								"namespaced": true,
								"kind": "Deployment",
								"verbs": [
									"create",
									"delete",
									"deletecollection",
									"get",
									"list",
									"patch",
									"update",
									"watch"
								],
								"shortNames": [
									"deploy"
								],
								"categories": [
									"all"
								],
								"storageVersionHash": "asdf"
							}
						]
					}
					`),
			)
		})
		fakeCluster.RouteToHandler(http.MethodGet, "/apis/c/v1", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`
					{
						"kind": "APIResourceList",
						"apiVersion": "v1",
						"groupVersion": "c/v1",
						"resources": [
							{
								"name": "c",
								"singularName": "",
								"namespaced": true,
								"kind": "C",
								"verbs": [
									"create",
									"delete",
									"deletecollection",
									"get",
									"list",
									"patch",
									"update",
									"watch"
								],
								"shortNames": [
									"deploy"
								],
								"categories": [
									"all"
								],
								"storageVersionHash": "asdf"
							}
						]
					}
					`),
			)
		})
		fakeCluster.RouteToHandler(http.MethodGet, "/api/v1", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`
						{
							"kind": "APIResourceList",
							"groupVersion": "v1",
							"resources": [
								{
									"name": "namespaces",
									"singularName": "",
									"namespaced": false,
									"kind": "Namespace",
									"verbs": [
										"create",
										"delete",
										"get",
										"list",
										"patch",
										"update",
										"watch"
									],
									"shortNames": [
										"ns"
									],
									"storageVersionHash": "Q3oi5N2YM8M="
								}
							]
						}
					`),
			)
		})
		config = &restclient.Config{
			Host: fakeCluster.URL(),
		}
		c = NewCache(1 * time.Second)
		mc, err = c.NewClientForConfig(config)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		fakeCluster.Close()
	})

	Describe("#Fresh", func() {
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
				Expect(fakeCluster.ReceivedRequests()).To(HaveLen(2))
				Expect(fakeCluster.ReceivedRequests()[0].URL.Path).To(Equal("/api"))
				Expect(fakeCluster.ReceivedRequests()[1].URL.Path).To(Equal("/apis"))
			})
		})

		When("server groups is called twice", func() {
			BeforeEach(func() {
				mc.ServerGroups()
				mc.ServerGroups()
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(fakeCluster.ReceivedRequests()).To(HaveLen(2))
			})
		})

		When("resources is called", func() {
			BeforeEach(func() {
				mc.ServerResources()
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(fakeCluster.ReceivedRequests()).To(HaveLen(4))
			})
		})

		When("resources is called twice", func() {
			BeforeEach(func() {
				mc.ServerResources()
				mc.ServerResources()
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(fakeCluster.ReceivedRequests()).To(HaveLen(4))
			})
		})

		When("the cache is valid but a request is made for a non-existing resource", func() {
			BeforeEach(func() {
				mc.ServerResources()
				mc.ServerResourcesForGroupVersion("c/v1")
			})

			It("should be fresh", func() {
				Expect(mc.Fresh()).To(BeTrue())
				Expect(fakeCluster.ReceivedRequests()).To(HaveLen(5))
			})
		})

		Context("client is recreated", func() {
			BeforeEach(func() {
				mc.ServerGroups()
				mc.ServerResources()
				mc, err = c.NewClientForConfig(config)
				Expect(err).To(BeNil())
			})

			When("server groups is called", func() {
				BeforeEach(func() {
					mc.ServerGroups()
				})

				It("should not be fresh", func() {
					Expect(mc.Fresh()).To(BeFalse())
					Expect(fakeCluster.ReceivedRequests()).To(HaveLen(4))
				})
			})

			When("resources is called", func() {
				BeforeEach(func() {
					mc.ServerResources()
				})

				It("should not be fresh", func() {
					Expect(mc.Fresh()).To(BeFalse())
					Expect(fakeCluster.ReceivedRequests()).To(HaveLen(4))
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
					Expect(fakeCluster.ReceivedRequests()).To(HaveLen(4))
				})

				It("should ignore existing resources cache after invalidation", func() {
					_, err = mc.ServerResources()
					Expect(err).To(BeNil())
					Expect(mc.Fresh()).To(BeTrue())
					Expect(fakeCluster.ReceivedRequests()).To(HaveLen(8))
				})
			})
		})
	})

	Describe("#TTL", func() {
		It("respects the ttl", func() {
			_, err = mc.ServerGroups()
			Expect(err).To(BeNil())
			Expect(fakeCluster.ReceivedRequests()).To(HaveLen(2))
			time.Sleep(1 * time.Second)
			_, err = mc.ServerGroups()
			Expect(err).To(BeNil())
			Expect(fakeCluster.ReceivedRequests()).To(HaveLen(4))
		})
	})
})
