package middleware

import (
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

// responseWriterWrapper captures the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func NewLatencyMiddleware(meter metric.Meter) Middleware {
	return func(next http.Handler) http.Handler {

		// Create histogram for request duration
		histogram, err := meter.Float64Histogram(
			"http_request_duration_seconds",
			metric.WithDescription("HTTP request latency in seconds"),
			metric.WithUnit("s"),
		)
		if err != nil {

			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start measuring latency
			start := time.Now()

			// Call the next handler
			pattern := chi.RouteContext(r.Context()).RoutePattern()
			wrapper := &responseWriterWrapper{w, http.StatusOK}

			next.ServeHTTP(wrapper, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			histogram.Record(r.Context(), duration,
				metric.WithAttributes(
					attribute.String("path", pattern),
					attribute.String("method", r.Method),
					attribute.Int("status_code", wrapper.statusCode),
				),
			)
		})
	}
}
