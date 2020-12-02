package clouddriver

import (
	"time"

	"github.com/gin-gonic/gin"
)

// This represents a clouddriver error.
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Timestamp int64  `json:"timestamp"`
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
	c.Status(status)
	_ = c.Error(err).SetType(gin.ErrorTypePublic)
}
