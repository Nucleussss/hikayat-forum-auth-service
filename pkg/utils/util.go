package utils

import (
	"context"
	"regexp"

	contextKey "github.com/Nucleussss/hikayat-forum/auth/internal/context"
	"github.com/Nucleussss/hikayat-forum/auth/internal/models"
	authpb "github.com/Nucleussss/hikayat-proto/gen/go/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func EnsureUserAuthorized(ctx context.Context, reqID string) error {
	// get the user ID from the context
	userIDFromToken, ok := ctx.Value(contextKey.UserIDContextKey).(string)

	// check if the user ID is present in the context and is valid
	if !ok {
		return status.Error(codes.Internal, "user ID not found in context")
	}

	// check if the user ID in the token matches the user ID in the request	if not, return an error
	if reqID != userIDFromToken {
		return status.Error(codes.PermissionDenied, "cannot update another user's profile")
	}

	return nil
}

func AuthModelToPB(p *models.User) *authpb.User {
	if p == nil {
		return nil
	}

	return &authpb.User{
		Id:        p.ID.String(),
		Name:      p.Name,
		Email:     p.Email,
		IsActive:  p.IsActive,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}
