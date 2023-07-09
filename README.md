#gRPC Worker

The grpcworker package provides a worker implementation for running a gRPC server with unary and stream interceptors. It is designed to be flexible and customizable through the use of options.

## Usage

```go
package main

import (
	"fmt"
	"net"

	"github.com/voi-oss/svc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/alesr/grpcworker"
)

type MyApp struct{}

// Implement methods of your gRPC service interface

func main() {
    // Create a logger
    logger, _ := zap.NewDevelopment()
    
    // Create a listener on a specific address and port
    lis, _ := net.Listen("tcp", ":50051")

    // Create your application instance
    app := MyApp{}

    // Create a gRPC service descriptor
    serviceDesc := grpc.ServiceDesc{
        ServiceName: "service-name",
        HandlerType: (*MyApp)(nil),
    }

    // Create your unary interceptor
    yourUnaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Do something before handling the request
        resp, _ := handler(ctx, req)
        // Do something after handling the request
        return resp, nil
    }

    // Create your stream interceptor
    yourStreamInterceptor := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
        // Do something before handling the request
        _ := handler(srv, ss)
        // Do something after handling the request
        return nil
    }

    // Create a new worker with optional interceptors
    worker := grpcworker.New(logger, lis, &app, &serviceDesc,
        grpcworker.WithUnaryInterceptor(yourUnaryInterceptor),
        grpcworker.WithStreamInterceptor(yourStreamInterceptor),
    )

    // Initialize the worker
    if err := worker.Init(logger); err != nil {
        logger.Fatal("failed to initialize worker", zap.Error(err))
    }

    // Start the gRPC server
    if err := worker.Run(); err != nil {
        logger.Fatal("failed to start server", zap.Error(err))
    }

    // Terminate the server gracefully
    if err := worker.Terminate(); err != nil {
        logger.Fatal("failed to terminate server", zap.Error(err))
    }
}
```
