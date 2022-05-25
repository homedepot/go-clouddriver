package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/homedepot/go-clouddriver/internal/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CacheControl", func() {
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
		mw := CacheControl(30)
		mw(c)
	})

	It("sets the Cache-Control header", func() {
		Expect(c.Writer.Header().Get("Cache-Control")).To(Equal("public, max-age=30"))
	})
})
