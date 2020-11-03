package helm_test

import (
	"net/http"

	"github.com/billiford/go-clouddriver/pkg/helm"
	. "github.com/billiford/go-clouddriver/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server
		client Client
		err    error
		index  helm.Index
		b      []byte
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("#GetIndex", func() {
		JustBeforeEach(func() {
			index, err = client.GetIndex()
		})

		When("the uri is invalid", func() {
			BeforeEach(func() {
				client = NewClient("::haha")
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("the server is not reachable", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("the response is not 2XX", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error getting helm index: 500 Internal Server Error"))
			})
		})

		When("the response is status not modified", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.RespondWith(http.StatusNotModified, nil),
				)
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})

		When("the server returns bad data", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.RespondWith(http.StatusOK, ";{["),
				)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `;{[` into helm.Index"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/index.yaml"),
					ghttp.RespondWith(http.StatusOK, `apiVersion: v1
entries:
  hello-app:
  - apiVersion: v2
    appVersion: "1.0"
    created: 2020-03-20T13:13:42.956Z
    description: A Helm chart for deploying hello-app to Kuberntes
    digest: a934a39ba7f77e8e16b609d01c493d7ead5ec24c78d2918d3a00ebaa3c3fdfd2
    home: https://pages.github.com/test
    maintainers:
    - email: me@test.com
      name: CD Team
    name: hello-app
    urls:
    - https://helm..com/artifactory/helm/test-1.0.0.tgz
    version: 1.0.0
  thd-cd-hello-app:
  - apiVersion: v2
    appVersion: "2.0"
    created: 2020-08-14T22:11:19.894Z
    description: A Helm chart for deploying hello-app to Kubernetes
    digest: c9c2663958ffdf2c790ebfd0f502176c3ac663014e9d6ce49955c4178c2f43d7
    home: https://pages.github.com/test
    maintainers:
    - email: me@test.com
      name: CD Team
    name: thd-cd-hello-app
    urls:
    - https://helm..com/artifactory/helm/test-3.0.2.tgz
    version: 3.0.2`),
				))
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(index.Entries).To(HaveLen(2))
			})
		})
	})

	Describe("#GetChart", func() {
		JustBeforeEach(func() {
			b, err = client.GetChart("test-name", "0.1.0")
		})

		When("the uri is invalid", func() {
			BeforeEach(func() {
				client = NewClient("::haha")
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("the server is not reachable", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("the response is not 2XX", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				)
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error getting helm chart: 500 Internal Server Error"))
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/test-name-0.1.0.tgz"),
					ghttp.RespondWith(http.StatusOK, `some-binary-data`),
				))
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("some-binary-data"))
			})
		})
	})
})
