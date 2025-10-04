package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nucleussss/hikayat-forum/auth/internal/models"
	"github.com/Nucleussss/hikayat-forum/auth/internal/repository"
	"github.com/google/uuid"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

// Register a new user in the database
func (r *userRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email, password_hash, is_active, created_at, updated_at  FROM users 
		WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Check if the row was found or not
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, err
}

// FindUserById function to find user by ID in database.
func (r *userRepo) FindUserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email, password_hash, is_active, created_at, updated_at 
		FROM users WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Check if the row was found or not
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, err
}

// CreateNewUser function creates a new user in the database and returns an error if it fails to create the user.
func (r *userRepo) CreateNewUser(ctx context.Context, user *models.RegisterRequest) error {
	query := `
		INSERT INTO users (name, email, password_hash) 
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	return err
}

// ExistByEmail function checks if a user with the given email exists in the database
func (r *userRepo) ExistByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE email = $1
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)

	return exists, err
}

// UpdateUserProfile function updates the user profile in the database
func (r *userRepo) UpdateUserProfile(ctx context.Context, user *models.UpdateUserProfileRequest) (*models.UpdateUserProfileResponse, error) {
	query := `
		UPDATE users 
		SET name = $1 
		WHERE id = $2
		RETURNING id, email, name, is_active, created_at, updated_at
	`
	var updatedUser models.UpdateUserProfileResponse
	err := r.db.QueryRowContext(ctx, query, user.Name, user.ID).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.Name,
		&updatedUser.IsActive,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &updatedUser, nil
}

// ChangeUserPassword changes the password of a user.
func (r *userRepo) ChangeUserPassword(ctx context.Context, user *models.ChangeUserPasswordRequest) error {
	query := `
		UPDATE users 
		SET password_hash = $1 
		WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, user.NewPassword, user.ID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("error user not found")
	}

	return nil
}

// ChangeUserEmail changes the email of a user in
func (r *userRepo) ChangeUserEmail(ctx context.Context, user *models.ChangeUserEmailRequest) error {
	query := `
		UPDATE users 
		SET email = $1 
		WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, user.Email, user.ID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("error user not found")
	}

	return nil
}

// DeleteUser deletes a user from the database.
func (r *userRepo) DeleteUser(ctx context.Context, user *models.DeleteUserRequest) error {
	query := `
		DELETE FROM users 
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("error user not found")
	}

	return nil
}

func (r *userRepo) GetUserPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	query := `
		SELECT password_hash FROM users 
		WHERE id = $1
	`

	var PasswordHash string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("error user not found")
		}
		return "", err
	}

	return PasswordHash, nil
}
