package grpc

import (
	"os"

	"github.com/Nucleussss/hikayat-forum/auth/internal/middleware"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func NewServer() *grpc.Server {
	// 1. Create gRPC server options slice (if needed)
	var opts []grpc.ServerOption

	interceptor := []grpc.UnaryServerInterceptor{
		// 2. Add interceptors/middleware here
		middleware.AuthInterceptor(os.Getenv("JWT_SECRET_KEY")),
		// middleware.RateLimitInterceptor,
	}

	// 3. Append interceptors to options slice
	opts = append(opts, grpc.ChainUnaryInterceptor(interceptor...))

	// 4. Create the gRPC server instance with options
	grpcServer := grpc.NewServer(opts...)

	// 5. Enable gRPC reflection (optional, useful during development/debugging with tools like grpcurl)
	// Remove this in production if not needed for introspection.
	reflection.Register(grpcServer)

	// 6. Return the configured server instance
	return grpcServer
}
