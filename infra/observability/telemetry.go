package observability

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gostarter/infra/config"
	"time"
)

type telemetryService struct {
	tp *sdktrace.TracerProvider
	mp *sdkmetric.MeterProvider
}

type TelemetryService interface {
	GetTracerProvider() *sdktrace.TracerProvider
	GetMeterProvider() *sdkmetric.MeterProvider
	Stop()
}

func (o *telemetryService) GetTracerProvider() *sdktrace.TracerProvider {
	return o.tp
}

func (o *telemetryService) GetMeterProvider() *sdkmetric.MeterProvider {
	return o.mp
}

func (o *telemetryService) Stop() {
	_ = o.tp.Shutdown(context.Background())
	_ = o.mp.Shutdown(context.Background())
}

func NewTelemetryService(cfg *config.ObservabilityConfig) (TelemetryService, error) {
	tp, err := newTracerProvider(cfg)
	if err != nil {
		return nil, err
	}

	mp, err := newMeterProvider(cfg)
	if err != nil {
		return nil, err
	}

	//tracer := tp.Tracer(cfg.TracerName)
	//meter := mp.Meter(cfg.MeterName)

	// TODO: Implement telemetry init and usage.
	ts := &telemetryService{
		tp: tp,
		mp: mp,
	}

	return ts, nil
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("gostarter"),
	)
}

func newTracerProvider(cfg *config.ObservabilityConfig) (*sdktrace.TracerProvider, error) {
	// Setup Trace Exporter ( tempo )
	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(cfg.TraceExporter),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(
			newResource(),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func newMeterProvider(cfg *config.ObservabilityConfig) (*sdkmetric.MeterProvider, error) {
	// Set up Metric Exporter ( otlp collector )
	metricExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(cfg.MeterExporter),
	)
	if err != nil {
		return nil, err
	}

	// Create a MeterProvider with the exporter
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(1*time.Second)),
		),
		sdkmetric.WithResource(
			newResource(),
		),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}

//func newLoggerProvider(cfg *config.ObservabilityConfig) (*sdklog.LoggerProvider, error) {
//	// Setup Log Exporter ( tempo )
//	logExporter, err := otlptracehttp.New(
//		context.Background(),
//		otlptracehttp.WithInsecure(),
//		otlptracehttp.WithEndpoint(cfg.LogExporter),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	lp := sdklog.NewLoggerProvider(
//		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
//	)
//
//	otel.SetLogger(lp)
//
//	return lp, nil
//
//}
