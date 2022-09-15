package util

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	initialized = false
)

func InitXray() {
	if initialized {
		return
	}

	initialized = true
	s, _ := sampling.NewLocalizedStrategyFromFilePath("../../xray.json")
	if err := xray.Configure(xray.Config{
		DaemonAddr:       "127.0.0.1:2000", // default
		ServiceVersion:   "1.2.3",
		SamplingStrategy: s,
	}); err != nil {
		panic(err)
	}
}

func BeginSegmentWithTraceID(ctx context.Context, traceID string, name string) (context.Context, *xray.Segment) {
	con, seg := xray.BeginSegment(ctx, name)

	seg.TraceID = traceID

	return con, seg
}

func BeginSubSegment(ctx context.Context, name string) (context.Context, *xray.Segment) {
	return xray.BeginSubsegment(ctx, name)
}
