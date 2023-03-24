package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func OtelTracing(service string) gin.HandlerFunc {
	return otelgin.Middleware(service, otelgin.WithFilter(func(r *http.Request) bool {
		// Don't trace health checks
		return r.URL.Path != "/healthy"
	}))
}
