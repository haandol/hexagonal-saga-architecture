package util

import (
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	UTC       bool
	SkipPaths []string
}

func GinSlog(logger *slog.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return GinSlogWithConfig(logger, &Config{UTC: utc})
}

func GinSlogWithConfig(logger *slog.Logger, conf *Config) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			} else {
				fields := []any{
					slog.Int("status", c.Writer.Status()),
					slog.String("request-id", c.Request.Header.Get("X-Request-ID")),
					slog.String("method", c.Request.Method),
					slog.String("path", path),
					slog.String("query", query),
					slog.String("ip", c.ClientIP()),
					slog.String("user-agent", c.Request.UserAgent()),
					slog.Duration("latency", latency),
					slog.String("time", end.Format(time.RFC3339)),
				}
				logger.Info(path, fields...)
			}
		}
	}
}

func RecoveryWithSlog(logger *slog.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						slog.Time("time", time.Now()),
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
						slog.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						slog.Time("time", time.Now()),
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
