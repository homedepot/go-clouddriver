package v0_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestV0(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V0 Suite")
}
