package core_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Search", func() {
	Describe("#Search", func() {
		accountsHeader := ""

		BeforeEach(func() {
			setup()
			uri = svr.URL + "/search?pageSize=500&q=default&type=pod"
			accountsHeader = "account1,account2"
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			createRequest(http.MethodGet)
			req.Header.Add("X-Spinnaker-Accounts", accountsHeader)
			doRequest()
		})

		When("kind and namespace are not provided", func() {
			BeforeEach(func() {
				uri = svr.URL + "/search?pageSize=500"
			})

			It("returns status bad request", func() {
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				ce := getClouddriverError()
				Expect(ce.Error).To(Equal("Bad Request"))
				Expect(ce.Message).To(Equal("must provide query params 'q' to specify the namespace and 'type' to specify the kind"))
				Expect(ce.Status).To(Equal(http.StatusBadRequest))
			})
		})

		When("an empty account is provided", func() {
			BeforeEach(func() {
				accountsHeader = "account1,,account2"
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearchEmptyResponse)
			})
		})

		When("pageSize is not provided", func() {
			BeforeEach(func() {
				uri = svr.URL + "/search?q=default&type=pod"
			})

			It("returns an empty response", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearchEmptyResponseWithPageSizeZero)
			})
		})

		When("listing resource names returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourceNamesByAccountNameAndKindAndNamespaceReturns(nil, errors.New("error listing names"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearchEmptyResponse)
			})
		})

		When("it reached the pageSize limit while listing names", func() {
			BeforeEach(func() {
				uri = svr.URL + "/search?pageSize=3&q=default&type=pod"
				fakeSQLClient.ListKubernetesResourceNamesByAccountNameAndKindAndNamespaceReturns([]string{"test-name1", "test-name2"}, nil)
			})

			It("limits the response", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearchWithPageSizeThree)
			})
		})

		When("it succeeds", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourceNamesByAccountNameAndKindAndNamespaceReturns([]string{"test-name1", "test-name2"}, nil)
			})

			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				validateResponse(payloadSearch)
			})
		})
	})
})
