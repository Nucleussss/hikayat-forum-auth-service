package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nucleussss/hikayat-forum/auth/db"
	"github.com/Nucleussss/hikayat-forum/auth/internal/delivery/grpc"
	"github.com/Nucleussss/hikayat-forum/auth/internal/repository/postgres"
	"github.com/Nucleussss/hikayat-forum/auth/internal/service"

	pb "github.com/Nucleussss/hikayat-forum/auth/api/auth/v1"
)

func main() {

	log.Println("Starting Auth Service")

	// load configurations from .env file
	// cfg := config.LoadConfig()

	// initiate database connection
	connStr := db.ConnectionString()
	dbConn, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer func() {
		log.Println("Closing database connection...")
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// initiate user repository using the PostgreSQL database connection
	userRepo := postgres.NewUserRepository(dbConn)

	// initiate service layer
	authServie := service.NewAuthService(userRepo)

	// initiate auth handler
	authHandler := grpc.NewAuthHandler(authServie)

	grpcServer := grpc.NewServer()

	// register gRPC server with reflection for easy discovery and access
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// start gRPC server on the specified port
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen on port %s : %v", os.Getenv("GRPC_PORT"), err)
	}
	log.Printf("Starting gRPC server at %s\n", os.Getenv("GRPC_PORT"))

	// gracefull shutdown setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	grpcStopped := make(chan struct{})

	// start gRPC server in a separate goroutine
	go func() {
		defer close(grpcStopped)
		log.Printf("Starting GRPC Server..")
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("grpc server stopped with error: %v", err)
		} else {
			log.Printf("grpc server stopped gracefully")
		}
	}()

	// wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping gRPC server...")

	// graceful shutdown setup
	grpcServer.GracefulStop()
	log.Printf("GRPC Server stopped gracefully\n")

	// wait for server goroutine to finish
	<-grpcStopped
	log.Println("Auth Service Exited")

}
