package receivers

import (
	userpb "github.com/artem-webdev/otel_demo/internal/controller/grpc_ctrl/pb/user"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
)

func UserFromEntityUser(data *entity.User) *userpb.UserResponseMessage {
	return &userpb.UserResponseMessage{
		Id:        data.Id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Age:       uint32(data.Age),
	}
}
