package core_test

import (
	// . "github.com/billiford/go-clouddriver/pkg/http/v0"

	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Credential", func() {
	Describe("#CreateKubernetesDeployment", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/credentials"
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
				validateResponse(payloadCredentials)
			})
		})
	})
})
