package server

import (
	"GRPC_Server/internal/kafka"
	"GRPC_Server/internal/kafka/kafkaEntities"
	mock_kafka "GRPC_Server/internal/kafka/mocks"
	"GRPC_Server/internal/logger"
	mock_repository "GRPC_Server/internal/repository/mocks"
	"GRPC_Server/internal/server/serviceEntities"
	"GRPC_Server/pkg/pkgErrors"
	"GRPC_Server/proto/gRPCService/studentProto"
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestStudentHandler_Get(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{
			Request:     "studentID:1",
			RequestType: "/StudentService/GetStudent",
			Timestamp:   time.Now().Round(time.Minute),
		}
		topic = "CRUD_events"
	)

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}

	logger.SetGlobal(zapLogger.With(zap.String("compenent", "compenent")))
	type mockExpected struct {
		result serviceEntities.StudentRequest
		error  error
	}
	//
	tests := []struct {
		description          string
		mockArguments        int64
		mockExpectedEntities mockExpected
		result               serviceEntities.StudentRequest
		expectedGrpcStatus   *status.Status
		expectedMessage      *studentProto.StudentGetResponse
		repoFunc             func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:          "Student not found",
			mockArguments:        4,
			mockExpectedEntities: mockExpected{error: pkgErrors.ErrNotFound},
			result:               serviceEntities.StudentRequest{},
			expectedGrpcStatus:   status.New(codes.NotFound, "Student with such student_id: Not Found"),
			expectedMessage:      nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
		{
			description:   "Student exists",
			mockArguments: 1,
			mockExpectedEntities: mockExpected{
				result: serviceEntities.StudentRequest{
					StudentID:   1,
					StudentName: "Test",
					Grade:       90,
				}, error: nil,
			},
			result:             serviceEntities.StudentRequest{StudentID: 1, StudentName: "Test", Grade: 90},
			expectedGrpcStatus: nil,
			expectedMessage:    &studentProto.StudentGetResponse{Message: "Successfully Got Student"},
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
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)

			grpcServer := grpc.NewServer()
			studentService := NewStudentServer(mockRepo, mockKafka)
			studentProto.RegisterStudentServiceServer(grpcServer, studentService)
			mockRepo.EXPECT().GetByID(
				gomock.Any(),
				tc.mockArguments,
			).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)

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
			client := studentProto.NewStudentServiceClient(conn)
			resp, err := client.GetStudent(context.Background(), &studentProto.StudentGetRequest{
				StudentID: tc.mockArguments,
			})
			// assert
			if err != nil {
				statusError, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedGrpcStatus.Code(), statusError.Code())
				assert.Equal(t, tc.expectedGrpcStatus.Err(), statusError.Err())

				return
			}
			assert.Equal(t, tc.expectedMessage.Message, resp.Message)
			assert.Equal(t, tc.result.StudentID, resp.Student.StudentID)
			assert.Equal(t, tc.result.StudentName, resp.Student.StudentName)
			assert.Equal(t, tc.result.Grade, resp.Student.Grade)
		})
	}
}

func TestStudentHandler_Create(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{
			Request:     "studentName:\"Test\"  grade:90",
			RequestType: "/StudentService/CreateStudent",
			Timestamp:   time.Now().Round(time.Minute),
		}
		topic = "CRUD_events"
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
		mockArguments        serviceEntities.StudentRequest
		mockExpectedEntities mockExpected
		result               serviceEntities.StudentRequest
		expectedGrpcStatus   *status.Status
		expectedMessage      *studentProto.StudentCreateResponse
		repoFunc             func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:          "Successfully Added into Database",
			mockArguments:        serviceEntities.StudentRequest{StudentID: 0, StudentName: "Test", Grade: 90},
			mockExpectedEntities: mockExpected{result: 1, error: nil},
			result:               serviceEntities.StudentRequest{StudentID: 1, StudentName: "Test", Grade: 90},
			expectedGrpcStatus:   nil,
			expectedMessage:      &studentProto.StudentCreateResponse{Message: "Successfully Created Student"},
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				m := mock_kafka.NewMockProducerInterface(ctrl)
				m.EXPECT().SendMessage(topic, request).Return(nil)
				m.EXPECT().Topic().Return(topic)
				return m
			},
		},
		{
			description:          "Failed database unable to add",
			mockArguments:        serviceEntities.StudentRequest{StudentID: 0, StudentName: "test", Grade: 98},
			mockExpectedEntities: mockExpected{result: 1, error: assert.AnError},
			result:               serviceEntities.StudentRequest{},
			expectedGrpcStatus:   status.New(codes.Internal, "Failed to add student: assert.AnError general error for testing"),
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
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)

			grpcServer := grpc.NewServer()
			studentService := NewStudentServer(mockRepo, mockKafka)
			studentProto.RegisterStudentServiceServer(grpcServer, studentService)

			mockRepo.EXPECT().Add(
				gomock.Any(),
				tc.mockArguments,
			).Return(tc.mockExpectedEntities.result, tc.mockExpectedEntities.error)

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

			client := studentProto.NewStudentServiceClient(conn)
			resp, err := client.CreateStudent(
				context.Background(),
				&studentProto.StudentCreateRequest{
					StudentName: tc.mockArguments.StudentName,
					Grade:       tc.mockArguments.Grade,
				},
			)
			// assert
			if err != nil {
				statusError, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedGrpcStatus.Code(), statusError.Code())

				return
			}
			assert.Equal(t, tc.expectedMessage.Message, resp.Message)
			assert.Equal(t, tc.result.StudentName, resp.Student.StudentName)
			assert.Equal(t, tc.result.Grade, resp.Student.Grade)
		})
	}
}

func TestStudentHandler_Delete(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{
			Request:     "studentID:1",
			RequestType: "/StudentService/DeleteStudent",
			Timestamp:   time.Now().Round(time.Minute),
		}
		topic = "CRUD_events"
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
		expectedMessage    *studentProto.StudentDeleteResponse
		repoFunc           func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:       "Unable to delete",
			mockArguments:     4,
			mockExpectedError: assert.AnError,
			expectedGrpcStatus: status.New(
				codes.Internal,
				"Failed to delete studentByID: assert.AnError general error for testing",
			),
			expectedMessage: nil,
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				// Return a mock producer without setting any expectations
				return mock_kafka.NewMockProducerInterface(ctrl)
			},
		},
		{
			description:        "Student not found",
			mockArguments:      4,
			mockExpectedError:  pkgErrors.ErrNotFound,
			expectedGrpcStatus: status.New(codes.NotFound, "Student with such student_id: Not Found"),
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
			expectedMessage:    &studentProto.StudentDeleteResponse{Message: "Successfully Deleted Student"},
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
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			mockKafka := tc.repoFunc(ctrl, topic, request)
			grpcServer := grpc.NewServer()
			studentService := NewStudentServer(mockRepo, mockKafka)
			studentProto.RegisterStudentServiceServer(grpcServer, studentService)
			mockRepo.EXPECT().Delete(gomock.Any(), tc.mockArguments).Return(tc.mockExpectedError)
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
			client := studentProto.NewStudentServiceClient(conn)
			// assert
			resp, err := client.DeleteStudent(
				context.Background(),
				&studentProto.StudentDeleteRequest{
					StudentID: tc.mockArguments,
				},
			)
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

func TestStudentHandler_Update(t *testing.T) {
	t.Parallel()
	var (
		request = kafkaEntities.Message{
			Request:     "studentName:\"Test2\" grade:92",
			RequestType: "/StudentService/UpdateStudent",
			Timestamp:   time.Now().Round(time.Minute),
		}
		topic = "CRUD_events"
	)

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}

	logger.SetGlobal(zapLogger.With(zap.String("compenent", "compenent")))

	tests := []struct {
		description        string
		mockArguments      serviceEntities.StudentRequest
		mockExpectedError  error
		expectedGrpcStatus *status.Status
		expectedMessage    *studentProto.StudentUpdateResponse
		repoFunc           func(*gomock.Controller, string, kafkaEntities.Message) kafka.ProducerInterface
	}{
		{
			description:        "Succesfully Updated in Database",
			mockArguments:      serviceEntities.StudentRequest{StudentID: 12, StudentName: "Test2", Grade: 92},
			mockExpectedError:  nil,
			expectedGrpcStatus: nil,
			expectedMessage:    &studentProto.StudentUpdateResponse{Message: "Successfully Updated Student"},
			repoFunc: func(ctrl *gomock.Controller, topic string, request kafkaEntities.Message) kafka.ProducerInterface {
				m := mock_kafka.NewMockProducerInterface(ctrl)
				m.EXPECT().SendMessage(topic, gomock.Any()).Return(nil)
				m.EXPECT().Topic().Return(topic)
				return m
			},
		},
		{
			description:        "Not Found",
			mockArguments:      serviceEntities.StudentRequest{StudentID: 21, StudentName: "Test2", Grade: 92},
			mockExpectedError:  pkgErrors.ErrNotFound,
			expectedGrpcStatus: status.New(codes.NotFound, "Student with such student_id: Not Found"),
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
			mockRepo := mock_repository.NewMockStudentPgRepo(ctrl)
			require.NoError(t, err)
			request = kafkaEntities.Message{
				Request:     "studentName:\"Test2\" grade:92",
				RequestType: "/StudentService/UpdateStudent",
				Timestamp:   time.Now().Round(time.Minute),
			}
			mockKafka := tc.repoFunc(ctrl, topic, request)
			grpcServer := grpc.NewServer()
			studentService := NewStudentServer(mockRepo, mockKafka)
			studentProto.RegisterStudentServiceServer(grpcServer, studentService)
			mockRepo.EXPECT().Update(gomock.Any(), tc.mockArguments.StudentID, tc.mockArguments).Return(tc.mockExpectedError)
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
			client := studentProto.NewStudentServiceClient(conn)
			resp, err := client.UpdateStudent(
				context.Background(),
				&studentProto.StudentUpdateRequest{
					StudentID:   tc.mockArguments.StudentID,
					StudentName: tc.mockArguments.StudentName,
					Grade:       tc.mockArguments.Grade,
				},
			)
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
