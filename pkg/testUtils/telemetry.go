package testUtils

import (
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/metric"
	noopMetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	noopTrace "go.opentelemetry.io/otel/trace/noop"
)

// NoopTracer returns a tracer that does nothing
func NewNoopTracer() trace.Tracer {
	return noopTrace.NewTracerProvider().Tracer("")
}

func NewNoopMeter() metric.Meter {
	return noopMetric.NewMeterProvider().Meter("")
}

func NewNoopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
