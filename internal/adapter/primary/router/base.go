package router

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/internal/constant"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/cerrors"
)

type BaseRouter struct{}

func (r BaseRouter) WrappedHandler(f func(c *gin.Context) *cerrors.CodedError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			httpStatusCode := http.StatusInternalServerError
			if err.Code == constant.ErrUnAuthorized {
				httpStatusCode = http.StatusUnauthorized
			}
			logger := util.LoggerFromContext(c.Request.Context()).WithGroup("GinRouter")
			logger.Error("HTTP Error",
				slog.Any("error", err),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.RawPath),
				slog.String("query", c.Request.URL.RawQuery),
				slog.String("stack", string(debug.Stack())),
			)

			c.AbortWithStatusJSON(
				httpStatusCode,
				gin.H{"status": false, "code": err.Code, "message": err.Error()},
			)
		}
	}
}

func (r BaseRouter) Success(c *gin.Context, data any) *cerrors.CodedError {
	c.JSON(http.StatusOK, gin.H{"status": true, "data": data})
	return nil
}
