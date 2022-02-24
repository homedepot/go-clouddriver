package core_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/http/httptest"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v32/github"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/api"
	"github.com/homedepot/go-clouddriver/internal/arcade/arcadefakes"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/artifact/artifactfakes"
	"github.com/homedepot/go-clouddriver/internal/fiat/fiatfakes"
	"github.com/homedepot/go-clouddriver/internal/front50/front50fakes"
	"github.com/homedepot/go-clouddriver/internal/helm/helmfakes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/kubernetesfakes"
	"github.com/homedepot/go-clouddriver/internal/sql/sqlfakes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var (
	err                               error
	svr                               *httptest.Server
	uri                               string
	req                               *http.Request
	body                              *bytes.Buffer
	res                               *http.Response
	fakeArcadeClient                  *arcadefakes.FakeClient
	fakeArtifactCredentialsController *artifactfakes.FakeCredentialsController
	fakeFiatClient                    *fiatfakes.FakeClient
	fakeFront50Client                 *front50fakes.FakeClient
	fakeStorageClient                 *storage.Client
	fakeGithubClient                  *github.Client
	fakeHelmClient                    *helmfakes.FakeClient
	fakeSQLClient                     *sqlfakes.FakeClient
	fakeKubeClient                    *kubernetesfakes.FakeClient
	fakeKubeClientset                 *kubernetesfakes.FakeClientset
	fakeKubeController                *kubernetesfakes.FakeController
	fakeStorageServer                 *fakestorage.Server
	fakeGithubServer                  *ghttp.Server
	fakeFileServer                    *ghttp.Server
)

func setup() {
	// Setup fakes.
	fakeSQLClient = &sqlfakes.FakeClient{}
	fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
		Name:   "test-account",
		Host:   "http://localhost",
		CAData: "",
	}, nil)
	fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
		{
			AccountName: "test-account-name",
		},
	}, nil)
	fakeSQLClient.ListKubernetesProvidersReturns([]kubernetes.Provider{
		{
			Name:   "account1",
			Host:   "http://localhost",
			CAData: "",
		},
		{
			Name:   "account2",
			Host:   "http://localhost",
			CAData: "",
		},
		{
			Name:   "account3",
			Host:   "http://localhost",
			CAData: "",
		},
		{
			Name:   "account4",
			Host:   "http://localhost",
			CAData: "",
		},
	}, nil)

	fakeKubeClient = &kubernetesfakes.FakeClient{}
	fakeKubeClient.GetReturns(&unstructured.Unstructured{Object: map[string]interface{}{}}, nil)
	fakeKubeClient.ListByGVRReturns(&unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
							kubernetes.AnnotationSpinnakerMonikerApplication: "wrong-application",
						},
						"name":      "test-name",
						"namespace": "test-namespace",
					},
				},
			},
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
						},
						"name":      "test-name",
						"namespace": "test-namespace",
					},
				},
			},
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
						},
						"name":      "test-name",
						"namespace": "test-namespace",
					},
				},
			},
		},
	}, nil)
	fakeKubeClient.ListByGVRWithContextReturns(&unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster: "deployment test-deployment",
						},
					},
				},
			},
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster: "deployment test-deployment",
						},
					},
				},
			},
		},
	}, nil)
	fakeKubeClient.ListResourceReturns(&unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster: "deployment test-deployment",
						},
					},
				},
			},
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster: "deployment test-deployment",
						},
					},
				},
			},
		},
	}, nil)
	fakeKubeClient.ListResourcesByKindAndNamespaceReturns(&unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
						},
						"name":      "test-name",
						"namespace": "test-namespace",
					},
				},
			},
		},
	}, nil)
	fakeKubeClient.ListResourcesByKindAndNamespaceWithContextReturns(&unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerMonikerCluster:     "deployment test-deployment",
							kubernetes.AnnotationSpinnakerMonikerApplication: "test-application",
						},
						"name":      "test-name",
						"namespace": "test-namespace",
					},
				},
			},
		},
	}, nil)

	fakeKubeClientset = &kubernetesfakes.FakeClientset{}
	fakeKubeClientset.PodLogsReturns("log output", nil)

	fakeKubeController = &kubernetesfakes.FakeController{}
	fakeKubeController.NewClientReturns(fakeKubeClient, nil)
	fakeKubeController.NewClientsetReturns(fakeKubeClientset, nil)

	fakeArcadeClient = &arcadefakes.FakeClient{}
	fakeFiatClient = &fiatfakes.FakeClient{}
	fakeFront50Client = &front50fakes.FakeClient{}

	fakeHelmClient = &helmfakes.FakeClient{}

	fakeStorageServer = fakestorage.NewServer([]fakestorage.Object{
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "fake-bucket",
				Name:       "fake-path/fake-file.txt",
			},
			Content: []byte("fake contents"),
		},
	})
	defer fakeStorageServer.Stop()
	fakeStorageClient = fakeStorageServer.Client()

	fakeGithubServer = ghttp.NewServer()
	fakeFileServer = ghttp.NewServer()

	fakeGithubClient, err = github.NewEnterpriseClient(fakeGithubServer.URL(), fakeGithubServer.URL(), nil)
	Expect(err).To(BeNil())

	fakeArtifactCredentialsController = &artifactfakes.FakeCredentialsController{}
	fakeArtifactCredentialsController.ListArtifactCredentialsNamesAndTypesReturns([]artifact.Credentials{
		{
			Name: "helm-stable",
			Types: []artifact.Type{
				artifact.TypeHelmChart,
			},
		},
		{
			Name: "embedded-artifact",
			Types: []artifact.Type{
				artifact.TypeEmbeddedBase64,
			},
		},
	})
	fakeArtifactCredentialsController.GitClientForAccountNameReturns(fakeGithubClient, nil)
	fakeArtifactCredentialsController.GitRepoClientForAccountNameReturns(http.DefaultClient, nil)
	fakeArtifactCredentialsController.HelmClientForAccountNameReturns(fakeHelmClient, nil)
	fakeArtifactCredentialsController.HTTPClientForAccountNameReturns(http.DefaultClient, nil)
	fakeArtifactCredentialsController.GCSClientForAccountNameReturns(fakeStorageClient, nil)

	// Disable debug logging.
	gin.SetMode(gin.ReleaseMode)

	// Create new gin instead of using gin.Default().
	// This disables request logging which we don't want for tests.
	r := gin.New()
	r.Use(gin.Recovery())
	// Sometimes when pipelines get canceled (or for unknown reasons)
	// c.Errors is not nil, even though it contains no errors. This
	// validates that we return a task during these circumstances.
	r.Use(setContextErrors())

	c := &internal.Controller{
		ArcadeClient:                  fakeArcadeClient,
		ArtifactCredentialsController: fakeArtifactCredentialsController,
		FiatClient:                    fakeFiatClient,
		Front50Client:                 fakeFront50Client,
		SQLClient:                     fakeSQLClient,
		KubernetesController:          fakeKubeController,
	}

	// Create server.
	server := api.NewServer(r)
	server.WithController(c)
	server.Setup()

	svr = httptest.NewServer(r)
	body = &bytes.Buffer{}
	// Ignore log output.
	log.SetOutput(ioutil.Discard)
}

func teardown() {
	svr.Close()
	res.Body.Close()
}

func createRequest(method string) {
	req, _ = http.NewRequest(method, uri, ioutil.NopCloser(body))
	req.Header.Set("API-Key", "test")
}

func doRequest() {
	res, err = http.DefaultClient.Do(req)
}

func validateResponse(expected string) {
	mt, mtp, _ := mime.ParseMediaType(res.Header.Get("content-type"))
	Expect(mt).To(Equal("application/json"), "content-type")
	Expect(mtp["charset"]).To(Equal("utf-8"), "charset")

	actual, _ := ioutil.ReadAll(res.Body)
	Expect(actual).To(MatchJSON(expected), "correct body")
}

func validateTextResponse(expected string) {
	mt, mtp, _ := mime.ParseMediaType(res.Header.Get("content-type"))
	Expect(mt).To(Equal("text/plain"), "content-type")
	Expect(mtp["charset"]).To(Equal("utf-8"), "charset")

	actual, _ := ioutil.ReadAll(res.Body)
	Expect(string(actual)).To(Equal(expected), "correct body")
}

func validateGZipResponse(expected []byte) {
	mt, _, _ := mime.ParseMediaType(res.Header.Get("content-type"))
	Expect(mt).To(Equal("application/x-gzip"), "content-type")

	actual, _ := ioutil.ReadAll(res.Body)
	Expect(actual).To(Equal(expected), "correct body")
}

func getClouddriverError() clouddriver.ErrorResponse {
	ce := clouddriver.ErrorResponse{}
	b, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(b, &ce)

	return ce
}

func setContextErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Errors = []*gin.Error{}
		c.Next()
	}
}
