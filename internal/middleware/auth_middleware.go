package middleware

import (
	"context"
	"log"
	"os"
	"strings"

	contextKey "github.com/Nucleussss/hikayat-forum/auth/internal/context"
	"github.com/Nucleussss/hikayat-forum/auth/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a gRPC unary server interceptor that provides authentication for incoming requests.
// It checks if a method is publicly accessible, and if not, it extracts and validates the JWT token
// from the authorization header. If the token is valid, it extracts the user ID and adds it to the request context
// before proceeding with the original gRPC handler. If authentication fails at any step,
// it returns an appropriate unauthorized or invalid argument status error.
func AuthInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	op := "server.AuthInterceptor"
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		// Define a map of publicly accessible methods that do not require authentication.
		publicMethod := map[string]bool{
			// Publicly accessible methods that do not require authentication
			"/hikayat.forum.v1.AuthService/Register": true,
			"/hikayat.forum.v1.AuthService/Login":    true,
		}
		// If the current method is in the publicMethod map, proceed without authentication.
		if publicMethod[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract metadata from the incoming context. Metadata typically contains request headers.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Printf("%s: %v", op, err)
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// Get the "authorization" header from the metadata.
		authHeader := md.Get("authorization")
		// Check if the authorization header is present.
		if len(authHeader) == 0 {
			log.Printf("%s: %v", op, err)
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Extract the JWT token by removing the "Bearer " prefix.
		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		// Validate the JWT token using the provided secret key.
		mapClaims, err := utils.ValidateJWTToken(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			log.Printf("%s: %v", op, err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		// Extract the "user_id" from the token claims.
		userID, ok := (*mapClaims)["user_id"].(string)
		// Check if the user_id is present and of the correct type.
		if !ok {
			log.Printf("%s: %v", op, err)
			return nil, status.Errorf(codes.InvalidArgument, "missing user_id in token claims")
		}

		// Set the extracted user ID into the context for downstream handlers to access.
		ctx = context.WithValue(ctx, contextKey.UserIDContextKey, userID)
		// Proceed with the original handler with the updated context.
		return handler(ctx, req)
	}
}
