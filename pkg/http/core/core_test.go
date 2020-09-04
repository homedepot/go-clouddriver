package core_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/arcade/arcadefakes"
	"github.com/billiford/go-clouddriver/pkg/helm/helmfakes"
	kubefakes "github.com/billiford/go-clouddriver/pkg/http/core/kubernetes/kubernetesfakes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/kubernetesfakes"
	"github.com/billiford/go-clouddriver/pkg/server"
	"github.com/billiford/go-clouddriver/pkg/sql/sqlfakes"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	err                   error
	svr                   *httptest.Server
	uri                   string
	req                   *http.Request
	body                  *bytes.Buffer
	res                   *http.Response
	fakeArcadeClient      *arcadefakes.FakeClient
	fakeHelmClient        *helmfakes.FakeClient
	fakeSQLClient         *sqlfakes.FakeClient
	fakeKubeClient        *kubernetesfakes.FakeClient
	fakeKubeController    *kubernetesfakes.FakeController
	fakeKubeActionHandler *kubefakes.FakeActionHandler
	fakeAction            *kubefakes.FakeAction
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

	fakeKubeClient = &kubernetesfakes.FakeClient{}
	fakeKubeClient.GetReturns(&unstructured.Unstructured{Object: map[string]interface{}{}}, nil)

	fakeKubeController = &kubernetesfakes.FakeController{}
	fakeKubeController.NewClientReturns(fakeKubeClient, nil)

	fakeAction = &kubefakes.FakeAction{}
	fakeKubeActionHandler = &kubefakes.FakeActionHandler{}

	fakeKubeActionHandler.NewDeployManifestActionReturns(fakeAction)
	fakeKubeActionHandler.NewScaleManifestActionReturns(fakeAction)
	fakeKubeActionHandler.NewRollingRestartActionReturns(fakeAction)
	fakeKubeActionHandler.NewRollbackActionReturns(fakeAction)

	fakeArcadeClient = &arcadefakes.FakeClient{}

	fakeHelmClient = &helmfakes.FakeClient{}

	// Disable debug logging.
	gin.SetMode(gin.ReleaseMode)

	// Create new gin instead of using gin.Default().
	// This disables request logging which we don't want for tests.
	r := gin.New()
	r.Use(gin.Recovery())

	c := &server.Config{
		ArcadeClient:      fakeArcadeClient,
		HelmClient:        fakeHelmClient,
		SQLClient:         fakeSQLClient,
		KubeController:    fakeKubeController,
		KubeActionHandler: fakeKubeActionHandler,
	}

	// Create server.
	server.Setup(r, c)
	svr = httptest.NewServer(r)
	body = &bytes.Buffer{}
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

func getClouddriverError() clouddriver.Error {
	ce := clouddriver.Error{}
	b, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(b, &ce)
	return ce
}
