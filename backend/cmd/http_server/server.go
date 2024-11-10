package http_server

import (
	"context"
	"errors"
	"github.com/artem-webdev/otel_demo/internal/controller/http_ctrl/handler"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
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
		return errors.New("nil context in NewHttpServer.run ")
	}
	s.route()
	go func() {
		select {
		case <-ctx.Done():
			ctx.Err().Error()
			if err := s.App.ShutdownWithContext(ctx); err != nil {
				if !strings.Contains(err.Error(), "context") {
					log.Printf("failed to shutdown shutdownMeterProvider: %s", err)
				}
			}
		}
	}()
	if err := s.App.Listen(addr); err != nil {
		return err
	}
	return nil
}
