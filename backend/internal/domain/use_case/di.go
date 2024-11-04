package use_case

import (
	"context"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
	"go.opentelemetry.io/otel/trace"
)

type TraceProvider interface {
	trace.Tracer
}

type UserStoreProvider interface {
	WhoIsCool(ctx context.Context, data entity.User) (*entity.User, error)
}
