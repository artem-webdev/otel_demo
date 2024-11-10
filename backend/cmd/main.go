package main

import (
	"context"
	"github.com/artem-webdev/otel_demo/cmd/http_server"
	"github.com/artem-webdev/otel_demo/internal/adapter/store"
	grpchandler "github.com/artem-webdev/otel_demo/internal/controller/grpc_ctrl/handler"
	userpb "github.com/artem-webdev/otel_demo/internal/controller/grpc_ctrl/pb/user"
	httphandler "github.com/artem-webdev/otel_demo/internal/controller/http_ctrl/handler"
	"github.com/artem-webdev/otel_demo/internal/domain/use_case"
	"github.com/artem-webdev/otel_demo/internal/pkg/otel_sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	HttpServerAddr = "0.0.0.0:7777"
	GrpcServerAddr = "0.0.0.0:50057"
	TracerName     = "otel-collector"
	MeterName      = "otel-collector"
	ServiceName    = "otel-demo-service"
	//AddrGrpcCollector = "otel-collector:4320"
	AddrGrpcCollector = "0.0.0.0:4320"
)

func tracerInit(ctx context.Context) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	exporterGrpc, err := otel_sdk.GrpcTraceExporter(ctx, AddrGrpcCollector, opts...)
	if err != nil {
		log.Fatal(err)
	}
	shutdownTracerProvider, err := otel_sdk.SetTracerProvider(ServiceName, exporterGrpc)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		select {
		case <-ctx.Done():
			if err = shutdownTracerProvider(ctx); err != nil {
				if !strings.Contains(err.Error(), "context") {
					log.Printf("failed to shutdown shutdownMeterProvider: %s", err)
				}
			}
		}
	}()
}

func metricInit(ctx context.Context) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	exporterGrpc, err := otel_sdk.GrpcMetricExporter(ctx, AddrGrpcCollector, opts...)
	if err != nil {
		log.Fatal(err)
	}
	shutdownMeterProvider, err := otel_sdk.SetMetricProvider(ServiceName, exporterGrpc)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		select {
		case <-ctx.Done():
			if err = shutdownMeterProvider(ctx); err != nil {
				if !strings.Contains(err.Error(), "context") {
					log.Printf("failed to shutdown shutdownMeterProvider: %s", err)
				}
			}
		}
	}()
}

func main() {
	ctxParent := context.Background()
	// set tracer
	tracerInit(ctxParent)
	tracer := otel_sdk.Tracer(TracerName)
	// set metrics
	metricInit(ctxParent)
	meter := otel_sdk.Meter(MeterName)
	// init di
	repo := store.NewUserRepo(nil)
	userUseCase := use_case.NewUserUseCase(repo, tracer)
	wg := &sync.WaitGroup{}
	// start http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		userHandlerHttp := httphandler.NewUserHandler(userUseCase, tracer, meter)
		server := http_server.NewHttpServer(userHandlerHttp)
		if err := server.Run(ctxParent, HttpServerAddr); err != nil {
			log.Fatal(err)
		}
	}()
	// start grpc server
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		grpcServer := grpc.NewServer()
		go func() {
			select {
			case <-ctx.Done():
				grpcServer.GracefulStop()
				log.Println("gRPC server stopped gracefully")
			}
		}()
		userHandlerGrpc := grpchandler.NewUserHandler(userUseCase, tracer, meter)
		userpb.RegisterUserServer(grpcServer, userHandlerGrpc)
		reflection.Register(grpcServer)
		lis, err := net.Listen("tcp", GrpcServerAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}(ctxParent)
	wg.Wait()
}
