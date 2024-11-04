package use_case

import (
	"context"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
)

type UserUseCase struct {
	userStore UserStoreProvider
	tracer    TraceProvider
}

func NewUserUseCase(userStore UserStoreProvider, tracer TraceProvider) *UserUseCase {
	return &UserUseCase{
		userStore,
		tracer,
	}
}

func (us *UserUseCase) WhoIsCool(ctx context.Context, data entity.User) (*entity.User, error) {
	ctx, span := us.tracer.Start(ctx, "UserUseCase.WhoIsCool")
	defer span.End()
	return us.userStore.WhoIsCool(ctx, data)
}
