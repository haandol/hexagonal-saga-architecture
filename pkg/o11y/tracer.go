package o11y

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/contrib/detectors/aws/eks"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

func initTracer(endpoint string) ShutdownFunc {
	ctx := context.Background()
	// Create and start new OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		otel.SetTracerProvider(nooptrace.NewTracerProvider())
		otel.SetTextMapPropagator(xray.Propagator{})
		tracer = otel.Tracer("hexagonal")
		return func(ctx context.Context) error { return nil }
	}

	// AWS EKS resource
	resourceDetector := eks.NewResourceDetector()
	resource, err := resourceDetector.Detect(context.Background())
	if err != nil {
		// just use nil-resource if failed to detect resource
		log.Printf("Failed to create new resource: %v", err)
	}

	xrayIDG := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(xrayIDG),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	tracer = otel.Tracer("hexagonal")

	return tp.Shutdown
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func BeginSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
	span.SetAttributes(
		attribute.String("service.name", name),
	)
	return ctx, span
}

func BeginSpanWithTraceID(ctx context.Context, corrID, parentID, name string) (context.Context, trace.Span) {
	traceID, err := trace.TraceIDFromHex(corrID)
	if err != nil {
		log.Printf("Failed to parse traceID: %v", err)
	}

	spanID, err := trace.SpanIDFromHex(parentID)
	if err != nil {
		log.Printf("Failed to parse spanID: %v", err)
	}

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled.WithSampled(true),
		Remote:     true,
	})

	ctx, span := tracer.Start(
		trace.ContextWithSpanContext(ctx, spanContext),
		name,
		trace.WithSpanKind(trace.SpanKindServer),
	)
	span.SetAttributes(
		attribute.String("TraceId", GetXrayTraceID(traceID.String())),
		attribute.String("ParentSpanId", parentID),
		attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue(name),
		},
	)

	return ctx, span
}

func BeginSubSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func BeginSubSpanWithNode(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
}

func GetTraceSpanID(ctx context.Context) (traceID, spanID string) {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return traceID, spanID
	}

	traceID = spanContext.TraceID().String()
	spanID = spanContext.SpanID().String()
	return traceID, spanID
}

func GetXrayTraceID(traceID string) string {
	if traceID == "" {
		return ""
	}
	return fmt.Sprintf("1-%s-%s", traceID[0:8], traceID[8:])
}

func AttrString(k, v string) attribute.KeyValue {
	return attribute.String(k, v)
}

func AttrInt(k string, v int) attribute.KeyValue {
	return attribute.Int(k, v)
}

func BuildKafkaMessageAttr(topic, key, id string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(topic),
		semconv.MessagingDestinationKindKey.String(key),
		semconv.MessagingMessageIDKey.String(id),
	}
	return attrs
}

func GetStatus(err error) (code codes.Code, msg string) {
	code = codes.Ok
	if err != nil {
		code = codes.Error
		msg = fmt.Sprintf("%v", err)
	}

	return
}
