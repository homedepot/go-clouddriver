package fiat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFiat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fiat Suite")
}
