package otel_sdk

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func Tracer(traceName string, opts ...trace.TracerOption) trace.Tracer {
	return otel.GetTracerProvider().Tracer(traceName, opts...)
}

func Meter(meterName string, opts ...metric.MeterOption) metric.Meter {
	return otel.GetMeterProvider().Meter(meterName, opts...)
}
