package receivers

import (
	"github.com/artem-webdev/otel_demo/internal/controller/http_ctrl/dto"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
)

func UserFromEntityUser(data *entity.User) *dto.User {
	return (*dto.User)(data)
}
