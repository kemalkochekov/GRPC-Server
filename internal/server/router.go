package server

import (
	"GRPC_Server/internal/kafka"
	"GRPC_Server/internal/logger"
	"GRPC_Server/internal/repository"
	"GRPC_Server/proto/gRPCService/classInfoProto"
	"GRPC_Server/proto/gRPCService/studentProto"
	"context"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
)

func logError(ctx context.Context, span opentracing.Span, err error, message string) {
	logger.Errorf(ctx, "%s: %v", message, err)
	span.SetTag("error", true)
	span.LogFields(log.Error(err))
}
func StartGRPCServer(studentStorage repository.StudentPgRepo, classInfoStorage repository.ClassInfoPgRepo, producer kafka.ProducerInterface, grpcPort string) error {

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	studentService := NewStudentServer(studentStorage, producer)
	classInfoService := NewClassInfoServer(classInfoStorage, producer)
	studentProto.RegisterStudentServiceServer(grpcServer, studentService)
	classInfoProto.RegisterClassInfoServiceServer(grpcServer, classInfoService)

	// Start grpc server
	return grpcServer.Serve(lis)
}
