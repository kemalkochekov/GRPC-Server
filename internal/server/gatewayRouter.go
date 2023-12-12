package server

import (
	"GRPC_Server/internal/logger"
	"GRPC_Server/proto/gRPCService/classInfoProto"
	"GRPC_Server/proto/gRPCService/studentProto"
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGatewayRouter(ctx context.Context, grpcServerEndpoint string, httpPort string) {
	loggerGateway := logger.FromContext(ctx)
	logger.ToContext(ctx, loggerGateway.With(zap.String("router", "gateway")))

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := studentProto.RegisterStudentServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		logger.Fatalf(ctx, "Failed to register Student HTTP server: %v", err)
	}
	err = classInfoProto.RegisterClassInfoServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcServerEndpoint, opts)
	if err != nil {
		logger.Fatalf(ctx, "Failed to register ClassInfo HTTP server: %v", err)
	}

	// Start HTTP server (and proxy calls to gRPC server)
	logger.Infof(ctx, "Starting HTTP server on port %v", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Fatalf("Failed to serve HTTP server over port %s: %v", httpPort, err)
	}
}
