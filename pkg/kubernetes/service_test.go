package kubernetes_test

import (
	. "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Service", func() {
	var (
		s Service
	)

	BeforeEach(func() {
		m := map[string]interface{}{}
		s = NewService(m)
	})

	Describe("#Selector", func() {
		BeforeEach(func() {
			_ = s.Selector()
		})

		It("succeeds", func() {
		})
	})
})
