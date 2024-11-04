package store

import (
	"context"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepo struct {
	dbConn interface{}
}

func NewUserRepo(dbConn interface{}) *UserRepo {
	return &UserRepo{
		dbConn,
	}
}

func (repo *UserRepo) WhoIsCool(ctx context.Context, data entity.User) (*entity.User, error) {
	data.Id = uuid.NewString()
	return &data, nil
}
