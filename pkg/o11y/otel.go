package o11y

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel/trace"
)

var (
	tracer         trace.Tracer
	initialized    = false
	tracerShutdown ShutdownFunc
	metricShutdown ShutdownFunc
)

func InitOtel() {
	if initialized {
		return
	}

	initialized = true

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "0.0.0.0:4317" // setting default endpoint for exporter
	}
	log.Printf("OTEL endpoint: %s", endpoint)

	tracerShutdown = initTracer(endpoint)
	metricShutdown = initMetricProvider(endpoint)

	log.Printf("OTEL initialized")
}

func Close(ctx context.Context) error {
	if tracerShutdown != nil {
		if err := tracerShutdown(ctx); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		} else {
			log.Println("tracer shutdown")
		}
	}

	if metricShutdown != nil {
		if err := metricShutdown(ctx); err != nil {
			log.Printf("failed to shutdown metric: %v", err)
		} else {
			log.Println("metric shutdown")
		}
	}

	return nil
}
