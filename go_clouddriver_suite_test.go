package go_clouddriver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoClouddriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoClouddriver Suite")
}
