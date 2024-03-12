package v1_test

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/arcade/pkg/arcadefakes"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/api"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/kubernetesfakes"
	"github.com/homedepot/go-clouddriver/internal/sql/sqlfakes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/onsi/gomega"
)

var (
	err                error
	svr                *httptest.Server
	uri                string
	req                *http.Request
	body               *bytes.Buffer
	res                *http.Response
	fakeSQLClient      *sqlfakes.FakeClient
	fakeArcadeClient   *arcadefakes.FakeClient
	fakeKubeClient     *kubernetesfakes.FakeClient
	fakeKubeController *kubernetesfakes.FakeController
	fakeDaemonSets     *unstructured.UnstructuredList
	fakeDeployments    *unstructured.UnstructuredList
	fakeIngresses      *unstructured.UnstructuredList
	fakePods           *unstructured.UnstructuredList
	fakeReplicaSets    *unstructured.UnstructuredList
	fakeServices       *unstructured.UnstructuredList
	fakeStatefulSets   *unstructured.UnstructuredList
)

func setup() {
	// Setup fake clients.
	fakeSQLClient = &sqlfakes.FakeClient{}
	fakeArcadeClient = &arcadefakes.FakeClient{}
	fakeKubeClient = &kubernetesfakes.FakeClient{}

	fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
		Name:   "test-account",
		Host:   "http://localhost",
		CAData: "",
	}, nil)

	fakeKubeController = &kubernetesfakes.FakeController{}
	fakeKubeController.NewClientReturns(fakeKubeClient, nil)

	// Disable debug logging.
	gin.SetMode(gin.ReleaseMode)

	// Create new gin instead of using gin.Default().
	// This disables request logging which we don't want for tests.
	r := gin.New()
	r.Use(gin.Recovery())

	c := &internal.Controller{
		SQLClient:            fakeSQLClient,
		ArcadeClient:         fakeArcadeClient,
		KubernetesController: fakeKubeController,
	}
	// Create server.
	server := api.NewServer(r)
	server.WithController(c)
	server.Setup()

	svr = httptest.NewServer(r)
	body = &bytes.Buffer{}

	fakeDaemonSets = &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"kind":       "DaemonSet",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"name":      "test-daemonset",
						"namespace": "test-namespace",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application",
						},
					},
				},
			},
		},
	}
	fakeDeployments = &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"name":      "test-deployment1",
						"namespace": "test-namespace1",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application1",
						},
					},
				},
			},
			{
				Object: map[string]interface{}{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"name":      "test-deployment2",
						"namespace": "test-namespace2",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application2",
						},
					},
				},
			},
		},
	}
	fakeIngresses = &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"kind":       "Ingress",
					"apiVersion": "apps/v1",
					"metadata": map[string]interface{}{
						"name":      "test-ingress",
						"namespace": "test-namespace",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application",
						},
					},
				},
			},
		},
	}
	fakeReplicaSets = &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"kind":       "ReplicaSet",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name":      "test-replicaset-v001",
						"namespace": "test-namespace",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application",
						},
					},
				},
			},
		},
	}
	fakeStatefulSets = &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"kind":       "StatefulSet",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name":      "test-statefulset",
						"namespace": "test-namespace",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application",
						},
					},
				},
			},
		},
	}
	fakeServices = &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"kind":       "Service",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name":      "test-service",
						"namespace": "test-namespace",
						"annotations": map[string]interface{}{
							"moniker.spinnaker.io/application": "test-application",
						},
					},
				},
			},
		},
	}
}

func teardown() {
	svr.Close()

	mt, mtp, _ := mime.ParseMediaType(res.Header.Get("content-type"))
	Expect(mt).To(Equal("application/json"), "content-type")
	Expect(mtp["charset"]).To(Equal("utf-8"), "charset")
	res.Body.Close()
}

func createRequest(method string) {
	req, _ = http.NewRequest(method, uri, io.NopCloser(body))
	req.Header.Set("API-Key", "test")
}

func doRequest() {
	res, err = http.DefaultClient.Do(req)
}

func validateResponse(expected string) {
	actual, _ := io.ReadAll(res.Body)
	Expect(actual).To(MatchJSON(expected), "correct body")
}
