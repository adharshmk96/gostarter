package middleware

import (
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

func NewCounterMiddleware(meter metric.Meter) Middleware {
	counter, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests made."),
	)
	if err != nil {
		return nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Increment the counter
			counter.Add(r.Context(), 1)

			next.ServeHTTP(w, r)
		})
	}
}
