package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/util"
)

const headerTraceID = "X-Amzn-Trace-Id"

func XrayTracing() gin.HandlerFunc {
	util.InitXray()

	return func(c *gin.Context) {
		name := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		ctx, seg := xray.BeginSegment(c.Request.Context(), name)
		defer seg.Close(nil)

		ctx = context.WithValue(ctx, constant.CtxTraceKey, seg.TraceID)

		c.Request = c.Request.WithContext(ctx)

		captureRequestData(c, seg)
		c.Next()
		captureResponseData(c, seg)
	}
}

// Write request data to segment
func captureRequestData(c *gin.Context, seg *xray.Segment) {
	r := c.Request
	seg.Lock()
	defer seg.Unlock()
	segmentRequest := seg.GetHTTP().GetRequest()
	segmentRequest.Method = r.Method
	segmentRequest.URL = r.URL.String()
	segmentRequest.XForwardedFor = hasXForwardedFor(r)
	segmentRequest.ClientIP = clientIP(r)
	segmentRequest.UserAgent = r.UserAgent()
	c.Writer.Header().Set(headerTraceID, fmt.Sprintf("Root=%s", seg.TraceID))
}

// Write response data to segment
func captureResponseData(c *gin.Context, seg *xray.Segment) {
	respStatus := c.Writer.Status()

	seg.Lock()
	defer seg.Unlock()
	seg.GetHTTP().GetResponse().Status = respStatus
	seg.GetHTTP().GetResponse().ContentLength = c.Writer.Size()

	if respStatus >= 400 && respStatus < 500 {
		seg.Error = true
	}
	if respStatus == 429 {
		seg.Throttle = true
	}
	if respStatus >= 500 && respStatus < 600 {
		seg.Fault = true
	}
}

func hasXForwardedFor(r *http.Request) bool {
	return r.Header.Get("X-Forwarded-For") != ""
}

func clientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
	}

	return r.RemoteAddr
}
