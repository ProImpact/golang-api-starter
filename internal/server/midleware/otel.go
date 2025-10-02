package midleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		defer span.End()

		// Set common attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.url", c.Request.URL.String()),
		)

		// Replace the request with context that has the span
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// Set status code after request is processed
		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
	}
}
