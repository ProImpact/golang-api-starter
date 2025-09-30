package config

import (
	"apistarter/internal/shutdown"
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func SetTracer(cfg *Configuration, shutdownFuns *shutdown.ShutdownManager) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.JeagerUrl),
		),
	)
	if err != nil {
		log.Fatal("Error al crear el exporter de Jaeger: ", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.AppName),
		attribute.String("service.version", "1.0.0"),
		attribute.String("deployment.environment", cfg.Mode),
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	shutdownFuns.CleanupFuncs = append(shutdownFuns.CleanupFuncs, func() error {
		return tp.Shutdown(context.Background())
	})
}
