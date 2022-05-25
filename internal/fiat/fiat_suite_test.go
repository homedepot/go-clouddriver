package fiat_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFiat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fiat Suite")
}
