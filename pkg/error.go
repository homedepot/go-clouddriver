package clouddriver

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This represents a clouddriver error.
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type ErrorMeta struct {
	FuncName string
	FileName string
	GUID     string
	LineNum  int
}

// Example.
// {
//   "error":"Forbidden",
//   "message":"Access denied to account spin-cluster-account - required authorization: READ",
//   "status":403,
//   "timestamp":1597608027851
// }
func NewError(err, message string, status int) ErrorResponse {
	return ErrorResponse{
		Error:     err,
		Message:   message,
		Status:    status,
		Timestamp: time.Now().UnixNano() / 1000000,
	}
}

// Error attaches a given Go error to a gin context and sets its type to public.
func Error(c *gin.Context, status int, err error) {
	pc, fn, ln, _ := runtime.Caller(1)
	m := ErrorMeta{
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fn,
		GUID:     uuid.New().String(),
		LineNum:  ln,
	}

	c.Status(status)
	_ = c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(m)
}

func Meta(msg *gin.Error) ErrorMeta {
	return msg.Meta.(ErrorMeta)
}
