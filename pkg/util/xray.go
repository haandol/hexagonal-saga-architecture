package util

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/awsplugins/ecs"
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

	ecs.Init()

	s, _ := sampling.NewLocalizedStrategyFromFilePath("../../xray.json")
	if err := xray.Configure(xray.Config{
		ServiceVersion:   "1.2.3",
		SamplingStrategy: s,
	}); err != nil {
		panic(err)
	}
}

func BeginSegmentWithTraceID(ctx context.Context,
	traceID string, parentID string, name string,
) (context.Context, *xray.Segment) {
	con, seg := xray.BeginSegment(ctx, name)

	seg.TraceID = traceID
	seg.ParentID = parentID

	return con, seg
}

func BeginSubSegment(ctx context.Context, name string) (context.Context, *xray.Segment) {
	return xray.BeginSubsegment(ctx, name)
}

func GetSegmentID(ctx context.Context) string {
	seg := xray.GetSegment(ctx)
	if seg != nil {
		return seg.ID
	}
	return ""
}
