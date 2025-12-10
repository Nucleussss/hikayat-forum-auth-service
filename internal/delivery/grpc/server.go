package grpc

import (
	"os"

	"github.com/Nucleussss/hikayat-forum/auth/internal/middleware"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func NewServer() *grpc.Server {
	// Create gRPC server options slice (if needed)
	var opts []grpc.ServerOption

	interceptor := []grpc.UnaryServerInterceptor{
		// Add interceptors/middleware here

		// This middleware was not activate bacause hikayat-gateway was already handle it.
		middleware.AuthInterceptor(os.Getenv("JWT_SECRET_KEY")),
	}

	// Append interceptors to options slice
	opts = append(opts, grpc.ChainUnaryInterceptor(interceptor...))

	// Create the gRPC server instance with options
	grpcServer := grpc.NewServer(opts...)

	// Enable gRPC reflection (optional, useful during development/debugging with tools like grpcurl)
	// Remove this in production if not needed for introspection.
	reflection.Register(grpcServer)

	// Return the configured server instance
	return grpcServer

}
