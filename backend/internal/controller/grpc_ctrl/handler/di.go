package handler

import (
	"context"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type MetricProvider interface {
	metric.Meter
}

type TraceProvider interface {
	trace.Tracer
}

type UserUseCase interface {
	WhoIsCool(ctx context.Context, data entity.User) (*entity.User, error)
}
