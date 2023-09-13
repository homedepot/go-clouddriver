package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/homedepot/go-clouddriver/internal/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vary", func() {
	var (
		c   *gin.Context
		r   *http.Request
		err error
	)

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)
		c, _ = gin.CreateTestContext(httptest.NewRecorder())
		r, err = http.NewRequest(http.MethodGet, "", nil)
		Expect(err).To(BeNil())
		c.Request = r
	})

	JustBeforeEach(func() {
		mw := Vary("header1", "header2")
		mw(c)
	})

	It("sets the Vary header", func() {
		Expect(c.Writer.Header().Get("Vary")).To(Equal("header1, header2"))
	})
})
