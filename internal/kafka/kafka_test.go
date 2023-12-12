//go:build integration
// +build integration

package kafka_test

import (
	kafka_test "GRPC_Server/internal/kafka"
	"GRPC_Server/internal/logger"
	"GRPC_Server/internal/repository"
	"GRPC_Server/internal/repository/postgres"
	"GRPC_Server/internal/server"
	"GRPC_Server/internal/server/serviceEntities"
	"GRPC_Server/proto/gRPCService/classInfoProto"
	"GRPC_Server/proto/gRPCService/studentProto"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/stretchr/testify/require"
)

// DB_HOST=localhost; DB_PORT=5432; DB_USER=postgres; DB_PASSWORD=test; DB_NAME=test
func setupGrpcServerAndClient(t *testing.T, registerService func(*grpc.Server)) *grpc.ClientConn {
	lis, err := net.Listen("tcp", "localhost:8080")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	registerService(grpcServer)

	go grpcServer.Serve(lis)
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		grpcServer.Stop()
		conn.Close()
		lis.Close()
	})
	return conn
}
func TestKafkaGet(t *testing.T) {

	var (
		topic      = "CRUD_events"
		brokers    = []string{"localhost:9092"}
		studentReq = serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		classInfoReq = serviceEntities.ClassInfo{
			StudentID: 1,
			ClassName: "MATH",
		}
		migrationPath = "./../repository/migrations"
	)
	setUpAll := func(t *testing.T) (*postgres.TDB, kafka_test.ProducerInterface, kafka_test.ConsumerInterface, *bytes.Buffer) {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to initialize Zap logger: %v\n", err)
		}
		logger.SetGlobal(zapLogger.With(zap.String("component", "main")))
		db := postgres.NewFromEnv()
		t.Cleanup(func() { db.DB.GetPool(context.Background()).Close() })

		producer, err := kafka_test.NewProducer(brokers, topic)
		require.NoError(t, err)

		var messages bytes.Buffer
		consumer, err := kafka_test.NewKafkaConsumer(brokers, &messages)
		require.NoError(t, err)

		db.SetUpDatabase(migrationPath)

		t.Cleanup(func() { db.TearDownDatabase(migrationPath) })
		return db, producer, consumer, &messages
	}
	testCases := []struct {
		name       string
		setup      func(*testing.T, *postgres.TDB, kafka_test.ProducerInterface)
		assertFunc func(*testing.T, *bytes.Buffer)
	}{
		{
			name: "Success Student Kafka Get",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				// Act
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls
				_, err := client.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = client.GetStudent(context.Background(), &studentProto.StudentGetRequest{StudentID: studentReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				fmt.Println(messagesSlice[0])
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /StudentService/GetStudent")
			},
		},
		{
			name: "Fail Student Kafka Get",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls
				_, err := client.GetStudent(context.Background(), &studentProto.StudentGetRequest{StudentID: studentReq.StudentID})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
		{
			name: "Success ClassInfo Kafka Get",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})
				studentClient := studentProto.NewStudentServiceClient(conn)
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)
				_, err := studentClient.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = classInfoClient.AddClass(context.Background(), &classInfoProto.ClassAddRequest{StudentID: classInfoReq.StudentID, ClassName: classInfoReq.ClassName})
				require.NoError(t, err)
				_, err = classInfoClient.GetAllClassesByStudent(context.Background(), &classInfoProto.ClassesGetRequest{StudentID: classInfoReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /ClassInfoService/AddClass")
				assert.Contains(t, messagesSlice[2], "Received Kafka Message: RequestType: /ClassInfoService/GetAllClassesByStudent")
			},
		},
		{
			name: "Fail ClassInfo Kafka Get",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)
				_, err := classInfoClient.GetAllClassesByStudent(context.Background(), &classInfoProto.ClassesGetRequest{StudentID: classInfoReq.StudentID})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, producer, consumer, messages := setUpAll(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := consumer.ReadMessages(ctx, topic); err != nil {
					err := producer.Close()
					require.NoError(t, err)
					err = consumer.ConsumerClose()
					require.NoError(t, err)
				}
			}()
			tc.setup(t, db, producer)
			wg.Wait()
			tc.assertFunc(t, messages)
		})
	}
}

func TestStudentKafkaCreate(t *testing.T) {
	var (
		topic      = "CRUD_events"
		brokers    = []string{"localhost:9092"}
		studentReq = serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		classInfoReq = serviceEntities.ClassInfo{
			StudentID: 1,
			ClassName: "MATH",
		}
		migrationPath = "./../repository/migrations"
	)

	setUpAll := func(t *testing.T) (*postgres.TDB, kafka_test.ProducerInterface, kafka_test.ConsumerInterface, *bytes.Buffer) {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to initialize Zap logger: %v\n", err)
		}
		logger.SetGlobal(zapLogger.With(zap.String("component", "main")))
		db := postgres.NewFromEnv()
		t.Cleanup(func() { db.DB.GetPool(context.Background()).Close() })

		producer, err := kafka_test.NewProducer(brokers, topic)
		require.NoError(t, err)

		var messages bytes.Buffer
		consumer, err := kafka_test.NewKafkaConsumer(brokers, &messages)
		require.NoError(t, err)

		db.SetUpDatabase(migrationPath)

		t.Cleanup(func() { db.TearDownDatabase(migrationPath) })
		return db, producer, consumer, &messages
	}
	testCases := []struct {
		name       string
		setup      func(*testing.T, *postgres.TDB, kafka_test.ProducerInterface)
		assertFunc func(*testing.T, *bytes.Buffer)
	}{
		{
			name: "Success Student Kafka Create",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls
				_, err := client.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = client.GetStudent(context.Background(), &studentProto.StudentGetRequest{StudentID: studentReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /StudentService/GetStudent")
			},
		},
		{
			name: "Fail Student Kafka Create",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls
				_, err := client.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: "",
					Grade:       studentReq.Grade,
				})
				require.Error(t, err)

			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
		{
			name: "Success ClassInfo Kafka Create",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})

				studentClient := studentProto.NewStudentServiceClient(conn)
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)
				_, err := studentClient.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = classInfoClient.AddClass(context.Background(), &classInfoProto.ClassAddRequest{StudentID: classInfoReq.StudentID, ClassName: classInfoReq.ClassName})
				require.NoError(t, err)
				_, err = classInfoClient.GetAllClassesByStudent(context.Background(), &classInfoProto.ClassesGetRequest{StudentID: classInfoReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /ClassInfoService/AddClass")
				assert.Contains(t, messagesSlice[2], "Received Kafka Message: RequestType: /ClassInfoService/GetAllClassesByStudent")
			},
		},
		{
			name: "Fail ClassInfo Kafka Create",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)

				_, err := classInfoClient.AddClass(context.Background(), &classInfoProto.ClassAddRequest{StudentID: classInfoReq.StudentID, ClassName: classInfoReq.ClassName})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, producer, consumer, messages := setUpAll(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := consumer.ReadMessages(ctx, topic); err != nil {
					err := producer.Close()
					require.NoError(t, err)
					err = consumer.ConsumerClose()
					require.NoError(t, err)
				}
			}()
			tc.setup(t, db, producer)
			wg.Wait()
			tc.assertFunc(t, messages)
		})
	}
}

func TestStudentKafkaDelete(t *testing.T) {
	var (
		topic      = "CRUD_events"
		brokers    = []string{"localhost:9092"}
		studentReq = serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		classInfoReq = serviceEntities.ClassInfo{
			StudentID: 1,
			ClassName: "MATH",
		}
		migrationPath = "./../repository/migrations"
	)
	setUpAll := func(t *testing.T) (*postgres.TDB, kafka_test.ProducerInterface, kafka_test.ConsumerInterface, *bytes.Buffer) {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to initialize Zap logger: %v\n", err)
		}
		logger.SetGlobal(zapLogger.With(zap.String("component", "main")))
		db := postgres.NewFromEnv()
		t.Cleanup(func() { db.DB.GetPool(context.Background()).Close() })

		producer, err := kafka_test.NewProducer(brokers, topic)
		require.NoError(t, err)

		var messages bytes.Buffer
		consumer, err := kafka_test.NewKafkaConsumer(brokers, &messages)
		require.NoError(t, err)

		db.SetUpDatabase(migrationPath)

		t.Cleanup(func() { db.TearDownDatabase(migrationPath) })
		return db, producer, consumer, &messages
	}
	testCases := []struct {
		name       string
		setup      func(*testing.T, *postgres.TDB, kafka_test.ProducerInterface)
		assertFunc func(*testing.T, *bytes.Buffer)
	}{
		{
			name: "Success Student Kafka Delete",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls
				_, err := client.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = client.DeleteStudent(context.Background(), &studentProto.StudentDeleteRequest{StudentID: studentReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /StudentService/DeleteStudent")
			},
		},
		{
			name: "Fail Student Kafka Delete",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls

				_, err := client.DeleteStudent(context.Background(), &studentProto.StudentDeleteRequest{StudentID: studentReq.StudentID})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
		{
			name: "Success ClassInfo Kafka Delete",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})

				studentClient := studentProto.NewStudentServiceClient(conn)
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)
				_, err := studentClient.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = classInfoClient.AddClass(context.Background(), &classInfoProto.ClassAddRequest{StudentID: classInfoReq.StudentID, ClassName: classInfoReq.ClassName})
				require.NoError(t, err)
				_, err = classInfoClient.DeleteClassByStudent(context.Background(), &classInfoProto.ClassDeleteRequest{StudentID: classInfoReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /ClassInfoService/AddClass")
				assert.Contains(t, messagesSlice[2], "Received Kafka Message: RequestType: /ClassInfoService/DeleteClassByStudent")
			},
		},
		{
			name: "Fail ClassInfo Kafka Delete",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)

				_, err := classInfoClient.DeleteClassByStudent(context.Background(), &classInfoProto.ClassDeleteRequest{StudentID: classInfoReq.StudentID})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, producer, consumer, messages := setUpAll(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := consumer.ReadMessages(ctx, topic); err != nil {
					err := producer.Close()
					require.NoError(t, err)
					err = consumer.ConsumerClose()
					require.NoError(t, err)
				}
			}()
			tc.setup(t, db, producer)
			wg.Wait()
			tc.assertFunc(t, messages)
		})
	}
}

func TestStudentKafkaUpdate(t *testing.T) {
	var (
		topic      = "CRUD_events"
		brokers    = []string{"localhost:9092", "localhost:9093"}
		studentReq = serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		classInfoReq = serviceEntities.ClassInfo{
			StudentID: 1,
			ClassName: "MATH",
		}
		migrationPath = "./../repository/migrations"
	)
	setUpAll := func(t *testing.T) (*postgres.TDB, kafka_test.ProducerInterface, kafka_test.ConsumerInterface, *bytes.Buffer) {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to initialize Zap logger: %v\n", err)
		}
		logger.SetGlobal(zapLogger.With(zap.String("component", "main")))
		db := postgres.NewFromEnv()
		t.Cleanup(func() { db.DB.GetPool(context.Background()).Close() })

		producer, err := kafka_test.NewProducer(brokers, topic)
		require.NoError(t, err)

		var messages bytes.Buffer
		consumer, err := kafka_test.NewKafkaConsumer(brokers, &messages)
		require.NoError(t, err)

		db.SetUpDatabase(migrationPath)

		t.Cleanup(func() { db.TearDownDatabase(migrationPath) })
		return db, producer, consumer, &messages
	}
	testCases := []struct {
		name       string
		setup      func(*testing.T, *postgres.TDB, kafka_test.ProducerInterface)
		assertFunc func(*testing.T, *bytes.Buffer)
	}{
		{
			name: "Success Student Kafka Update",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC calls
				_, err := client.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = client.UpdateStudent(context.Background(), &studentProto.StudentUpdateRequest{StudentID: studentReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /StudentService/UpdateStudent")
			},
		},
		{
			name: "Fail Student Kafka Update",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
				})
				client := studentProto.NewStudentServiceClient(conn)
				//gRPC call
				_, err := client.UpdateStudent(context.Background(), &studentProto.StudentUpdateRequest{StudentID: studentReq.StudentID})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
		{
			name: "Success ClassInfo Kafka Update",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					studentStorage := repository.NewStudentStorage(db.DB)
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					studentService := server.NewStudentServer(&studentStorage, producer)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					studentProto.RegisterStudentServiceServer(s, studentService)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})
				studentClient := studentProto.NewStudentServiceClient(conn)
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)
				_, err := studentClient.CreateStudent(context.Background(), &studentProto.StudentCreateRequest{
					StudentName: studentReq.StudentName,
					Grade:       studentReq.Grade,
				})
				require.NoError(t, err)
				_, err = classInfoClient.AddClass(context.Background(), &classInfoProto.ClassAddRequest{StudentID: classInfoReq.StudentID, ClassName: classInfoReq.ClassName})
				require.NoError(t, err)
				_, err = classInfoClient.UpdateClass(context.Background(), &classInfoProto.ClassUpdateRequest{StudentID: classInfoReq.StudentID})
				require.NoError(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				consumerMessages := messages.String()
				messagesSlice := strings.Split(consumerMessages, "\n")
				assert.Contains(t, messagesSlice[0], "Received Kafka Message: RequestType: /StudentService/CreateStudent")
				assert.Contains(t, messagesSlice[1], "Received Kafka Message: RequestType: /ClassInfoService/AddClass")
				assert.Contains(t, messagesSlice[2], "Received Kafka Message: RequestType: /ClassInfoService/UpdateClass")
			},
		},
		{
			name: "Fail ClassInfo Kafka Update",
			setup: func(t *testing.T, db *postgres.TDB, producer kafka_test.ProducerInterface) {
				conn := setupGrpcServerAndClient(t, func(s *grpc.Server) {
					classInfoStorage := repository.NewClassInfoStorage(db.DB)
					classInfoService := server.NewClassInfoServer(&classInfoStorage, producer)
					classInfoProto.RegisterClassInfoServiceServer(s, classInfoService)
				})
				classInfoClient := classInfoProto.NewClassInfoServiceClient(conn)

				_, err := classInfoClient.UpdateClass(context.Background(), &classInfoProto.ClassUpdateRequest{StudentID: classInfoReq.StudentID})
				require.Error(t, err)
			},
			assertFunc: func(t *testing.T, messages *bytes.Buffer) {
				require.Equal(t, 0, messages.Len())
				assert.Equal(t, "", messages.String())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, producer, consumer, messages := setUpAll(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := consumer.ReadMessages(ctx, topic); err != nil {
					err := producer.Close()
					require.NoError(t, err)
					err = consumer.ConsumerClose()
					require.NoError(t, err)
				}
			}()
			tc.setup(t, db, producer)
			wg.Wait()
			tc.assertFunc(t, messages)
		})
	}
}
