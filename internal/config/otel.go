package config

import (
	sd "apistarter/internal/shutdown"
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	golog "log"
)

func NewOtelSdk(ctx context.Context, closer *sd.ShutdownManager) trace.Tracer {
	// Skip initialization if disabled
	if os.Getenv("OTEL_ENABLED") != "true" {
		golog.Println("🔕 OpenTelemetry disabled (OTEL_ENABLED != true)")
		return otel.Tracer("noop")
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
		return otel.Tracer("noop")
	}

	// Set up trace provider with OTLP - SOLO TRAZAS
	tracerProvider, err := newTracerProvider(ctx, res)
	if err != nil {
		golog.Printf("⚠️ Failed to initialize tracer provider: %v", err)
	} else {
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
		golog.Println("✅ Tracing initialized with Jaeger OTLP HTTP")
	}

	golog.Println("✅ OpenTelemetry initialized for traces only")
	closer.CleanupFuncsWithContext = append(closer.CleanupFuncsWithContext, shutdown)

	return otel.Tracer("apistarter")
}

func createResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("apistarter"),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
	)
}

func newTracerProvider(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithURLPath("/v1/traces"), 
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
