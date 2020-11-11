package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/billiford/go-clouddriver/pkg/fiat"
	"github.com/billiford/go-clouddriver/pkg/fiat/fiatfakes"
	"github.com/billiford/go-clouddriver/pkg/http/core"
	"github.com/billiford/go-clouddriver/pkg/middleware"
	. "github.com/billiford/go-clouddriver/pkg/middleware"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	c                                                       *gin.Context
	hf                                                      gin.HandlerFunc
	fakeFiatClient                                          *fiatfakes.FakeClient
	r                                                       *http.Request
	err                                                     error
	testUser, testApplication, testAccount, authorizeErrMsg string
	allApps, filteredApps                                   core.Applications
	authorizedAppsMap                                       map[string]fiat.Application
)

var _ = Describe("Auth", func() {
	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)
		c, _ = gin.CreateTestContext(httptest.NewRecorder())
		fakeFiatClient = &fiatfakes.FakeClient{}
		c.Set(fiat.ClientInstanceKey, fakeFiatClient)
		r, err = http.NewRequest(http.MethodGet, "", nil)
		Expect(err).To(BeNil())
		c.Request = r
		testUser = "test-user"
		testAccount = "test-account"
		testApplication = "test-application"
		r.Header.Set("X-Spinnaker-User", testUser)
		r.Header.Set("X-Spinnaker-Application", testApplication)
	})

	Describe("#AuthApplication", func() {
		BeforeEach(func() {
			hf = AuthApplication("READ")
		})

		JustBeforeEach(func() {
			hf(c)
		})

		When("user is empty", func() {
			BeforeEach(func() {
				r.Header.Del("X-Spinnaker-User")
			})

			It("calls c.Next", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeFiatClient.AuthorizeCallCount()).To(BeZero())
			})
		})

		When("app is empty", func() {
			BeforeEach(func() {
				r.Header.Del("X-Spinnaker-Application")
			})

			It("calls c.Next", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeFiatClient.AuthorizeCallCount()).To(BeZero())
			})
		})

		When("fiatClient.Authorize returns an error", func() {
			BeforeEach(func() {
				fakeFiatClient.AuthorizeReturns(fiat.Response{}, errors.New("fake error"))
			})

			It("returns status unauthorized", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusUnauthorized))
				Expect(c.Errors[0].Error()).To(Equal("fake error"))
			})
		})

		When("the user doesn't have the permission", func() {
			BeforeEach(func() {
				fakeResp := fiat.Response{}
				fakeResp.Name = testUser
				fakeApplication := fiat.Application{
					Name:           testApplication,
					Authorizations: []string{"WRITE"},
				}
				fakeResp.Applications = []fiat.Application{fakeApplication}
				fakeFiatClient.AuthorizeReturns(fakeResp, nil)
			})

			It("returns status Forbidden", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusForbidden))
				Expect(c.Errors[0].Error()).To(Equal("Access denied to application test-application - required authorization: READ"))
			})
		})

		When("the user has the permission", func() {
			BeforeEach(func() {
				fakeResp := fiat.Response{}
				fakeResp.Name = testUser
				fakeApplication := fiat.Application{
					Name:           testApplication,
					Authorizations: []string{"READ"},
				}
				fakeResp.Applications = []fiat.Application{fakeApplication}
				fakeFiatClient.AuthorizeReturns(fakeResp, nil)
			})

			It("returns status Forbidden", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("#AuthAccount", func() {
		BeforeEach(func() {
			hf = AuthAccount("READ")
			c.Params = append(c.Params, gin.Param{
				Key:   "account",
				Value: testAccount,
			})
		})

		JustBeforeEach(func() {
			hf(c)
		})

		When("user is empty", func() {
			BeforeEach(func() {
				r.Header.Del("X-Spinnaker-User")
			})

			It("calls c.Next", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeFiatClient.AuthorizeCallCount()).To(BeZero())
			})
		})

		When("account is missing from path params", func() {
			BeforeEach(func() {
				newParams := []gin.Param{}
				for _, p := range c.Params {
					if p.Key != "account" {
						newParams = append(newParams, gin.Param{
							Key:   p.Key,
							Value: p.Value,
						})
					}
				}
				c.Params = newParams
			})

			It("calls c.Next", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
				Expect(fakeFiatClient.AuthorizeCallCount()).To(BeZero())
			})
		})

		When("fiatClient.Authorize returns an error", func() {
			BeforeEach(func() {
				fakeFiatClient.AuthorizeReturns(fiat.Response{}, errors.New("fake error"))
			})

			It("returns status unauthorized", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusUnauthorized))
				Expect(c.Errors[0].Error()).To(Equal("fake error"))
			})
		})

		When("the user doesn't have the permission", func() {
			BeforeEach(func() {
				fakeResp := fiat.Response{}
				fakeResp.Name = testUser
				fakeAccount := fiat.Account{
					Name:           testAccount,
					Authorizations: []string{"WRITE"},
				}
				fakeResp.Accounts = []fiat.Account{fakeAccount}
				fakeFiatClient.AuthorizeReturns(fakeResp, nil)
			})

			It("returns status Forbidden", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusForbidden))
				Expect(c.Errors[0].Error()).To(Equal("Access denied to account test-account - required authorization: READ"))
			})
		})

		When("the user has the permission", func() {
			BeforeEach(func() {
				fakeResp := fiat.Response{}
				fakeResp.Name = testUser
				fakeAccount := fiat.Account{
					Name:           testAccount,
					Authorizations: []string{"READ"},
				}
				fakeResp.Accounts = []fiat.Account{fakeAccount}
				fakeFiatClient.AuthorizeReturns(fakeResp, nil)
			})

			It("returns status OK", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("#FilterAuthorizedApplications", func() {
		BeforeEach(func() {
			hf = PostFilterAuthorizedApplications("READ")
			allApps = append(allApps, core.Application{
				Name: "test-app1",
			})
			c.Set(core.KeyAllApplications, allApps)
		})

		JustBeforeEach(func() {
			hf(c)
		})

		AfterEach(func() {
			c.Abort()
		})

		When("there is an error attached to the context", func() {
			BeforeEach(func() {
				c.Errors = append(c.Errors, c.Error(errors.New("fake error")))
			})

			It("returns", func() {

			})

		})

		When("user is missing from header", func() {
			BeforeEach(func() {
				r.Header.Del("X-Spinnaker-User")
			})

			It("returns the list of all apps without filtering", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))
			})
		})

		When("fiat returns an error", func() {
			BeforeEach(func() {
				fakeApplication := fiat.Application{
					Name:           testApplication,
					Authorizations: []string{"WRITE"},
				}
				fakeResp := fiat.Response{}
				fakeResp.Name = testApplication
				fakeResp.Applications = []fiat.Application{fakeApplication}
				authorizeErrMsg = "user authorization error: 401"
				fakeFiatClient.AuthorizeReturns(fiat.Response{}, errors.New(authorizeErrMsg))
			})

			It("returns an error", func() {
				Expect(c.Writer.Status()).To(Equal(http.StatusUnauthorized))
				Expect(c.Errors[0].Error()).To(Equal(authorizeErrMsg))
			})
		})
	})

	Describe("#FilterAuthorizedApps", func() {
		BeforeEach(func() {
			allApps = []core.Application{}
			authorizedAppsMap = make(map[string]fiat.Application)
		})

		When("allApps is empty", func() {
			BeforeEach(func() {
				allApps = []core.Application{}
				filteredApps = middleware.FilterAuthorizedApps(authorizedAppsMap, allApps, "READ")
			})

			It("returns an empty list", func() {
				Expect(len(filteredApps)).To(BeZero())
			})
		})

		When("authorizedAppsMap is empty", func() {
			BeforeEach(func() {
				allApps = append(allApps, core.Application{
					Name: "app1",
				})
				allApps = append(allApps, core.Application{
					Name: "app2",
				})
				filteredApps = middleware.FilterAuthorizedApps(authorizedAppsMap, allApps, "READ")
			})

			It("returns an empty list", func() {
				Expect(len(filteredApps)).To(BeZero())
			})
		})

		When("the user doesn't have the required permission", func() {
			BeforeEach(func() {
				authorizedAppsMap["app1"] = fiat.Application{
					Name:           "app1",
					Authorizations: []string{"READ"},
				}
				allApps = append(allApps, core.Application{
					Name: "app1",
				})
				allApps = append(allApps, core.Application{
					Name: "app2",
				})
				filteredApps = middleware.FilterAuthorizedApps(authorizedAppsMap, allApps, "WRITE")
			})

			It("returns an empty list", func() {
				Expect(len(filteredApps)).To(BeZero())
			})
		})

		When("the user has the required permission", func() {
			BeforeEach(func() {
				authorizedAppsMap = make(map[string]fiat.Application)
				authorizedAppsMap["app1"] = fiat.Application{
					Name:           "app1",
					Authorizations: []string{"READ"},
				}
				allApps = append(allApps, core.Application{
					Name: "app1",
				})
				allApps = append(allApps, core.Application{
					Name: "app2",
				})
				filteredApps = middleware.FilterAuthorizedApps(authorizedAppsMap, allApps, "READ")
			})

			It("returns an empty list", func() {
				Expect(len(filteredApps)).To(Equal(1))
			})
		})
	})
})
