package core_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/homedepot/go-clouddriver/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Artifacts", func() {
	Describe("#ListArtifactCredentials", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/credentials"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadArtifactCredentials)
			})
		})
	})

	Describe("#ListHelmArtifactAccountNames", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/account/helm-stable/names"
			createRequest(http.MethodGet)
			fakeHelmClient.GetIndexReturns(helm.Index{
				Entries: map[string][]helm.Resource{
					"prometheus-operator": {},
					"minecraft":           {},
				},
			}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the helm client returns an error", func() {
			BeforeEach(func() {
				fakeArtifactCredentialsController.HelmClientForAccountNameReturns(nil, errors.New("error getting helm client"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting helm client"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("getting the index returns an error", func() {
			BeforeEach(func() {
				fakeHelmClient.GetIndexReturns(helm.Index{}, errors.New("error getting index"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting index"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListHelmArtifactAccountNames)
			})
		})
	})

	Describe("#ListHelmArtifactAccountVersions", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/account/helm-stable/versions?artifactName=minecraft"
			createRequest(http.MethodGet)
			fakeHelmClient.GetIndexReturns(helm.Index{
				Entries: map[string][]helm.Resource{
					"prometheus-operator": {
						{
							Version: "2.0.0",
						},
						{
							Version: "2.1.0",
						},
					},
					"minecraft": {
						{
							Version: "1.0.0",
						},
						{
							Version: "1.1.0",
						},
					},
				}}, nil)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the helm client returns an error", func() {
			BeforeEach(func() {
				fakeArtifactCredentialsController.HelmClientForAccountNameReturns(nil, errors.New("error getting helm client"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("error getting helm client"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("getting the index returns an error", func() {
			BeforeEach(func() {
				fakeHelmClient.GetIndexReturns(helm.Index{}, errors.New("error getting index"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Internal Server Error"))
				Expect(ce.Message).To(Equal("error getting index"))
				Expect(ce.Status).To(Equal(http.StatusInternalServerError))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadListHelmArtifactAccountVersions)
			})
		})
	})

	Describe("#GetArtifact", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/artifacts/fetch/"
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the request contains bad data", func() {
			BeforeEach(func() {
				setup()
				body.Write([]byte(";[]---"))
				createRequest(http.MethodPut)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Bad Request"))
				Expect(ce.Message).To(Equal("invalid character ';' looking for beginning of value"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the artifact is type helm/chart", func() {
			BeforeEach(func() {
				body.Write([]byte(payloadRequestFetchHelmArtifact))
				createRequest(http.MethodPut)
				fakeHelmClient.GetChartReturns([]byte("some-binary-data"), nil)
			})

			When("getting the helm client returns an error", func() {
				BeforeEach(func() {
					fakeArtifactCredentialsController.HelmClientForAccountNameReturns(nil, errors.New("error getting helm client"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal("error getting helm client"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("getting the chart returns an error", func() {
				BeforeEach(func() {
					fakeHelmClient.GetChartReturns(nil, errors.New("error getting chart"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal("error getting chart"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("some-binary-data")
				})
			})
		})

		Context("when the artifact is type embedded/base64", func() {
			BeforeEach(func() {
				body.Write([]byte(payloadRequestFetchBase64Artifact))
				createRequest(http.MethodPut)
			})

			When("the reference contains bad data", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(payloadRequestFetchBase64ArtifactBadReference))
					createRequest(http.MethodPut)
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal("illegal base64 data at input byte 3"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("it succeeds", func() {
				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("helloworld\n")
				})
			})
		})

		Context("when the artifact is type github/file", func() {
			BeforeEach(func() {
				body.Write([]byte(fmt.Sprintf(payloadRequestFetchGithubFileArtifact, fakeGithubServer.URL())))
				createRequest(http.MethodPut)
			})

			When("getting the client returns an error", func() {
				BeforeEach(func() {
					fakeArtifactCredentialsController.GitClientForAccountNameReturns(nil, errors.New("error getting git client"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal("error getting git client"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("the reference is incorrect", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(fmt.Sprintf(payloadRequestFetchGithubFileArtifact, "https://bad-reference")))
					createRequest(http.MethodPut)
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal(fmt.Sprintf("content URL https://bad-reference/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml should have base URL %s",
						fakeGithubClient.BaseURL.String())))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("creating the github request returns an error", func() {
				BeforeEach(func() {
					fakeGithubClient.BaseURL, err = url.Parse(strings.TrimSuffix(fakeGithubClient.BaseURL.String(), "/"))
					Expect(err).To(BeNil())
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal(fmt.Sprintf(`BaseURL must have a trailing slash, but "%s" does not`, fakeGithubClient.BaseURL.String())))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("the branch is set in the version", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(fmt.Sprintf(payloadRequestFetchGithubFileArtifactTestBranch, fakeGithubServer.URL())))
					createRequest(http.MethodPut)
					fakeGithubServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml", "ref=test"),
						ghttp.RespondWith(http.StatusOK, `{
							"name": "helloweb-deployment.yaml",
							"path": "hello-app/manifests/helloweb-deployment.yaml",
							"sha": "53de8cefaaadead83771a4e7ec0e3024510a8969",
							"size": 1035,
							"url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
							"html_url": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml",
							"git_url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
							"download_url": "https://raw.githubusercontent.com/GoogleCloudPlatform/kubernetes-engine-samples/master/hello-app/manifests/helloweb-deployment.yaml",
							"type": "file",
							"content": "IyBDb3B5cmlnaHQgMjAyMCBHb29nbGUgTExDCiMKIyBMaWNlbnNlZCB1bmRl\nciB0aGUgQXBhY2hlIExpY2Vuc2UsIFZlcnNpb24gMi4wICh0aGUgIkxpY2Vu\nc2UiKTsKIyB5b3UgbWF5IG5vdCB1c2UgdGhpcyBmaWxlIGV4Y2VwdCBpbiBj\nb21wbGlhbmNlIHdpdGggdGhlIExpY2Vuc2UuCiMgWW91IG1heSBvYnRhaW4g\nYSBjb3B5IG9mIHRoZSBMaWNlbnNlIGF0CiMKIyAgICAgaHR0cDovL3d3dy5h\ncGFjaGUub3JnL2xpY2Vuc2VzL0xJQ0VOU0UtMi4wCiMKIyBVbmxlc3MgcmVx\ndWlyZWQgYnkgYXBwbGljYWJsZSBsYXcgb3IgYWdyZWVkIHRvIGluIHdyaXRp\nbmcsIHNvZnR3YXJlCiMgZGlzdHJpYnV0ZWQgdW5kZXIgdGhlIExpY2Vuc2Ug\naXMgZGlzdHJpYnV0ZWQgb24gYW4gIkFTIElTIiBCQVNJUywKIyBXSVRIT1VU\nIFdBUlJBTlRJRVMgT1IgQ09ORElUSU9OUyBPRiBBTlkgS0lORCwgZWl0aGVy\nIGV4cHJlc3Mgb3IgaW1wbGllZC4KIyBTZWUgdGhlIExpY2Vuc2UgZm9yIHRo\nZSBzcGVjaWZpYyBsYW5ndWFnZSBnb3Zlcm5pbmcgcGVybWlzc2lvbnMgYW5k\nCiMgbGltaXRhdGlvbnMgdW5kZXIgdGhlIExpY2Vuc2UuCgojIFtTVEFSVCBj\nb250YWluZXJfaGVsbG9hcHBfZGVwbG95bWVudF0KYXBpVmVyc2lvbjogYXBw\ncy92MQpraW5kOiBEZXBsb3ltZW50Cm1ldGFkYXRhOgogIG5hbWU6IGhlbGxv\nd2ViCiAgbGFiZWxzOgogICAgYXBwOiBoZWxsbwpzcGVjOgogIHNlbGVjdG9y\nOgogICAgbWF0Y2hMYWJlbHM6CiAgICAgIGFwcDogaGVsbG8KICAgICAgdGll\ncjogd2ViCiAgdGVtcGxhdGU6CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxz\nOgogICAgICAgIGFwcDogaGVsbG8KICAgICAgICB0aWVyOiB3ZWIKICAgIHNw\nZWM6CiAgICAgIGNvbnRhaW5lcnM6CiAgICAgIC0gbmFtZTogaGVsbG8tYXBw\nCiAgICAgICAgaW1hZ2U6IGdjci5pby9nb29nbGUtc2FtcGxlcy9oZWxsby1h\ncHA6MS4wCiAgICAgICAgcG9ydHM6CiAgICAgICAgLSBjb250YWluZXJQb3J0\nOiA4MDgwCiMgW0VORCBjb250YWluZXJfaGVsbG9hcHBfZGVwbG95bWVudF0K\n",
							"encoding": "base64",
							"_links": {
								"self": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
								"git": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
								"html": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml"
							}
						}`),
					))
				})

				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("# Copyright 2020 Google LLC\n#\n# Licensed under the Apache License, Version 2.0 (the \"License\");\n# you may not use this file except in compliance with the License.\n# You may obtain a copy of the License at\n#\n#     http://www.apache.org/licenses/LICENSE-2.0\n#\n# Unless required by applicable law or agreed to in writing, software\n# distributed under the License is distributed on an \"AS IS\" BASIS,\n# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n# See the License for the specific language governing permissions and\n# limitations under the License.\n\n# [START container_helloapp_deployment]\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: helloweb\n  labels:\n    app: hello\nspec:\n  selector:\n    matchLabels:\n      app: hello\n      tier: web\n  template:\n    metadata:\n      labels:\n        app: hello\n        tier: web\n    spec:\n      containers:\n      - name: hello-app\n        image: gcr.io/google-samples/hello-app:1.0\n        ports:\n        - containerPort: 8080\n# [END container_helloapp_deployment]\n")
				})
			})

			When("the server is not reachable", func() {
				var url, addr string
				BeforeEach(func() {
					url = fakeGithubServer.URL()
					addr = fakeGithubServer.Addr()
					fakeGithubServer.Close()
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal(fmt.Sprintf(`Get "%s/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master": dial tcp %s: connect: connection refused`, url, addr)))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("the base64 encoded content is bad", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(fmt.Sprintf(payloadRequestFetchGithubFileArtifactTestBranch, fakeGithubServer.URL())))
					createRequest(http.MethodPut)
					fakeGithubServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml", "ref=test"),
						ghttp.RespondWith(http.StatusOK, `{
							"name": "helloweb-deployment.yaml",
							"path": "hello-app/manifests/helloweb-deployment.yaml",
							"sha": "53de8cefaaadead83771a4e7ec0e3024510a8969",
							"size": 1035,
							"url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
							"html_url": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml",
							"git_url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
							"download_url": "https://raw.githubusercontent.com/GoogleCloudPlatform/kubernetes-engine-samples/master/hello-app/manifests/helloweb-deployment.yaml",
							"type": "file",
							"content": "{}",
							"encoding": "base64",
							"_links": {
								"self": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
								"git": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
								"html": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml"
							}
						}`),
					))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal("illegal base64 data at input byte 0"))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("the content is not encoded", func() {
				BeforeEach(func() {
					fakeGithubServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml", "ref=master"),
						ghttp.RespondWith(http.StatusOK, `{
							"name": "helloweb-deployment.yaml",
							"path": "hello-app/manifests/helloweb-deployment.yaml",
							"sha": "53de8cefaaadead83771a4e7ec0e3024510a8969",
							"size": 1035,
							"url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
							"html_url": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml",
							"git_url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
							"download_url": "https://raw.githubusercontent.com/GoogleCloudPlatform/kubernetes-engine-samples/master/hello-app/manifests/helloweb-deployment.yaml",
							"type": "file",
							"content": "helloworld",
							"_links": {
								"self": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
								"git": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
								"html": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml"
							}
						}`),
					))
				})

				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("helloworld")
				})
			})

			When("it succeeds", func() {
				BeforeEach(func() {
					fakeGithubServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml", "ref=master"),
						ghttp.RespondWith(http.StatusOK, `{
							"name": "helloweb-deployment.yaml",
							"path": "hello-app/manifests/helloweb-deployment.yaml",
							"sha": "53de8cefaaadead83771a4e7ec0e3024510a8969",
							"size": 1035,
							"url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
							"html_url": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml",
							"git_url": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
							"download_url": "https://raw.githubusercontent.com/GoogleCloudPlatform/kubernetes-engine-samples/master/hello-app/manifests/helloweb-deployment.yaml",
							"type": "file",
							"content": "IyBDb3B5cmlnaHQgMjAyMCBHb29nbGUgTExDCiMKIyBMaWNlbnNlZCB1bmRl\nciB0aGUgQXBhY2hlIExpY2Vuc2UsIFZlcnNpb24gMi4wICh0aGUgIkxpY2Vu\nc2UiKTsKIyB5b3UgbWF5IG5vdCB1c2UgdGhpcyBmaWxlIGV4Y2VwdCBpbiBj\nb21wbGlhbmNlIHdpdGggdGhlIExpY2Vuc2UuCiMgWW91IG1heSBvYnRhaW4g\nYSBjb3B5IG9mIHRoZSBMaWNlbnNlIGF0CiMKIyAgICAgaHR0cDovL3d3dy5h\ncGFjaGUub3JnL2xpY2Vuc2VzL0xJQ0VOU0UtMi4wCiMKIyBVbmxlc3MgcmVx\ndWlyZWQgYnkgYXBwbGljYWJsZSBsYXcgb3IgYWdyZWVkIHRvIGluIHdyaXRp\nbmcsIHNvZnR3YXJlCiMgZGlzdHJpYnV0ZWQgdW5kZXIgdGhlIExpY2Vuc2Ug\naXMgZGlzdHJpYnV0ZWQgb24gYW4gIkFTIElTIiBCQVNJUywKIyBXSVRIT1VU\nIFdBUlJBTlRJRVMgT1IgQ09ORElUSU9OUyBPRiBBTlkgS0lORCwgZWl0aGVy\nIGV4cHJlc3Mgb3IgaW1wbGllZC4KIyBTZWUgdGhlIExpY2Vuc2UgZm9yIHRo\nZSBzcGVjaWZpYyBsYW5ndWFnZSBnb3Zlcm5pbmcgcGVybWlzc2lvbnMgYW5k\nCiMgbGltaXRhdGlvbnMgdW5kZXIgdGhlIExpY2Vuc2UuCgojIFtTVEFSVCBj\nb250YWluZXJfaGVsbG9hcHBfZGVwbG95bWVudF0KYXBpVmVyc2lvbjogYXBw\ncy92MQpraW5kOiBEZXBsb3ltZW50Cm1ldGFkYXRhOgogIG5hbWU6IGhlbGxv\nd2ViCiAgbGFiZWxzOgogICAgYXBwOiBoZWxsbwpzcGVjOgogIHNlbGVjdG9y\nOgogICAgbWF0Y2hMYWJlbHM6CiAgICAgIGFwcDogaGVsbG8KICAgICAgdGll\ncjogd2ViCiAgdGVtcGxhdGU6CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxz\nOgogICAgICAgIGFwcDogaGVsbG8KICAgICAgICB0aWVyOiB3ZWIKICAgIHNw\nZWM6CiAgICAgIGNvbnRhaW5lcnM6CiAgICAgIC0gbmFtZTogaGVsbG8tYXBw\nCiAgICAgICAgaW1hZ2U6IGdjci5pby9nb29nbGUtc2FtcGxlcy9oZWxsby1h\ncHA6MS4wCiAgICAgICAgcG9ydHM6CiAgICAgICAgLSBjb250YWluZXJQb3J0\nOiA4MDgwCiMgW0VORCBjb250YWluZXJfaGVsbG9hcHBfZGVwbG95bWVudF0K\n",
							"encoding": "base64",
							"_links": {
								"self": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml?ref=master",
								"git": "https://api.github.com/repos/GoogleCloudPlatform/kubernetes-engine-samples/git/blobs/53de8cefaaadead83771a4e7ec0e3024510a8969",
								"html": "https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/blob/master/hello-app/manifests/helloweb-deployment.yaml"
							}
						}`),
					))
				})

				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("# Copyright 2020 Google LLC\n#\n# Licensed under the Apache License, Version 2.0 (the \"License\");\n# you may not use this file except in compliance with the License.\n# You may obtain a copy of the License at\n#\n#     http://www.apache.org/licenses/LICENSE-2.0\n#\n# Unless required by applicable law or agreed to in writing, software\n# distributed under the License is distributed on an \"AS IS\" BASIS,\n# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n# See the License for the specific language governing permissions and\n# limitations under the License.\n\n# [START container_helloapp_deployment]\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: helloweb\n  labels:\n    app: hello\nspec:\n  selector:\n    matchLabels:\n      app: hello\n      tier: web\n  template:\n    metadata:\n      labels:\n        app: hello\n        tier: web\n    spec:\n      containers:\n      - name: hello-app\n        image: gcr.io/google-samples/hello-app:1.0\n        ports:\n        - containerPort: 8080\n# [END container_helloapp_deployment]\n")
				})
			})
		})

		Context("when the artifact is type git/repo", func() {
			BeforeEach(func() {
				body.Write([]byte(fmt.Sprintf(payloadRequestFetchGitRepoArtifact, fakeFileServer.URL())))
				createRequest(http.MethodPut)
			})

			When("getting the client returns an error", func() {
				BeforeEach(func() {
					fakeArtifactCredentialsController.GitRepoClientForAccountNameReturns(nil, errors.New("error getting http client"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal("error getting http client"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("the server is not reachable", func() {
				var url, addr string

				BeforeEach(func() {
					url = fakeFileServer.URL()
					addr = fakeFileServer.Addr()
					fakeFileServer.Close()
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal(fmt.Sprintf(`Get "%s/git-repo/archive/master.tar.gz": dial tcp %s: connect: connection refused`, url, addr)))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("the repo is not readable", func() {
				var url string

				BeforeEach(func() {
					url = fakeFileServer.URL()
					fakeFileServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/git-repo/archive/master.tar.gz"),
						ghttp.RespondWith(http.StatusNotFound, ""),
					))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal(fmt.Sprintf(`error getting git/repo (repo: %s/git-repo, branch: master): 404 Not Found`, url)))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("the branch is set in the version", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(fmt.Sprintf(payloadRequestFetchGitRepoArtifactBranch, fakeFileServer.URL())))
					createRequest(http.MethodPut)

					actual, err := ioutil.ReadFile("test/git-repo-test.tar.gz")
					Expect(err).To(BeNil())
					fakeFileServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/git-repo/archive/test.tar.gz"),
						ghttp.RespondWith(http.StatusOK, actual),
					))
				})

				It("succeeds", func() {
					expected, err := ioutil.ReadFile("test/expected-git-repo-test.tar.gz")
					Expect(err).To(BeNil())

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateGZipResponse(expected)
				})
			})

			When("the subpath is set in the location", func() {
				BeforeEach(func() {
					body = &bytes.Buffer{}
					body.Write([]byte(fmt.Sprintf(payloadRequestFetchGitRepoArtifactSubPath, fakeFileServer.URL())))
					createRequest(http.MethodPut)

					actual, err := ioutil.ReadFile("test/git-repo-master.tar.gz")
					Expect(err).To(BeNil())
					fakeFileServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/git-repo/archive/master.tar.gz"),
						ghttp.RespondWith(http.StatusOK, actual),
					))
				})

				It("succeeds", func() {
					expected, err := ioutil.ReadFile("test/expected-git-repo-master-subpath.tar.gz")
					Expect(err).To(BeNil())

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateGZipResponse(expected)
				})
			})

			When("it succeeds", func() {
				BeforeEach(func() {
					actual, err := ioutil.ReadFile("test/git-repo-master.tar.gz")
					Expect(err).To(BeNil())
					fakeFileServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/git-repo/archive/master.tar.gz"),
						ghttp.RespondWith(http.StatusOK, actual),
					))
				})

				It("succeeds", func() {
					expected, err := ioutil.ReadFile("test/expected-git-repo-master.tar.gz")
					Expect(err).To(BeNil())

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateGZipResponse(expected)
				})
			})
		})

		Context("when the artifact is type http/file", func() {
			BeforeEach(func() {
				body.Write([]byte(fmt.Sprintf(payloadRequestFetchHTTPFileArtifact, fakeFileServer.URL())))
				createRequest(http.MethodPut)
			})

			When("getting the client returns an error", func() {
				BeforeEach(func() {
					fakeArtifactCredentialsController.HTTPClientForAccountNameReturns(nil, errors.New("error getting http client"))
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Bad Request"))
					Expect(ce.Message).To(Equal("error getting http client"))
					Expect(ce.Status).To(Equal(http.StatusBadRequest))
				})
			})

			When("the server is not reachable", func() {
				var url, addr string

				BeforeEach(func() {
					url = fakeFileServer.URL()
					addr = fakeFileServer.Addr()
					fakeFileServer.Close()
				})

				It("returns an error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					ce := getClouddriverError()
					Expect(ce.Error).To(HavePrefix("Internal Server Error"))
					Expect(ce.Message).To(Equal(fmt.Sprintf(`Get "%s/hello": dial tcp %s: connect: connection refused`, url, addr)))
					Expect(ce.Status).To(Equal(http.StatusInternalServerError))
				})
			})

			When("it succeeds", func() {
				BeforeEach(func() {
					fakeFileServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/hello"),
						ghttp.RespondWith(http.StatusOK, `world`),
					))
				})

				It("succeeds", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					validateTextResponse("world")
				})
			})
		})

		Context("when the artifact is not an implemented type", func() {
			BeforeEach(func() {
				body.Write([]byte(payloadRequestFetchNotImplementedArtifact))
				createRequest(http.MethodPut)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotImplemented))
				ce := getClouddriverError()
				Expect(ce.Error).To(HavePrefix("Not Implemented"))
				Expect(ce.Message).To(Equal("getting artifact of type unknown/type not implemented"))
				Expect(ce.Status).To(Equal(http.StatusNotImplemented))
			})
		})
	})
})
