package http_server

import (
	"context"
	"github.com/artem-webdev/otel_demo/internal/controller/http_ctrl/handler"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

type HttpServer struct {
	App         *fiber.App
	userHandler *handler.UserHandler
}

func NewHttpServer(userHandler *handler.UserHandler) *HttpServer {
	return &HttpServer{
		fiber.New(fiber.Config{
			Prefork:      false,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}),
		userHandler,
	}
}

// route mok routing
func (s *HttpServer) route() {
	s.App.Get("/user/cool/", s.userHandler.WhoIsCool)
}

func (s *HttpServer) Run(ctx context.Context, addr string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	s.route()
	go func() {
		if err := s.App.ShutdownWithContext(ctx); err != nil {
			log.Println(err)
		}
	}()
	if err := s.App.Listen(addr); err != nil {
		return err
	}
	return nil
}
