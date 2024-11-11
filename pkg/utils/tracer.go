package utils

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type StopFunc func()

func TraceSpan(ctx context.Context, tracer trace.Tracer, name string) (context.Context, StopFunc) {
	if tracer == nil {
		return ctx, func() {}
	}

	_, span := tracer.Start(ctx, name)
	return ctx, func() {
		span.End()
	}
}
