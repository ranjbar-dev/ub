package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit returns a Gin middleware that limits request body size.
// Bodies exceeding maxBytes are rejected with HTTP 413.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}
