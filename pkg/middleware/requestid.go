package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	xRequestIDKey = "X-Request-ID"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get(xRequestIDKey) == "" {
			xRequestId := uuid.NewString()
			c.Request.Header.Set(xRequestIDKey, xRequestId)
			c.Set(xRequestIDKey, xRequestId)
		}
		c.Next()
	}
}
