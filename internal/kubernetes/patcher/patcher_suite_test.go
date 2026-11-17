package patcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Patcher Suite")
}
