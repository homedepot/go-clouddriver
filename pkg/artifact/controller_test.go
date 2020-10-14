package artifact_test

import (
	"io/ioutil"
	"log"
	"os"

	. "github.com/billiford/go-clouddriver/pkg/artifact"
	"github.com/billiford/go-clouddriver/pkg/helm"
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
				Expect(artifactCredentials).To(HaveLen(9))
				for _, ac := range artifactCredentials {
					Expect(ac.Repository).To(BeEmpty())
				}
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
})
