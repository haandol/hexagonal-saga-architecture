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
		xRequestId := uuid.New().String()
		c.Request.Header.Set(xRequestIDKey, xRequestId)
		c.Set(xRequestIDKey, xRequestId)
		c.Next()
	}
}
