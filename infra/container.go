package infra

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"gostarter/infra/config"
	"log/slog"
)

type Container struct {
	Cfg    *config.Config
	Logger *slog.Logger
	Tracer *trace.Tracer
	Meter  *metric.Meter
}
