package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Request
type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type GetUserRequest struct {
	ID uuid.UUID
}

type UpdateUserProfileRequest struct {
	Name string
	ID   uuid.UUID
}

type ChangeUserEmailRequest struct {
	Email string
	ID    uuid.UUID
}

type ChangeUserPasswordRequest struct {
	CurrentPassword string
	NewPassword     string
	ID              uuid.UUID
}

type DeleteUserRequest struct {
	ID uuid.UUID
}

// Response
type RegisterResponse struct {
	Message string
}

type LoginResponse struct {
	Message string
	Token   string
}

type UpdateUserProfileResponse struct {
	ID        uuid.UUID
	Name      string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
