package front50_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFront50(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Front50 Suite")
}
