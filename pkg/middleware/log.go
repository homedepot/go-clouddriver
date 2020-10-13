package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

var (
	bold = color.New(color.FgWhite, color.Bold).SprintFunc()
)

// A verbose request/response logger.
func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
			buf bytes.Buffer
		)

		clone := c.Request.Clone(context.TODO())

		buf.ReadFrom(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(&buf)
		clone.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		if err != nil {
			log.Println("[MIDDLEWARE] error getting request body in verbose request logger:", err.Error())
		} else {
			b, _ := ioutil.ReadAll(clone.Body)
			buffer := &bytes.Buffer{}

			buffer.WriteString(bold("REQUEST: ["+time.Now().In(time.UTC).Format(time.RFC3339)) + bold("]"))
			buffer.WriteByte('\n')
			buffer.WriteString(fmt.Sprintf("%s %s %s", clone.Method, clone.URL, clone.Proto))
			buffer.WriteByte('\n')
			buffer.WriteString(fmt.Sprintf("Host: %s", clone.Host))
			buffer.WriteByte('\n')
			buffer.WriteString(fmt.Sprintf("Accept: %s", clone.Header.Get("Accept")))
			buffer.WriteByte('\n')
			buffer.WriteString(fmt.Sprintf("User-Agent: %s", clone.Header.Get("User-Agent")))
			buffer.WriteByte('\n')
			buffer.WriteString(fmt.Sprintf("Headers: %s", clone.Header))
			buffer.WriteByte('\n')
			if len(b) > 0 {
				j, err := json.MarshalIndent(b, "", "    ")
				if err != nil {
					log.Fatal("[MIDDLEWARE] failed to generate json", err)
				} else {
					body, _ := base64.StdEncoding.DecodeString(string(j[1 : len(j)-1]))
					buffer.Write(body)
					buffer.WriteByte('\n')
				}
			}

			fmt.Println(buffer.String())
		}

		c.Next() // execute all the handlers
	}
}
