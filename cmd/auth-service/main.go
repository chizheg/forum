package main

import (
	"fmt"
	"github.com/artemxgod/forum/internal/auth/delivery/grpc"
	"github.com/artemxgod/forum/internal/auth/repository/postgres"
	"github.com/artemxgod/forum/internal/auth/service"
	"github.com/artemxgod/forum/pkg/database"
	"github.com/artemxgod/forum/pkg/logger"
	pb "github.com/artemxgod/forum/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort = "50051"
)

func main() {
	// Initialize logger
	log, err := logger.New(true)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Initialize database connection
	db, err := database.NewPostgresDB(database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "forum_auth",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal("Failed to initialize database", err)
	}
	defer db.Close()

	// Initialize repository
	repo := postgres.NewRepository(db)

	// Initialize service
	svc := service.NewService(repo)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":"+defaultPort)
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, grpc.NewAuthServer(svc, log.Logger))

	// Start server
	go func() {
		log.Info("Starting gRPC server on port " + defaultPort)
		if err := s.Serve(lis); err != nil {
			log.Fatal("Failed to serve", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gRPC server...")
	s.GracefulStop()
} 