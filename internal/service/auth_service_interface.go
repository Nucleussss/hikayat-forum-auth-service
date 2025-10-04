package service

import (
	"context"

	"github.com/Nucleussss/hikayat-forum/auth/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, userRegister *models.RegisterRequest) (*models.RegisterResponse, error)
	Login(ctx context.Context, userLogin *models.LoginRequest) (*models.LoginResponse, error)
	GetUser(ctx context.Context, usrID *models.GetUserRequest) (*models.User, error)
	UpdateUserProfile(ctx context.Context, user *models.UpdateUserProfileRequest) (*models.UpdateUserProfileResponse, error)
	ChangeUserPassword(ctx context.Context, user *models.ChangeUserPasswordRequest) error
	ChangeUserEmail(ctx context.Context, user *models.ChangeUserEmailRequest) error
	DeleteUser(ctx context.Context, user *models.DeleteUserRequest) error
}
