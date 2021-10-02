package arcade_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArcade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arcade Suite")
}
