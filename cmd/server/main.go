package main

import (
	"GRPC_Server/internal/configs"
	"GRPC_Server/internal/kafka"
	"GRPC_Server/internal/logger"
	"GRPC_Server/internal/repository"
	"GRPC_Server/internal/server"
	"GRPC_Server/pkg/connection"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func main() {
	// Initialize Zap logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v\n", err)
	}
	logger.SetGlobal(zapLogger.With(zap.String("component", "main")))
	// Jagaer configurations
	cfg := config.Configuration{
		ServiceName: "messages-service",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		zapLogger.Fatal("cannot create tracer", zap.Error(err))
	}

	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)
	// Use Zap logger for the application logs
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	defer cancel()

	if err := godotenv.Load(".env"); err != nil {
		logger.Fatalf(ctx, "Could not set up environment variable: %s", err)
	}

	brokers := []string{}
	brokers = append(brokers, os.Getenv("BROKER"))
	topic := os.Getenv("TOPIC")
	grpcPort := os.Getenv("GRPCPORT")
	httpPort := os.Getenv("HTTPPORT")

	dbConfig, err := configs.FromEnv()
	if err != nil {
		logger.Fatalf(ctx, "Could not get environment variable: %v", err)
	}
	// Create Kafka Producer
	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		logger.Fatalf(ctx, "Failed to create Kafka Producer: %v\n", err)
	}
	//Creating Kafka Consumer
	consumer, err := kafka.NewKafkaConsumer(brokers, os.Stdout)
	if err != nil {
		producer.Close()
		logger.Fatalf(ctx, "Failed to create Kafka Consumer: %v\n", err)
	}
	// Start Kafka message Consumer
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := consumer.ReadMessages(ctx, topic); err != nil {
			// Last closing Producer in order to prevent data lost
			if err := producer.Close(); err != nil {
				logger.Infof(ctx, "Error closing Kafka Producer: %v\n", err)
			}
			// First closing Consumer
			if err := consumer.ConsumerClose(); err != nil {
				logger.Infof(ctx, "Error closing Kafka Consumer: %v\n", err)
			}
			logger.Fatalf(ctx, "Closed Consumer: %v\n", err)
		}
	}()

	database, err := connection.NewDB(ctx, dbConfig)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect Database %s", err)
	}

	defer database.GetPool(ctx).Close()

	err = connection.MigrationUp(dbConfig)
	if err != nil {
		log.Fatalf("Failed to %v", err)
	}

	studentStorage := repository.NewStudentStorage(database)
	classInfoStorage := repository.NewClassInfoStorage(database)

	go server.StartGatewayRouter(ctx, "localhost"+grpcPort, httpPort)

	logger.Infof(ctx, "Starting gRPC server on port %s", grpcPort)

	err = server.StartGRPCServer(&studentStorage, &classInfoStorage, producer, grpcPort)
	if err != nil {
		logger.Fatalf(ctx, "gRPC Server forced to shutdown: %v", err)
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		// Wait for the termination signal
		<-quit
		log.Printf("Graceful shutdown signal received")
		// Perform graceful shut down
		if err := connection.MigrationDownAndCloseSql(dbConfig); err != nil {
			log.Printf("Error during migration down and closing SQL: %v", err)
		}

		// Close the database pool
		database.GetPool(ctx).Close()

		os.Exit(0)
	}()

	wg.Wait()
}
