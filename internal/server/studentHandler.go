package server

import (
	"GRPC_Server/internal/kafka"
	"GRPC_Server/internal/kafka/kafkaEntities"
	"GRPC_Server/internal/logger"
	"GRPC_Server/internal/repository"
	"GRPC_Server/internal/server/serviceEntities"
	"GRPC_Server/pkg/pkgErrors"
	"GRPC_Server/proto/gRPCService/studentProto"
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StudentServer struct {
	studentRepo repository.StudentPgRepo
	producer    kafka.ProducerInterface
	studentProto.UnimplementedStudentServiceServer
}

func NewStudentServer(studentStorage repository.StudentPgRepo, producer kafka.ProducerInterface) *StudentServer {
	return &StudentServer{
		studentRepo: studentStorage,
		producer:    producer,
	}
}

func (s *StudentServer) Main(ctx context.Context, _ *empty.Empty) (*studentProto.MainResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "Main")))

	span, _ := opentracing.StartSpanFromContext(ctx, "Main")
	defer span.Finish()

	return &studentProto.MainResponse{Message: "Welcome TO GPRC SERVER"}, nil
}

func (s *StudentServer) CreateStudent(
	ctx context.Context,
	req *studentProto.StudentCreateRequest,
) (*studentProto.StudentCreateResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "Student Create")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateStudent")
	defer span.Finish()

	if req.StudentName == "" || req.Grade < 0 {
		logError(
			ctx,
			span,
			errors.New("Failed Student name is empty or Grade is negative"),
			"Failed Student name is empty or Grade is negative",
		)

		return nil, status.Errorf(codes.InvalidArgument, "Failed Student name is empty or Grade is negative")
	}
	var studentReq serviceEntities.StudentRequest
	studentReq.StudentName = req.StudentName
	studentReq.Grade = req.Grade
	var err error

	studentReq.StudentID, err = s.studentRepo.Add(ctx, studentReq)
	if err != nil {
		logError(ctx, span, err, "Failed to add student: ")
		return nil, status.Errorf(codes.Internal, "Failed to add student: %v", err)
	}
	method, _ := grpc.Method(ctx)
	// Created Kafka Message
	err = s.producer.SendMessage(s.producer.Topic(), kafkaEntities.Message{
		RequestType: method,
		Request:     req.String(),
		Timestamp:   time.Now().Round(time.Minute),
	})
	if err != nil {
		logError(ctx, span, err, "Failed to send Kafka message: ")
	}
	student := &studentProto.Student{
		StudentID:   studentReq.StudentID,
		StudentName: req.StudentName,
		Grade:       req.Grade,
	}

	return &studentProto.StudentCreateResponse{Student: student, Message: "Successfully Created Student"}, nil
}

func (s *StudentServer) UpdateStudent(
	ctx context.Context,
	req *studentProto.StudentUpdateRequest,
) (*studentProto.StudentUpdateResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "Student Update")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateStudent")
	defer span.Finish()

	student := serviceEntities.StudentRequest{
		StudentID:   req.StudentID,
		StudentName: req.StudentName,
		Grade:       req.Grade,
	} // 1

	err := s.studentRepo.Update(ctx, req.StudentID, student)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			logError(ctx, span, err, "Student with such student_id: ")
			return nil, status.Errorf(codes.NotFound, "Student with such student_id: %v", err)
		}

		logError(ctx, span, err, "Failed to update studentByID: ")

		return nil, status.Errorf(codes.Internal, "Failed to update studentByID: %v", err)
	}
	method, _ := grpc.Method(ctx)
	// Created Kafka Message
	err = s.producer.SendMessage(s.producer.Topic(), kafkaEntities.Message{
		RequestType: method,
		Request:     req.String(),
		Timestamp:   time.Now().Round(time.Minute),
	})
	if err != nil {
		logError(ctx, span, err, "Failed to send Kafka message: ")
	}

	return &studentProto.StudentUpdateResponse{Message: "Successfully Updated Student"}, nil
}

func (s *StudentServer) GetStudent(
	ctx context.Context,
	req *studentProto.StudentGetRequest,
) (*studentProto.StudentGetResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "Student Get")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetStudent")
	defer span.Finish()

	userInfo, err := s.studentRepo.GetByID(ctx, req.StudentID)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			logError(ctx, span, err, "Student with such student_id: ")
			return nil, status.Errorf(codes.NotFound, "Student with such student_id: %v", err)
		}

		logError(ctx, span, err, "Failed to get record by StudentByID: ")

		return nil, status.Errorf(codes.Internal, "Failed to get record by StudentByID: %v", err)
	}
	method, _ := grpc.Method(ctx)
	// Created Kafka Message
	err = s.producer.SendMessage(s.producer.Topic(), kafkaEntities.Message{
		RequestType: method,
		Request:     req.String(),
		Timestamp:   time.Now().Round(time.Minute),
	})
	if err != nil {
		logError(ctx, span, err, "Failed to send Kafka message: ")
	}
	student := &studentProto.Student{
		StudentID:   userInfo.StudentID,
		StudentName: userInfo.StudentName,
		Grade:       userInfo.Grade,
	}

	return &studentProto.StudentGetResponse{Student: student, Message: "Successfully Got Student"}, nil
}

func (s *StudentServer) DeleteStudent(
	ctx context.Context,
	req *studentProto.StudentDeleteRequest,
) (*studentProto.StudentDeleteResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "Student Delete")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteStudent")
	defer span.Finish()

	err := s.studentRepo.Delete(ctx, req.StudentID)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			logError(ctx, span, err, "Student with such student_id: ")
			return nil, status.Errorf(codes.NotFound, "Student with such student_id: %v", err)
		}

		logError(ctx, span, err, "Failed to delete studentByID: ")

		return nil, status.Errorf(codes.Internal, "Failed to delete studentByID: %v", err)
	}
	method, _ := grpc.Method(ctx)
	// Created Kafka Message
	err = s.producer.SendMessage(s.producer.Topic(), kafkaEntities.Message{
		RequestType: method,
		Request:     req.String(),
		Timestamp:   time.Now().Round(time.Minute),
	})
	if err != nil {
		logError(ctx, span, err, "Failed to send Kafka message: ")
	}

	return &studentProto.StudentDeleteResponse{Message: "Successfully Deleted Student"}, nil
}
