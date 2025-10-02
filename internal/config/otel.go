package config

import (
	sd "apistarter/internal/shutdown"
	"context"
	"errors"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	golog "log"
)

func NewOtelSdk(ctx context.Context, closer *sd.ShutdownManager) trace.Tracer {
	tracer := otel.Tracer("apistarter")
	// Skip initialization if disabled
	if os.Getenv("OTEL_ENABLED") != "true" {
		golog.Println("🔕 OpenTelemetry disabled (OTEL_ENABLED != true)")
		return nil
	}

	var shutdownFuncs []func(context.Context) error

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		return err
	}

	// Set up propagator
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Create shared resource
	res, err := createResource(ctx)
	if err != nil {
		golog.Printf("⚠️ Failed to create OpenTelemetry resource: %v", err)
		return tracer
	}

	// Set up trace provider with OTLP
	tracerProvider, err := newTracerProvider(ctx, res)
	if err != nil {
		golog.Printf("⚠️ Failed to initialize tracer provider: %v", err)
	} else {
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
		golog.Println("✅ Tracing initialized with OTLP HTTP")
	}

	// Set up metric provider with OTLP
	meterProvider, err := newMetricProvider(ctx, res)
	if err != nil {
		golog.Printf("⚠️ Failed to initialize metric provider: %v", err)
	} else {
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
		golog.Println("✅ Metrics initialized with OTLP HTTP")
	}

	golog.Println("✅ OpenTelemetry fully initialized with OTLP")
	closer.CleanupFuncsWithContext = append(closer.CleanupFuncsWithContext, shutdown)
	return tracer
}

func createResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("apistarter"),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
		resource.WithFromEnv(), // Reads OTEL_RESOURCE_ATTRIBUTES
		resource.WithTelemetrySDK(),
		resource.WithHost(),    // Adds host information
		resource.WithOS(),      // Adds OS information
		resource.WithProcess(), // Adds process information
	)
}

func newTracerProvider(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	// OTLP HTTP exporter - compatible with Jaeger, Tempo, etc.
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"), // OTLP HTTP endpoint
		otlptracehttp.WithInsecure(),                 // Use WithTLS() in production
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)
	return tracerProvider, nil
}

func newMetricProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	// OTLP HTTP metric exporter
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("localhost:4318"), // OTLP HTTP endpoint
		otlpmetrichttp.WithInsecure(),                 // Use WithTLS() in production
		otlpmetrichttp.WithURLPath("/v1/traces"),      // OTLP path
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			// Export metrics every 15 seconds
			sdkmetric.WithInterval(15*time.Second),
		)),
	)
	return meterProvider, nil
}
