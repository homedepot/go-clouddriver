package artifact_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/v32/github"
	. "github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {
	var (
		cc  CredentialsController
		err error
		dir string
	)

	BeforeEach(func() {
		dir = "test"
		log.SetOutput(ioutil.Discard)
	})

	Describe("#NewCredentialsController", func() {
		JustBeforeEach(func() {
			cc, err = NewCredentialsController(dir)
		})

		When("the directory does not exist", func() {
			BeforeEach(func() {
				dir = "i-dont-exist"
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("open i-dont-exist: no such file or directory"))
			})
		})

		When("a file exists with bad json", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "cred*.json")
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("unexpected end of JSON input"))
			})
		})

		When("a file exists without specifying a credential name", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "cred*.json")
				_, err = tmpFile.WriteString("{}")
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("no \"name\" found in artifact config file test/cred"))
			})
		})

		When("a duplicate credential exists", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "cred*.json")
				_, err = tmpFile.WriteString(`{
					"name": "helm-test"
				}`)
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("duplicate artifact credential listed: helm-test"))
			})
		})

		When("a type gcs/object with keyfile that does not exist", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "gcs-fake.json")
				_, err = tmpFile.WriteString(`{
					"name": "gcs-fake",
					"types": [
					  "gcs/object"
					],
					"jsonPath": "fake.json"
				}`)
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(`dialing: cannot read credentials file: open fake.json: no such file or directory`))
			})
		})

		When("a enterprise github/file artifact does not set the baseURL", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "cred*.json")
				_, err = tmpFile.WriteString(`{
          "enterprise": true,
					"name": "github.example2.com",
					"types": [
					  "github/file"
					]
				}`)
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix(`github file github.example2.com missing required "baseURL" attribute`))
			})
		})

		When("a enterprise github/file artifact does not set the baseURL correctly", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "cred*.json")
				_, err = tmpFile.WriteString(`{
          "baseURL": ":haha",
          "enterprise": true,
					"name": "github.example2.com",
					"types": [
					  "github/file"
					]
				}`)
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix(`parse ":haha": missing protocol scheme`))
			})
		})

		When("a type helm/chart is missing the repository attribute", func() {
			var tmpFile *os.File

			BeforeEach(func() {
				tmpFile, err = ioutil.TempFile("test", "cred*.json")
				_, err = tmpFile.WriteString(`{
					"name": "helm-test2",
					"types": [
					  "helm/chart"
					]
				}`)
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				os.Remove(tmpFile.Name())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix(`helm chart helm-test2 missing required "repository" attribute`))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#ListArtifactCredentialsNamesAndTypes", func() {
		var artifactCredentials []Credentials

		BeforeEach(func() {
			cc, err = NewCredentialsController(dir)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			artifactCredentials = cc.ListArtifactCredentialsNamesAndTypes()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(artifactCredentials).To(HaveLen(14))
				for _, ac := range artifactCredentials {
					Expect(ac.Repository).To(BeEmpty())
					Expect(ac.Token).To(BeEmpty())
					Expect(ac.BaseURL).To(BeEmpty())
				}
			})
		})
	})

	Describe("#GCSClientForAccountName", func() {
		var (
			storageClient *storage.Client
			accountName   string
		)

		BeforeEach(func() {
			accountName = "gcs-spinnaker"
			cc, err = NewCredentialsController(dir)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			storageClient, err = cc.GCSClientForAccountName(accountName)
		})

		When("the account name does not exist in the cache", func() {
			BeforeEach(func() {
				accountName = "fake"
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("gcs account fake not found"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(storageClient).ToNot(BeNil())
			})
		})
	})

	Describe("#GitClientForAccountName", func() {
		var (
			gitClient   *github.Client
			accountName string
		)

		BeforeEach(func() {
			accountName = "github.com"
			cc, err = NewCredentialsController(dir)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			gitClient, err = cc.GitClientForAccountName(accountName)
		})

		When("the account name does not exist in the cache", func() {
			BeforeEach(func() {
				accountName = "fake"
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("git account fake not found"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(gitClient).ToNot(BeNil())
			})
		})
	})

	Describe("#GitRepoClientForAccountName", func() {
		var (
			gitRepoClient *http.Client
			accountName   string
		)

		BeforeEach(func() {
			accountName = "github-spinnaker"
			cc, err = NewCredentialsController(dir)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			gitRepoClient, err = cc.GitRepoClientForAccountName(accountName)
		})

		When("the account name does not exist in the cache", func() {
			BeforeEach(func() {
				accountName = "fake"
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("git/repo account fake not found"))
			})
		})

		When("the account requires token authorization", func() {
			BeforeEach(func() {
				accountName = "ghe-spinnaker"
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(gitRepoClient).ToNot(BeNil())
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(gitRepoClient).ToNot(BeNil())
			})
		})
	})

	Describe("#HelmClientForAccountName", func() {
		var (
			helmClient  helm.Client
			accountName string
		)

		BeforeEach(func() {
			accountName = "helm-test"
			cc, err = NewCredentialsController(dir)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			helmClient, err = cc.HelmClientForAccountName(accountName)
		})

		When("the account name does not exist in the cache", func() {
			BeforeEach(func() {
				accountName = "fake"
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("helm account fake not found"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(helmClient).ToNot(BeNil())
			})
		})
	})

	Describe("#HTTPClientForAccountName", func() {
		var (
			httpClient  *http.Client
			accountName string
		)

		BeforeEach(func() {
			accountName = "http"
			cc, err = NewCredentialsController(dir)
			Expect(err).To(BeNil())
		})

		JustBeforeEach(func() {
			httpClient, err = cc.HTTPClientForAccountName(accountName)
		})

		When("the account name does not exist in the cache", func() {
			BeforeEach(func() {
				accountName = "fake"
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("http account fake not found"))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(httpClient).ToNot(BeNil())
			})
		})
	})
})
