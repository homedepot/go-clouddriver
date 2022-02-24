package kubernetes

import (
	"time"

	"github.com/homedepot/go-clouddriver/internal/kubernetes/cached/memory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
)

var _ = Describe("Controller", func() {
	var (
		client     Client
		clientset  Clientset
		config     *rest.Config
		controller Controller
		err        error
	)

	Describe("#NewClient", func() {
		BeforeEach(func() {
			cachedMemCacheClients = map[string]memory.CachedDiscoveryClient{}
			config = &rest.Config{
				Host:        "https://test-host",
				BearerToken: "some.bearer.token",
				TLSClientConfig: rest.TLSClientConfig{
					CAData: []byte("test-ca-data"),
				},
			}
			controller = NewController()
		})

		JustBeforeEach(func() {
			client, err = controller.NewClient(config)
		})

		Context("memory cache", func() {
			When("generating the dynamic client returns an error", func() {
				BeforeEach(func() {
					config = &rest.Config{
						Host:        ":::badhost;",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte("test-ca-data"),
						},
					}
				})

				It("returns an error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("parse \"https://:::badhost;\": invalid port \":badhost;\" after host"))
				})
			})

			When("a call is made for a cached client", func() {
				JustBeforeEach(func() {
					client, err = controller.NewClient(config)
				})

				It("returns and cached client", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(cachedMemCacheClients).To(HaveLen(1))
					cachedClient := cachedMemCacheClients[config.Host]
					Expect(cachedClient.Fresh()).To(BeTrue())
				})
			})

			When("the bearer token for a client changes", func() {
				JustBeforeEach(func() {
					newConfig := &rest.Config{
						Host:        "https://test-host",
						BearerToken: "another.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte("test-ca-data"),
						},
					}
					client, err = controller.NewClient(newConfig)
				})

				It("resets the cache and returns the new client", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(cachedMemCacheClients).To(HaveLen(1))
					cachedClient := cachedMemCacheClients[config.Host]
					Expect(cachedClient.Fresh()).To(BeTrue())
				})
			})

			When("the CAData for a client changes", func() {
				JustBeforeEach(func() {
					newConfig := &rest.Config{
						Host:        "https://test-host",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte("different-ca-data"),
						},
					}
					client, err = controller.NewClient(newConfig)
				})

				It("resets the cache and returns the new client", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(cachedMemCacheClients).To(HaveLen(1))
					cachedClient := cachedMemCacheClients[config.Host]
					Expect(cachedClient.Fresh()).To(BeTrue())
				})
			})

			When("the same host has two defined timeouts", func() {
				JustBeforeEach(func() {
					newConfig := &rest.Config{
						Host:        "https://test-host",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte("test-ca-data"),
						},
						Timeout: 1 * time.Second,
					}
					client, err = controller.NewClient(newConfig)
				})

				It("caches a single copy of the client", func() {
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					Expect(cachedMemCacheClients).To(HaveLen(1))
				})
			})

			It("returns and caches client", func() {
				Expect(err).To(BeNil())
				Expect(client).ToNot(BeNil())
				Expect(cachedMemCacheClients).To(HaveLen(1))
				cachedClient := cachedMemCacheClients[config.Host]
				Expect(cachedClient.Fresh()).To(BeTrue())
			})
		})
	})

	Describe("#NewClientset", func() {
		BeforeEach(func() {
			config = &rest.Config{
				Host:        "https://test-host",
				BearerToken: "some.bearer.token",
				TLSClientConfig: rest.TLSClientConfig{
					CAData: []byte("test-ca-data"),
				},
			}
			controller = NewController()
		})

		JustBeforeEach(func() {
			clientset, err = controller.NewClientset(config)
		})

		Context("memory cache", func() {
			When("generating the clientset returns an error", func() {
				BeforeEach(func() {
					config = &rest.Config{
						Host:        ":::badhost;",
						BearerToken: "some.bearer.token",
						TLSClientConfig: rest.TLSClientConfig{
							CAData: []byte("test-ca-data"),
						},
					}
				})

				It("returns an error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("parse \"https://:::badhost;\": invalid port \":badhost;\" after host"))
				})
			})

			It("returns the clientset", func() {
				Expect(err).To(BeNil())
				Expect(clientset).ToNot(BeNil())
			})
		})
	})
})
