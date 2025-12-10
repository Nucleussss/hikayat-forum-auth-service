package repository

import (
	"context"

	authpb "github.com/Nucleussss/hikayat-proto/gen/go/auth/v1"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*authpb.User, error)
	FindUserById(ctx context.Context, id string) (*authpb.User, error)
	CreateNewUser(ctx context.Context, req *authpb.RegisterRequest) error
	ExistByEmail(ctx context.Context, email string) (bool, error)
	UpdateUserProfile(ctx context.Context, req *authpb.UpdateUserProfileRequest) (*authpb.UpdateUserProfileResponse, error)
	ChangeUserPassword(ctx context.Context, req *authpb.ChangeUserPasswordRequest) error
	ChangeUserEmail(ctx context.Context, req *authpb.ChangeUserEmailRequest) error
	DeleteUser(ctx context.Context, user *authpb.DeleteUserRequest) error
	GetUserPasswordHash(ctx context.Context, identifier interface{}) (string, error)
}
