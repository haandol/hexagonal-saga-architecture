package o11y

import (
	"context"
	"log"

	"go.opentelemetry.io/contrib/detectors/aws/eks"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func initMetricProvider(endpoint string) ShutdownFunc {
	exporter, err := otlpmetricgrpc.New(
		context.TODO(),
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		otel.SetMeterProvider(noopmetric.NewMeterProvider())
		return func(ctx context.Context) error { return nil }
	}

	// AWS EKS resource
	resourceDetector := eks.NewResourceDetector()
	resource, err := resourceDetector.Detect(context.Background())
	if err != nil {
		// just use nil-resource if failed to detect resource
		log.Printf("Failed to create new resource: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(provider)

	return exporter.Shutdown
}
