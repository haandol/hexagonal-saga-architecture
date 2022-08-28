package middleware

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/pkg/config"
)

func Timeout(cfg config.App) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(time.Duration(cfg.TimeoutSec)*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
	)
}
