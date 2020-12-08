package util_test

import (
	"time"

	. "github.homedepot.com/cd/skipper/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time", func() {
	It("formats and sets the time to UTC", func() {
		now := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		t, _ := time.Parse("2006-01-02T15:04:05.999Z", now)
		utc := CurrentTimeUTC()

		Expect(utc).To(BeTemporally("~", t))
	})
})
