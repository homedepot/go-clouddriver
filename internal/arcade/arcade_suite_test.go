package arcade_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestArcade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arcade Suite")
}
