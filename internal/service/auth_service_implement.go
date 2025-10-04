package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Nucleussss/hikayat-forum/auth/internal/models"
	"github.com/Nucleussss/hikayat-forum/auth/internal/repository"
	"github.com/Nucleussss/hikayat-forum/auth/pkg/utils"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Register handles new user registration. It first checks if the provided email already exists in the database.
// If not, it hashes the user's password for security and then proceeds to create a new user entry in the database.
// Upon successful creation, it returns a confirmation message.
func (s *authService) Register(ctx context.Context, newUser *models.RegisterRequest) (*models.RegisterResponse, error) {
	op := "authService.Register"

	// check if email was exist
	exists, err := s.userRepo.ExistByEmail(ctx, newUser.Email)
	if err != nil {
		log.Printf("%s Error checking user existence: %v", op, err)
		return nil, err
	}

	if exists {
		log.Printf("%s Email already exists", op)
		return nil, fmt.Errorf("%v Email already exists", op)
	}

	// hash password
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("%s Error hashing password: %v ", op, err)
		return nil, err
	}

	createNewUser := &models.RegisterRequest{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: hashedPassword,
	}

	// create user in the database
	err = s.userRepo.CreateNewUser(ctx, createNewUser)
	if err != nil {
		log.Printf("%s Error creating new user: % v", op, err)
		return nil, err
	}

	response := &models.RegisterResponse{
		Message: "User created successfully",
	}

	log.Println("User created successfully")

	return response, nil
}

// Login authenticates a user by verifying their email and password against the database.
// Upon successful verification, it generates a JSON Web Token (JWT) for the user, allowing them to access protected resources,
// and returns this token along with a success message.
func (s *authService) Login(ctx context.Context, userLogin *models.LoginRequest) (*models.LoginResponse, error) {
	op := "authService.Login"

	// find user by email
	user, err := s.userRepo.FindUserByEmail(ctx, userLogin.Email)
	if err != nil {
		log.Printf("%s Error finding user by email: %v", op, err)
		return nil, err
	}

	// verify password
	if !utils.VerifyPassword(user.PasswordHash, userLogin.Password) {
		log.Printf(" %s Error verifying password", op)
		return nil, fmt.Errorf("%s Invalid credentials", op)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	// generate JWT token
	generatedToken, err := utils.GenerateJWTToken(user.ID, jwtSecret)
	if err != nil {
		log.Printf("%s Error generating JWT token: % v", op, err)
		return nil, err
	}

	log.Printf(" %s Login successful for user %v", op, user.Name)

	response := &models.LoginResponse{
		Message: "Login successful",
		Token:   generatedToken,
	}

	return response, nil
}

// GetUser retrieves a user's profile information from the database based on their unique user ID.
// It queries the repository for the user and returns the user object if found,
// or an error if the user does not exist or a database issue occurs.
func (s *authService) GetUser(ctx context.Context, usr *models.GetUserRequest) (*models.User, error) {
	op := "authService.GetUser"

	user, err := s.userRepo.FindUserById(ctx, usr.ID.String())
	if err != nil {
		log.Printf("%s Error finding user by id: %s, error: %v", op, usr.ID, err)
		return nil, err
	}

	return user, nil
}

// UpdateUserProfile updates a user's profile details in the database.
// This service takes the updated user information and uses the repository to persist these changes,
// returning the updated user profile or an error if the update operation fails.
func (s *authService) UpdateUserProfile(ctx context.Context, usr *models.UpdateUserProfileRequest) (*models.UpdateUserProfileResponse, error) {
	op := "authService.UpdateUser"

	user, err := s.userRepo.UpdateUserProfile(ctx, usr)
	if err != nil {
		log.Printf("%s Error update profile for user by id: %s, error: %v", op, usr.ID, err)
		return nil, err
	}

	return user, nil
}

// ChangeUserPassword allows a user to update their account password.
// It first retrieves the user's current password hash from the database and verifies it against the provided current password.
// If correct, the new password is hashed and then updated in the database, ensuring secure password management.
func (s *authService) ChangeUserPassword(ctx context.Context, usr *models.ChangeUserPasswordRequest) error {
	op := "authService.ChangeUserPassword"

	// get the current password hash from database
	CurrHashPass, err := s.userRepo.GetUserPasswordHash(ctx, usr.ID)
	if err != nil {
		log.Printf(" %s Error getting user password hash for user by id: %s, error: %v", op, usr.ID, err)
		return err
	}

	// check if current password is correct
	if !utils.VerifyPassword(CurrHashPass, usr.CurrentPassword) {
		log.Printf("%s Error current password is incorrect for user by id: %s", op, usr.ID)
		return fmt.Errorf("current password is incorrect")
	}

	// hash password
	newHashedPassword, err := utils.HashPassword(usr.NewPassword)
	if err != nil {
		log.Printf("%s Error hashing password: %v ", op, err)
		return err
	}

	usrNewPassword := &models.ChangeUserPasswordRequest{
		ID:              usr.ID,
		CurrentPassword: usr.CurrentPassword,
		NewPassword:     newHashedPassword,
	}

	err = s.userRepo.ChangeUserPassword(ctx, usrNewPassword)
	if err != nil {
		log.Printf("%s Error change password for user by id: %s, error: %v", op, usr.ID, err)
		return err
	}

	return nil
}

// ChangeUserEmail allows a user to update their registered email address.
// This service first checks if the new email is already in use by another account.
// If the email is unique, it proceeds to update the user's email in the database,
// ensuring data integrity and preventing duplicate email addresses.
func (s *authService) ChangeUserEmail(ctx context.Context, usr *models.ChangeUserEmailRequest) error {
	op := "authService.ChangeUserEmail"

	// check if email was already used
	exist, err := s.userRepo.ExistByEmail(ctx, usr.Email)
	if err != nil {
		log.Printf("%s error check if email was exist: %s", op, usr.Email)
		return fmt.Errorf("error check if email was exist: %s, error : %v", usr.Email, err)
	}

	if exist {
		log.Printf("%s email was already exist: %s", op, usr.Email)
		return fmt.Errorf("email was already exist: %s, error : %v", usr.Email, err)
	}

	err = s.userRepo.ChangeUserEmail(ctx, usr)
	if err != nil {
		log.Printf("%s Error change email for user by id: %s, error: %v", op, usr.ID, err)
		return err
	}

	return nil
}

// DeleteUser removes a user account from the database.
// This service takes a request containing the user ID and instructs the repository to delete the corresponding user record.
// It handles any errors that may occur during the deletion process, such as if the user does not exist or a database issue arises.
func (s *authService) DeleteUser(ctx context.Context, user *models.DeleteUserRequest) error {
	op := "authService.DeleteUser"

	err := s.userRepo.DeleteUser(ctx, user)
	if err != nil {
		log.Printf("%s Error delete user by id: %s, error: %v", op, user.ID, err)
		return err
	}

	return nil

}
