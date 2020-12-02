package v1_test

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/pkg/server"
	"github.com/homedepot/go-clouddriver/pkg/sql/sqlfakes"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	err           error
	svr           *httptest.Server
	uri           string
	req           *http.Request
	body          *bytes.Buffer
	res           *http.Response
	fakeSQLClient *sqlfakes.FakeClient
)

func setup() {
	// Setup fake SQL client.
	fakeSQLClient = &sqlfakes.FakeClient{}

	// Disable debug logging.
	gin.SetMode(gin.ReleaseMode)

	// Create new gin instead of using gin.Default().
	// This disables request logging which we don't want for tests.
	r := gin.New()
	r.Use(gin.Recovery())

	c := &server.Config{
		SQLClient: fakeSQLClient,
	}
	// Create server.
	server.Setup(r, c)
	svr = httptest.NewServer(r)
	body = &bytes.Buffer{}
}

func teardown() {
	svr.Close()

	mt, mtp, _ := mime.ParseMediaType(res.Header.Get("content-type"))
	Expect(mt).To(Equal("application/json"), "content-type")
	Expect(mtp["charset"]).To(Equal("utf-8"), "charset")
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
	actual, _ := ioutil.ReadAll(res.Body)
	Expect(actual).To(MatchJSON(expected), "correct body")
}
