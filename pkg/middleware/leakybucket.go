package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/util"
	"go.uber.org/ratelimit"
)

func LeakBucket(cfg *config.App) gin.HandlerFunc {
	logger := util.GetLogger().With(
		"module", "Middleware",
		"func", "LeakBucket",
	)

	var limiter ratelimit.Limiter
	if cfg.RPS == 0 {
		limiter = ratelimit.NewUnlimited()
	} else {
		limiter = ratelimit.New(cfg.RPS)
	}

	prev := time.Now()
	return func(c *gin.Context) {
		now := limiter.Take()
		logger.Debugf("%v", now.Sub(prev))
		prev = now
		c.Next()
	}
}
