package internal_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/arcade/arcadefakes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/kubernetesfakes"
	"github.com/homedepot/go-clouddriver/internal/sql/sqlfakes"
)

var _ = Describe("Controller", func() {
	var (
		c                        *internal.Controller
		fakeSQLClient            *sqlfakes.FakeClient
		fakeArcadeClient         *arcadefakes.FakeClient
		fakeKubernetesController *kubernetesfakes.FakeController
		fakeKubernetesClient     *kubernetesfakes.FakeClient
		fakeKubernetesClientset  *kubernetesfakes.FakeClientset
		provider                 *kubernetes.Provider
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

		fakeKubernetesController.NewClientReturns(fakeKubernetesClient, nil)
		fakeKubernetesController.NewClientsetReturns(fakeKubernetesClientset, nil)

		c = &internal.Controller{
			ArcadeClient:         fakeArcadeClient,
			KubernetesController: fakeKubernetesController,
			SQLClient:            fakeSQLClient,
		}
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
})
