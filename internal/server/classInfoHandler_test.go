package server

import (
	"GRPC_Server/internal/kafka"
	"GRPC_Server/internal/kafka/kafkaEntities"
	mock_kafka "GRPC_Server/internal/kafka/mocks"
	"GRPC_Server/internal/logger"
	mock_repository "GRPC_Server/internal/repository/mocks"
	"GRPC_Server/internal/server/serviceEntities"
	"GRPC_Server/pkg/pkgErrors"
	"GRPC_Server/proto/gRPCService/classInfoProto"
	"context"
	"log"
	"net"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"testing"
	"time"
)

func TestClassInfoHandler_AddClass(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{Request: "studentID:1  className:\"math\"", RequestType: "/ClassInfoService/AddClass", Timestamp: time.Now().Round(time.Minute)}
		topic   = "CRUD_events"
	)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}
	logger.SetGlobal(zapLogger.With(zap.String("compenent", "compenent")))
	type mockExpected struct {
		result int64
		error  error
	}
	tests := []struct {
		description          string
		mockArguments        serviceEntities.ClassInfo
		mockExpectedEntities mockExpected
		result               serviceEntities.ClassInfo
		expectedGrpcStatus   *status.Status
		expectedMessage      *classInfoProto.ClassAddResponse
		repoFunc             func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:          "Succesfully Added into Database",
			mockArguments:        serviceEntities.ClassInfo{StudentID: 1, ClassName: "math"},
			mockExpectedEntities: mockExpected{result: 0, error: nil},
			result:               serviceEntities.ClassInfo{StudentID: 1, ClassName: "math"},
			expectedGrpcStatus:   nil,
			expectedMessage:      &classInfoProto.ClassAddResponse{Message: "Successfully Added ClassInfo"},
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				m := mock_kafka.NewMockProducerInterface(ctrl)
				m.EXPECT().SendMessage(topic, request).Return(nil)
				m.EXPECT().Topic().Return(topic)
				return m
			},
		},
		{
			description:          "ForeignKey Error",
			mockArguments:        serviceEntities.ClassInfo{StudentID: 2, ClassName: "math"},
			mockExpectedEntities: mockExpected{result: -1, error: pkgErrors.ErrForeignKey},
			result:               serviceEntities.ClassInfo{},
			expectedGrpcStatus:   status.New(codes.Internal, "Failed to add class_info: ERROR: insert or update on table \"class_info\" violates foreign key constraint \"fk_student\" (SQLSTATE 23503)"),
			expectedMessage:      nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)
			grpcServer := grpc.NewServer()
			classInfoService := NewClassInfoServer(mockRepo, mockKafka)
			classInfoProto.RegisterClassInfoServiceServer(grpcServer, classInfoService)

			mockRepo.EXPECT().Add(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)
			defer ctrl.Finish()

			lis, err := net.Listen("tcp", "localhost:0")
			require.NoError(t, err)
			go func() {
				err := grpcServer.Serve(lis)
				require.NoError(t, err)
			}()
			defer grpcServer.Stop()

			// Create a gRPC client
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()
			client := classInfoProto.NewClassInfoServiceClient(conn)
			resp, err := client.AddClass(context.Background(), &classInfoProto.ClassAddRequest{StudentID: tc.mockArguments.StudentID, ClassName: tc.mockArguments.ClassName})
			// assert
			if err != nil {
				statusError, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedGrpcStatus.Err(), statusError.Err())
				assert.Equal(t, tc.expectedGrpcStatus.Code(), statusError.Code())
				return
			}
			assert.Equal(t, tc.expectedMessage.Message, resp.Message)
			assert.Equal(t, tc.result.StudentID, resp.ClassInfo.StudentID)
			assert.Equal(t, tc.result.ClassName, resp.ClassInfo.ClassName)
		})
	}
}

func TestClassInfoHandler_GetAllClassesByStudent(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{Request: "studentID:1", RequestType: "/ClassInfoService/GetAllClassesByStudent", Timestamp: time.Now().Round(time.Minute)}
		topic   = "CRUD_events"
	)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}
	logger.SetGlobal(zapLogger.With(zap.String("compenent", "compenent")))

	type mockExpected struct {
		result []serviceEntities.ClassInfo
		error  error
	}
	tests := []struct {
		description          string
		mockArguments        int64
		mockExpectedEntities mockExpected
		result               []serviceEntities.ClassInfo
		expectedGrpcStatus   *status.Status
		expectedMessage      *classInfoProto.ClassesGetResponse
		repoFunc             func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:          "Succesfully Get ClassInfo By StudentID",
			mockArguments:        1,
			mockExpectedEntities: mockExpected{result: []serviceEntities.ClassInfo{{StudentID: 1, ClassName: "math"}}, error: nil},
			result:               []serviceEntities.ClassInfo{{StudentID: 1, ClassName: "math"}},
			expectedGrpcStatus:   nil,
			expectedMessage:      &classInfoProto.ClassesGetResponse{Message: "Successfully Got All ClassInfo by StudentID"},
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				m := mock_kafka.NewMockProducerInterface(ctrl)
				m.EXPECT().SendMessage(topic, request).Return(nil)
				m.EXPECT().Topic().Return(topic)
				return m
			},
		},
		{
			description:          "In ClassInfo with StudentID does not exist",
			mockArguments:        2,
			mockExpectedEntities: mockExpected{},
			result:               []serviceEntities.ClassInfo{},
			expectedGrpcStatus:   status.New(codes.NotFound, "No existing student in class_info"),
			expectedMessage:      nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)
			grpcServer := grpc.NewServer()
			classInfoService := NewClassInfoServer(mockRepo, mockKafka)
			classInfoProto.RegisterClassInfoServiceServer(grpcServer, classInfoService)

			mockRepo.EXPECT().GetByStudentID(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)
			defer ctrl.Finish()

			lis, err := net.Listen("tcp", "localhost:0")
			require.NoError(t, err)
			go func() {
				err := grpcServer.Serve(lis)
				require.NoError(t, err)
			}()
			defer grpcServer.Stop()

			// Create a gRPC client
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()
			client := classInfoProto.NewClassInfoServiceClient(conn)
			resp, err := client.GetAllClassesByStudent(context.Background(), &classInfoProto.ClassesGetRequest{StudentID: tc.mockArguments})
			// act
			if err != nil {
				statusError, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedGrpcStatus.Code(), statusError.Code())
				assert.Equal(t, tc.expectedGrpcStatus.Err(), statusError.Err())
				return
			}
			assert.Equal(t, tc.expectedMessage.Message, resp.Message)
			assert.Equal(t, tc.result[0].StudentID, resp.ClassInfo[0].StudentID)
			assert.Equal(t, tc.result[0].ClassName, resp.ClassInfo[0].ClassName)
		})
	}
}

func TestClassInfoHandler_DeleteClassByStudent(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{Request: "studentID:1", RequestType: "/ClassInfoService/DeleteClassByStudent", Timestamp: time.Now().Round(time.Minute)}
		topic   = "CRUD_events"
	)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}
	logger.SetGlobal(zapLogger.With(zap.String("compenent", "compenent")))
	tests := []struct {
		description        string
		mockArguments      int64
		mockExpectedError  error
		expectedGrpcStatus *status.Status
		expectedMessage    *classInfoProto.ClassDeleteResponse
		repoFunc           func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:        "Unable to delete",
			mockArguments:      4,
			mockExpectedError:  assert.AnError,
			expectedGrpcStatus: status.New(codes.Internal, "Failed to delete record from class_info by StudentByID: assert.AnError general error for testing"),
			expectedMessage:    nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
		{
			description:        "In ClassInfo not found StudentByID",
			mockArguments:      4,
			mockExpectedError:  pkgErrors.ErrNotFound,
			expectedGrpcStatus: status.New(codes.NotFound, "Cannot delete the class_info due to non existing studentID Not Found"),
			expectedMessage:    nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
		{
			description:        "Succesfully Deleted StudentID",
			mockArguments:      1,
			mockExpectedError:  nil,
			expectedGrpcStatus: nil,
			expectedMessage:    &classInfoProto.ClassDeleteResponse{Message: "Successfully Deleted ClassInfo By StudentID"},
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				m := mock_kafka.NewMockProducerInterface(ctrl)
				m.EXPECT().SendMessage(topic, request).Return(nil)
				m.EXPECT().Topic().Return(topic)
				return m
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)
			grpcServer := grpc.NewServer()
			classInfoService := NewClassInfoServer(mockRepo, mockKafka)
			classInfoProto.RegisterClassInfoServiceServer(grpcServer, classInfoService)

			mockRepo.EXPECT().DeleteClassByStudentID(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedError)
			defer ctrl.Finish()
			lis, err := net.Listen("tcp", "localhost:0")
			require.NoError(t, err)
			go func() {
				err := grpcServer.Serve(lis)
				require.NoError(t, err)
			}()
			defer grpcServer.Stop()

			// Create a gRPC client
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()
			client := classInfoProto.NewClassInfoServiceClient(conn)
			resp, err := client.DeleteClassByStudent(context.Background(), &classInfoProto.ClassDeleteRequest{StudentID: tc.mockArguments})
			if err != nil {
				statusError, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedGrpcStatus.Code(), statusError.Code())
				assert.Equal(t, tc.expectedGrpcStatus.Err(), statusError.Err())
				return
			}
			// assert
			assert.Equal(t, tc.expectedMessage.Message, resp.Message)
		})
	}
}

func TestClassInfoHandler_UpdateClass(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{Request: "studentID:1 className:\"math\"", RequestType: "/ClassInfoService/UpdateClass", Timestamp: time.Now().Round(time.Minute)}
		topic   = "CRUD_events"
	)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}
	logger.SetGlobal(zapLogger.With(zap.String("compenent", "compenent")))

	tests := []struct {
		description        string
		mockArguments      serviceEntities.ClassInfo
		mockExpectedError  error
		expectedGrpcStatus *status.Status
		expectedMessage    *classInfoProto.ClassUpdateResponse
		repoFunc           func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:        "Succesfully Updated in Database",
			mockArguments:      serviceEntities.ClassInfo{StudentID: 1, ClassName: "math"},
			mockExpectedError:  nil,
			expectedGrpcStatus: nil,
			expectedMessage:    &classInfoProto.ClassUpdateResponse{Message: "Successfully Updated ClassInfo By StudentID"},
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				m := mock_kafka.NewMockProducerInterface(ctrl)
				m.EXPECT().SendMessage(topic, gomock.Any()).Return(nil)
				m.EXPECT().Topic().Return(topic)
				return m
			},
		},
		{
			description:        "Not Found",
			mockArguments:      serviceEntities.ClassInfo{StudentID: 2, ClassName: "math"},
			mockExpectedError:  pkgErrors.ErrNotFound,
			expectedGrpcStatus: status.New(codes.NotFound, "Cannot update the student in class_info due to existing references (foreign key constraint). Not Found"),
			expectedMessage:    nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockRepo := mock_repository.NewMockClassInfoPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)
			grpcServer := grpc.NewServer()
			classInfoService := NewClassInfoServer(mockRepo, mockKafka)
			classInfoProto.RegisterClassInfoServiceServer(grpcServer, classInfoService)

			mockRepo.EXPECT().Update(gomock.Any(), tc.mockArguments.StudentID, tc.mockArguments).Return(tc.mockExpectedError)
			defer ctrl.Finish()
			// act
			lis, err := net.Listen("tcp", "localhost:0")
			require.NoError(t, err)
			go func() {
				err := grpcServer.Serve(lis)
				require.NoError(t, err)
			}()
			defer grpcServer.Stop()

			// Create a gRPC client
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()
			client := classInfoProto.NewClassInfoServiceClient(conn)

			resp, err := client.UpdateClass(context.Background(), &classInfoProto.ClassUpdateRequest{StudentID: tc.mockArguments.StudentID, ClassName: tc.mockArguments.ClassName})
			// assert
			if err != nil {
				statusError, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedGrpcStatus.Code(), statusError.Code())
				assert.Equal(t, tc.expectedGrpcStatus.Err(), statusError.Err())
				return
			}
			assert.Equal(t, tc.expectedMessage.Message, resp.Message)
		})
	}
}
