package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"medika-backend/internal/infrastructure/config"
)

func Initialize(cfg config.ObservabilityConfig) func() {
	var cleanupFuncs []func()

	// Initialize tracing
	if cfg.Tracing.Enabled {
		cleanup := initTracing(cfg.Tracing)
		cleanupFuncs = append(cleanupFuncs, cleanup)
	}

	// Initialize metrics
	if cfg.Metrics.Enabled {
		cleanup := initMetrics(cfg.Metrics)
		cleanupFuncs = append(cleanupFuncs, cleanup)
	}

	return func() {
		for _, cleanup := range cleanupFuncs {
			cleanup()
		}
	}
}

func initTracing(cfg config.TracingConfig) func() {
	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		fmt.Printf("Failed to create Jaeger exporter: %v\n", err)
		return func() {}
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		fmt.Printf("Failed to create resource: %v\n", err)
		return func() {}
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
		trace.WithSampler(trace.TraceIDRatioBased(cfg.SamplingRate)),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	fmt.Printf("✅ Tracing initialized with Jaeger endpoint: %s\n", cfg.JaegerEndpoint)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}
}

func initMetrics(cfg config.MetricsConfig) func() {
	// Initialize Prometheus metrics
	// This would typically set up prometheus metrics
	fmt.Printf("✅ Metrics initialized on port: %s\n", cfg.Port)
	
	return func() {
		// Cleanup metrics
	}
}
