package service

import (
	"context"

	authpb "github.com/Nucleussss/hikayat-proto/gen/go/auth/v1"
)

type AuthService interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error)
	GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.User, error)
	UpdateUserProfile(ctx context.Context, req *authpb.UpdateUserProfileRequest) (*authpb.UpdateUserProfileResponse, error)
	ChangeUserPassword(ctx context.Context, req *authpb.ChangeUserPasswordRequest) error
	ChangeUserEmail(ctx context.Context, req *authpb.ChangeUserEmailRequest) error
	DeleteUser(ctx context.Context, user *authpb.DeleteUserRequest) error
}
