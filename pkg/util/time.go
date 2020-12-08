package util

import (
	"time"
)

func CurrentTimeUTC() time.Time {
	utc := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	t, _ := time.Parse("2006-01-02T15:04:05.999Z", utc)

	return t
}
