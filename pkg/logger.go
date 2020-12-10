package clouddriver

import (
	"log"
	"runtime"
	"time"

	"github.com/google/uuid"
)

const (
	white = "\033[90;47m"
	reset = "\033[0m"
)

// Log logs a given error. It checks if meta was passed in. If
// no meta was passed in, it defines the meta as the function and
// line number of what called the Log func.
//
// See https://stackoverflow.com/questions/24809287/how-do-you-get-a-golang-program-to-print-the-line-number-of-the-error-it-just-ca
func Log(err error, meta ...ErrorMeta) {
	if len(meta) == 0 {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, ln, _ := runtime.Caller(1)
		m := ErrorMeta{
			FuncName: runtime.FuncForPC(pc).Name(),
			FileName: fn,
			GUID:     uuid.New().String(),
			LineNum:  ln,
		}
		meta = append(meta, m)
	}

	m := meta[0]

	log.SetFlags(0)
	log.Printf("[CLOUDDRIVER] %v |%s %s %s| %s | %s | %s:%d | %v\n",
		time.Now().In(time.UTC).Format("2006/01/02 - 15:04:05"),
		white, "LOG", reset,
		m.GUID,
		m.FuncName,
		m.FileName,
		m.LineNum,
		err,
	)
}
