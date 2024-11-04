package handler

import (
	"github.com/artem-webdev/otel_demo/internal/controller/http_ctrl/receivers"
	"github.com/artem-webdev/otel_demo/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/metric"
	"net/http"
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
	userUseCase UserUseCase
	tracer      TraceProvider
	meter       MetricProvider
}

func NewUserHandler(userUseCase UserUseCase, tracer TraceProvider, meter MetricProvider) *UserHandler {
	return &UserHandler{
		userUseCase,
		tracer,
		meter,
	}
}

func (h *UserHandler) WhoIsCool(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "http.UserHandler.WhoIsCool")
	defer span.End()
	demoCount, err := h.meter.Int64Counter("demo_http_request_counter", metric.WithDescription("demo counter"))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err:": err.Error(),
		})
	}
	demoCount.Add(ctx, 1)
	ret, err := h.userUseCase.WhoIsCool(ctx, dataMok)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err:": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(receivers.UserFromEntityUser(ret))
}
