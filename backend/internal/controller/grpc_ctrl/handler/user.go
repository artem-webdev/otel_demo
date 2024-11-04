package handler

import (
	"context"
	userpb "github.com/artem-webdev/otel_demo/internal/controller/grpc_ctrl/pb/user"
	"github.com/artem-webdev/otel_demo/internal/controller/grpc_ctrl/receivers"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	dataMok = entity.User{
		FirstName: "Johnny",
		LastName:  "Cage",
		Email:     "mortal-kombat-ultimate@gmail.com",
		Age:       30,
	}
)

type UserHandler struct {
	userpb.UnimplementedUserServer
	userUseCase UserUseCase
	tracer      TraceProvider
	meter       MetricProvider
}

func NewUserHandler(userUseCase UserUseCase, tracer TraceProvider, meter MetricProvider) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		tracer:      tracer,
		meter:       meter,
	}
}

func (h *UserHandler) WhoIsCool(ctx context.Context, req *emptypb.Empty) (*userpb.UserResponseMessage, error) {
	ctx, span := h.tracer.Start(ctx, "grpc.UserHandler.WhoIsCool")
	defer span.End()
	demoCount, err := h.meter.Int64Counter("demo_grpc_request_counter", metric.WithDescription("demo counter"))
	if err != nil {
		return nil, err
	}
	demoCount.Add(ctx, 1)
	ret, err := h.userUseCase.WhoIsCool(ctx, dataMok)
	if err != nil {
		return nil, err
	}
	return receivers.UserFromEntityUser(ret), nil
}
