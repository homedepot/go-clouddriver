package kubernetes_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scale", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		kubernetesController.Scale(c, scaleManifestRequest)
	})

	When("getting the provider returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("internal: error getting kubernetes provider spin-cluster-account: error getting provider"))
		})
	})

	When("getting the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturns(nil, errors.New("error getting manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error getting manifest"))
		})
	})

	When("converting the replicas returns an error", func() {
		BeforeEach(func() {
			scaleManifestRequest.Replicas = "asdf"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("strconv.Atoi: parsing \"asdf\": invalid syntax"))
		})
	})

	When("applying the manifest returns an error", func() {
		BeforeEach(func() {
			fakeKubeClient.ApplyReturns(kubernetes.Metadata{}, errors.New("error applying manifest"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error applying manifest"))
		})
	})

	When("creating the resource returns an error", func() {
		BeforeEach(func() {
			fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(c.Errors.Last().Error()).To(Equal("error creating resource"))
		})
	})

	When("the kind is not supported to scale", func() {
		BeforeEach(func() {
			scaleManifestRequest.ManifestName = "not-supported-kind test-name"
		})

		It("returns an error", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
			Expect(c.Errors.Last().Error()).To(Equal("scaling kind not-supported-kind not currently supported"))
		})
	})

	When("The kind is ReplicaSet", func() {
		BeforeEach(func() {
			scaleManifestRequest.ManifestName = "replicaset someReplicaSet"
		})

		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
		})
	})

	When("The kind is StatefulSet", func() {
		BeforeEach(func() {
			scaleManifestRequest.ManifestName = "statefulset someStatefulSet"
		})

		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
		})
	})

	When("it succeeds", func() {
		It("succeeds", func() {
			Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			_, _, namespace := fakeKubeClient.GetArgsForCall(0)
			Expect(namespace).To(Equal(""))
			u := fakeKubeClient.ApplyArgsForCall(0)
			b, _ := json.Marshal(&u)
			Expect(string(b)).To(Equal("{\"spec\":{\"replicas\":16}}"))
		})
	})

	When("Using a namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(namespaceScopedProvider, nil)
		})

		When("the kind is not supported", func() {
			BeforeEach(func() {
				scaleManifestRequest.ManifestName = "namespace someNamespace"
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'namespace'"))
			})
		})

		When("the kind is supported", func() {
			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				_, _, namespace := fakeKubeClient.GetArgsForCall(0)
				Expect(namespace).To(Equal("provider-namespace"))
				u := fakeKubeClient.ApplyArgsForCall(0)
				b, _ := json.Marshal(&u)
				Expect(string(b)).To(Equal("{\"spec\":{\"replicas\":16}}"))
			})
		})
	})

	When("Using a multiple namespace-scoped provider", func() {
		BeforeEach(func() {
			fakeSQLClient.GetKubernetesProviderReturns(multipleNamespaceScopedProvider, nil)
		})

		When("the kind is not supported", func() {
			BeforeEach(func() {
				scaleManifestRequest.ManifestName = "namespace someNamespace"
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusBadRequest))
				Expect(c.Errors.Last().Error()).To(Equal("namespace-scoped account not allowed to access cluster-scoped kind: 'namespace'"))
			})
		})

		When("the kind is supported", func() {
			BeforeEach(func() {
				scaleManifestRequest.Location = "provider-namespace"
			})
			It("succeeds", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				_, _, namespace := fakeKubeClient.GetArgsForCall(0)
				Expect(namespace).To(Equal("provider-namespace"))
				u := fakeKubeClient.ApplyArgsForCall(0)
				b, _ := json.Marshal(&u)
				Expect(string(b)).To(Equal("{\"spec\":{\"replicas\":16}}"))
			})
		})
	})
})
