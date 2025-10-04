package repository

import (
	"context"

	"github.com/Nucleussss/hikayat-forum/auth/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserById(ctx context.Context, id string) (*models.User, error)
	CreateNewUser(ctx context.Context, user *models.RegisterRequest) error
	ExistByEmail(ctx context.Context, email string) (bool, error)
	UpdateUserProfile(ctx context.Context, user *models.UpdateUserProfileRequest) (*models.UpdateUserProfileResponse, error)
	ChangeUserPassword(ctx context.Context, user *models.ChangeUserPasswordRequest) error
	ChangeUserEmail(ctx context.Context, user *models.ChangeUserEmailRequest) error
	DeleteUser(ctx context.Context, user *models.DeleteUserRequest) error
	GetUserPasswordHash(ctx context.Context, id uuid.UUID) (string, error)
}
