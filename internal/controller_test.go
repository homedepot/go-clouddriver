package internal_test

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/arcade/arcadefakes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/kubernetesfakes"
	"github.com/homedepot/go-clouddriver/internal/sql/sqlfakes"
)

var _ = Describe("Controller", func() {
	Describe("#const", func() {
		Expect(internal.DefaultListTimeoutSeconds).To(Equal(10))
		Expect(internal.DefaultChanSize).To(Equal(100000))
	})

	var (
		c                        *internal.Controller
		fakeSQLClient            *sqlfakes.FakeClient
		fakeArcadeClient         *arcadefakes.FakeClient
		fakeKubernetesController *kubernetesfakes.FakeController
		fakeKubernetesClient     *kubernetesfakes.FakeClient
		fakeKubernetesClientset  *kubernetesfakes.FakeClientset
		provider                 *kubernetes.Provider
		providers                []*kubernetes.Provider
		err                      error
	)

	BeforeEach(func() {
		fakeSQLClient = &sqlfakes.FakeClient{}
		fakeArcadeClient = &arcadefakes.FakeClient{}
		fakeKubernetesController = &kubernetesfakes.FakeController{}
		fakeKubernetesClient = &kubernetesfakes.FakeClient{}
		fakeKubernetesClientset = &kubernetesfakes.FakeClientset{}

		fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
			Name:          "test-name",
			Host:          "test-host",
			CAData:        "12341234",
			TokenProvider: "test-token-provider",
		}, nil)
		fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
			{
				Name:          "test-name1",
				Host:          "test-host1",
				CAData:        "12341234",
				TokenProvider: "test-token-provider1",
			},
			{
				Name:          "test-name2",
				Host:          "test-host2",
				CAData:        "56785678",
				TokenProvider: "test-token-provider2",
			},
		}, nil)

		fakeKubernetesController.NewClientReturns(fakeKubernetesClient, nil)
		fakeKubernetesController.NewClientsetReturns(fakeKubernetesClientset, nil)

		c = &internal.Controller{
			ArcadeClient:         fakeArcadeClient,
			KubernetesController: fakeKubernetesController,
			SQLClient:            fakeSQLClient,
		}
		log.SetOutput(ioutil.Discard)
	})

	Describe("#KubernetesProvider", func() {
		JustBeforeEach(func() {
			provider, err = c.KubernetesProvider("test-account")
		})

		It("succeeds", func() {
			Expect(err).To(BeNil())
			Expect(provider).ToNot(BeNil())
			Expect(provider.Name).To(Equal("test-name"))
			Expect(provider.Host).To(Equal("test-host"))
			Expect(provider.CAData).To(Equal("12341234"))
			Expect(provider.TokenProvider).To(Equal("test-token-provider"))
			Expect(provider.Client).ToNot(BeNil())
			Expect(provider.Clientset).ToNot(BeNil())
			config := fakeKubernetesController.NewClientArgsForCall(0)
			Expect(config.Timeout).To(BeZero())
		})
	})

	Describe("#KubernetesProviderWithTimeout", func() {
		JustBeforeEach(func() {
			provider, err = c.KubernetesProviderWithTimeout("test-account", time.Second*1)
		})

		When("getting the provider from sql returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error getting kubernetes provider test-account: error getting provider"))
			})
		})

		When("the ca data is bad", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					CAData: "{}",
				}, nil)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error decoding provider CA data: illegal base64 data at input byte 0"))
			})
		})

		When("getting the arcade token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error getting token from arcade for provider test-token-provider: error getting token"))
			})
		})

		When("generating a new client returns an error", func() {
			BeforeEach(func() {
				fakeKubernetesController.NewClientReturns(nil, errors.New("error generating client"))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error creating new kubernetes client: error generating client"))
			})
		})

		When("generating a new clientset returns an error", func() {
			BeforeEach(func() {
				fakeKubernetesController.NewClientsetReturns(nil, errors.New("error generating clientset"))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error creating new kubernetes clientset: error generating clientset"))
			})
		})

		It("succeeds", func() {
			Expect(err).To(BeNil())
			Expect(provider).ToNot(BeNil())
			Expect(provider.Name).To(Equal("test-name"))
			Expect(provider.Host).To(Equal("test-host"))
			Expect(provider.CAData).To(Equal("12341234"))
			Expect(provider.TokenProvider).To(Equal("test-token-provider"))
			Expect(provider.Client).ToNot(BeNil())
			Expect(provider.Clientset).ToNot(BeNil())
			config := fakeKubernetesController.NewClientArgsForCall(0)
			Expect(config.Timeout).To(Equal(time.Second * 1))
		})
	})

	Describe("#KubernetesProvidersForAccountswithTimeout", func() {
		var accounts []string

		BeforeEach(func() {
			accounts = []string{
				"test-name1",
				"test-name2",
			}
		})

		JustBeforeEach(func() {
			providers, err = c.KubernetesProvidersForAccountsWithTimeout(accounts, time.Second*1)
		})

		When("listing the providers from sql returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{}, errors.New("error listing providers"))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error listing kubernetes providers: error listing providers"))
			})
		})

		When("the ca data is bad", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
					{
						Name:   "test-name1",
						CAData: "{}",
					},
					{
						Name:   "test-name2",
						CAData: "12341234",
					},
				}, nil)
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(1))
			})
		})

		When("getting the arcade token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		When("generating a new client returns an error", func() {
			BeforeEach(func() {
				fakeKubernetesController.NewClientReturns(nil, errors.New("error generating client"))
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		When("generating a new clientset returns an error", func() {
			BeforeEach(func() {
				fakeKubernetesController.NewClientsetReturns(nil, errors.New("error generating clientset"))
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		When("a subset of accounts are requested", func() {
			BeforeEach(func() {
				accounts = []string{
					"test-name2",
				}
			})

			It("only returns the requested providers", func() {
				Expect(err).To(BeNil())
				Expect(providers).ToNot(BeNil())
				Expect(providers).To(HaveLen(1))
				Expect(providers[0]).ToNot(BeNil())
				Expect(providers[0].Name).To(Equal("test-name2"))
				Expect(providers[0].Host).To(Equal("test-host2"))
				Expect(providers[0].CAData).To(Equal("56785678"))
				Expect(providers[0].TokenProvider).To(Equal("test-token-provider2"))
				Expect(providers[0].Client).ToNot(BeNil())
				Expect(providers[0].Clientset).ToNot(BeNil())
				config := fakeKubernetesController.NewClientArgsForCall(0)
				Expect(config.Timeout).To(Equal(time.Second * 1))
			})
		})

		It("succeeds", func() {
			Expect(err).To(BeNil())
			Expect(providers).ToNot(BeNil())
			Expect(providers).To(HaveLen(2))
			Expect(providers[0]).ToNot(BeNil())
			Expect(providers[1]).ToNot(BeNil())
			Expect(providers[0].Name).To(Equal("test-name1"))
			Expect(providers[0].Host).To(Equal("test-host1"))
			Expect(providers[0].CAData).To(Equal("12341234"))
			Expect(providers[0].TokenProvider).To(Equal("test-token-provider1"))
			Expect(providers[0].Client).ToNot(BeNil())
			Expect(providers[0].Clientset).ToNot(BeNil())
			Expect(providers[1].Name).To(Equal("test-name2"))
			Expect(providers[1].Host).To(Equal("test-host2"))
			Expect(providers[1].CAData).To(Equal("56785678"))
			Expect(providers[1].TokenProvider).To(Equal("test-token-provider2"))
			Expect(providers[1].Client).ToNot(BeNil())
			Expect(providers[1].Clientset).ToNot(BeNil())
			config := fakeKubernetesController.NewClientArgsForCall(0)
			Expect(config.Timeout).To(Equal(time.Second * 1))
			config = fakeKubernetesController.NewClientArgsForCall(1)
			Expect(config.Timeout).To(Equal(time.Second * 1))
		})
	})

	Describe("#AllKubernetesProvidersWithTimeout", func() {
		JustBeforeEach(func() {
			providers, err = c.AllKubernetesProvidersWithTimeout(time.Second * 1)
		})

		When("listing the providers from sql returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{}, errors.New("error listing providers"))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal: error listing kubernetes providers: error listing providers"))
			})
		})

		When("the ca data is bad", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
					{
						CAData: "{}",
					},
				}, nil)
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		When("getting the arcade token returns an error", func() {
			BeforeEach(func() {
				fakeArcadeClient.TokenReturns("", errors.New("error getting token"))
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		When("generating a new client returns an error", func() {
			BeforeEach(func() {
				fakeKubernetesController.NewClientReturns(nil, errors.New("error generating client"))
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		When("generating a new clientset returns an error", func() {
			BeforeEach(func() {
				fakeKubernetesController.NewClientsetReturns(nil, errors.New("error generating clientset"))
			})

			It("continues and returns an empty list", func() {
				Expect(err).To(BeNil())
				Expect(providers).To(HaveLen(0))
			})
		})

		It("succeeds", func() {
			Expect(err).To(BeNil())
			Expect(providers).ToNot(BeNil())
			Expect(providers).To(HaveLen(2))
			Expect(providers[0]).ToNot(BeNil())
			Expect(providers[1]).ToNot(BeNil())
			Expect(providers[0].Name).To(Equal("test-name1"))
			Expect(providers[0].Host).To(Equal("test-host1"))
			Expect(providers[0].CAData).To(Equal("12341234"))
			Expect(providers[0].TokenProvider).To(Equal("test-token-provider1"))
			Expect(providers[0].Client).ToNot(BeNil())
			Expect(providers[0].Clientset).ToNot(BeNil())
			Expect(providers[1].Name).To(Equal("test-name2"))
			Expect(providers[1].Host).To(Equal("test-host2"))
			Expect(providers[1].CAData).To(Equal("56785678"))
			Expect(providers[1].TokenProvider).To(Equal("test-token-provider2"))
			Expect(providers[1].Client).ToNot(BeNil())
			Expect(providers[1].Clientset).ToNot(BeNil())
			config := fakeKubernetesController.NewClientArgsForCall(0)
			Expect(config.Timeout).To(Equal(time.Second * 1))
			config = fakeKubernetesController.NewClientArgsForCall(1)
			Expect(config.Timeout).To(Equal(time.Second * 1))
		})
	})
})
