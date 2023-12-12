package server

import (
	"GRPC_Server/internal/kafka"
	"GRPC_Server/internal/kafka/kafkaEntities"
	"GRPC_Server/internal/logger"
	"GRPC_Server/internal/repository"
	"GRPC_Server/internal/server/serviceEntities"
	"GRPC_Server/pkg/pkgErrors"
	"GRPC_Server/proto/gRPCService/classInfoProto"
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClassInfoServer struct {
	classInfoStorage repository.ClassInfoPgRepo
	producer         kafka.ProducerInterface
	classInfoProto.UnimplementedClassInfoServiceServer
}

func NewClassInfoServer(classInfoStorage repository.ClassInfoPgRepo, producer kafka.ProducerInterface) *ClassInfoServer {
	return &ClassInfoServer{
		classInfoStorage: classInfoStorage,
		producer:         producer,
	}
}

// curl -X POST localhost:9000/class_info -d '{"student_id":1,"class_name":"math"}' -i
func (s *ClassInfoServer) AddClass(ctx context.Context, req *classInfoProto.ClassAddRequest) (*classInfoProto.ClassAddResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "ClassInfo Create")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "AddClass")
	defer span.Finish()

	var classInfo serviceEntities.ClassInfo
	classInfo.StudentID = req.StudentID
	classInfo.ClassName = req.ClassName
	var err error
	classInfo.ID, err = s.classInfoStorage.Add(ctx, classInfo)
	if err != nil {
		logError(ctx, span, err, "Failed to add class_info:")
		return nil, status.Errorf(codes.Internal, "Failed to add class_info: %v", err)
	}
	// Created Kafka Message
	method, _ := grpc.Method(ctx)
	err = s.producer.SendMessage(s.producer.Topic(), kafkaEntities.Message{
		RequestType: method,
		Request:     req.String(),
		Timestamp:   time.Now().Round(time.Minute),
	})
	if err != nil {
		logError(ctx, span, err, "Failed to send Kafka message: ")
	}
	classInfoProtoRes := &classInfoProto.ClassInfo{
		Id:        classInfo.ID,
		StudentID: classInfo.StudentID,
		ClassName: classInfo.ClassName,
	}
	return &classInfoProto.ClassAddResponse{ClassInfo: classInfoProtoRes, Message: "Successfully Added ClassInfo"}, nil
}

// curl -X PUT localhost:9000/class_info -d '{"student_id":4,"class_name":"phys"}' -i
func (s *ClassInfoServer) UpdateClass(ctx context.Context, req *classInfoProto.ClassUpdateRequest) (*classInfoProto.ClassUpdateResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "ClassInfo Update")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateClass")
	defer span.Finish()

	classInfo := serviceEntities.ClassInfo{
		StudentID: req.StudentID,
		ClassName: req.ClassName,
	}
	var err error
	err = s.classInfoStorage.Update(ctx, req.StudentID, classInfo)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			logError(ctx, span, err, "Cannot update the student in class_info due to existing references (foreign key constraint). ")
			return nil, status.Errorf(codes.NotFound, "Cannot update the student in class_info due to existing references (foreign key constraint). %v", err)
		}
		logError(ctx, span, err, "Failed to update class_info: ")
		return nil, status.Errorf(codes.Internal, "Failed to update class_info: %v", err)
	}
	method, _ := grpc.Method(ctx)
	err = s.producer.SendMessage(s.producer.Topic(), kafkaEntities.Message{
		RequestType: method,
		Request:     req.String(),
		Timestamp:   time.Now().Round(time.Minute),
	})
	if err != nil {
		logError(ctx, span, err, "Failed to send Kafka message: ")
	}
	return &classInfoProto.ClassUpdateResponse{Message: "Successfully Updated ClassInfo By StudentID"}, nil
}

// DeleteClassByStudent handles HTTP DELETE requests for deleting class information by student ID
func (s *ClassInfoServer) DeleteClassByStudent(ctx context.Context, req *classInfoProto.ClassDeleteRequest) (*classInfoProto.ClassDeleteResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "ClassInfo Delete")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteClassByStudent")
	defer span.Finish()

	err := s.classInfoStorage.DeleteClassByStudentID(ctx, req.StudentID)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			logError(ctx, span, err, "Cannot delete the class_info due to non existing studentID ")
			return nil, status.Errorf(codes.NotFound, "Cannot delete the class_info due to non existing studentID %v", err)
		}
		logError(ctx, span, err, "Failed to delete record from class_info by StudentByID: ")
		return nil, status.Errorf(codes.Internal, "Failed to delete record from class_info by StudentByID: %v", err)
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
	return &classInfoProto.ClassDeleteResponse{Message: "Successfully Deleted ClassInfo By StudentID"}, nil
}

// curl -X GET localhost:9000/class_info/4 -i
func (s *ClassInfoServer) GetAllClassesByStudent(ctx context.Context, req *classInfoProto.ClassesGetRequest) (*classInfoProto.ClassesGetResponse, error) {
	loggerGateway := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, loggerGateway.With(zap.String("method", "ClassInfo Get")))

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAllClassesByStudent")
	defer span.Finish()
	classesInfo, err := s.classInfoStorage.GetByStudentID(ctx, req.StudentID)
	if err != nil {
		logError(ctx, span, err, "Failed to get all records who is attending that class from class_info:")
		return nil, status.Errorf(codes.Internal, "Failed to get all records who is attending that class from class_info: %v", err)
	}

	if len(classesInfo) == 0 {
		// classesInfo is empty
		logError(ctx, span, err, "No existing student in class_info: ")
		return nil, status.Error(codes.NotFound, "No existing student in class_info")
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

	ClassesInfo := make([]*classInfoProto.ClassInfo, 0)

	for _, classInfo := range classesInfo {
		userInfo := &classInfoProto.ClassInfo{
			Id:        classInfo.ID,
			StudentID: classInfo.StudentID,
			ClassName: classInfo.ClassName,
		}
		ClassesInfo = append(ClassesInfo, userInfo)
	}
	return &classInfoProto.ClassesGetResponse{ClassInfo: ClassesInfo, Message: "Successfully Got All ClassInfo by StudentID"}, nil
}
