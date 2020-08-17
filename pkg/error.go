package clouddriver

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Error struct {
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
func NewError(e, m string, s int) Error {
	return Error{
		Error:     e,
		Message:   m,
		Status:    s,
		Timestamp: time.Now().UnixNano() / 1000000,
	}
}

func WriteError(c *gin.Context, status int, err error) {
	c.Status(status)
	c.Error(err).SetType(gin.ErrorTypePublic)
}
