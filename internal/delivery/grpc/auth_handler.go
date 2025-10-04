package grpc

import (
	"context"

	"log"

	"github.com/Nucleussss/hikayat-forum/auth/internal/models"
	"github.com/Nucleussss/hikayat-forum/auth/internal/service"
	"github.com/Nucleussss/hikayat-forum/auth/pkg/utils"
	"github.com/google/uuid"

	pb "github.com/Nucleussss/hikayat-forum/auth/api/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// contextKey is an unexported type for keys in context.
// This prevents collisions with other packages.
type contextKey string

const (
	UserIDContextKey contextKey = "user_id"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	op := "authhandler.register"

	log.Printf("Received register request received for : %v\n", req.GetEmail())

	// check input validation
	if req.GetName() == "" || req.GetEmail() == "" || req.GetPassword() == "" {
		log.Printf("%s Invalid input: name, email or password invalid\n", op)
		return nil, status.Error(codes.InvalidArgument, "Invalid input: name, email or password invalid")
	}

	// validate email and password format
	if !utils.IsValidEmail(req.GetEmail()) {
		log.Printf("%s Invalid email format\n", op)
		return nil, status.Error(codes.InvalidArgument, "Invalid input: email invalid")
	}

	// validate password format
	if !utils.IsValidPassword(req.GetPassword()) {
		log.Printf("%s Invalid password format\n", op)
		return nil, status.Error(codes.InvalidArgument, "Invalid input: password invalid")
	}

	// create a new register request message
	userRegister := &models.RegisterRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	// register the user
	res, err := h.authService.Register(ctx, userRegister)
	if err != nil {
		log.Printf("%s Error registering user: %v\n", op, err)
		return nil, status.Error(codes.Internal, "Error registering user")
	}

	// log the registration success
	response := &pb.RegisterResponse{
		Message: res.Message,
	}
	log.Printf("%s Registration successful for : %v\n", op, req.GetEmail())
	return response, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	op := "authHandler.Login"
	log.Printf("%s Received login request for email: %v ", op, req.GetEmail())

	// check input validation
	if req.GetEmail() == "" || req.GetPassword() == "" {
		log.Printf("%s Invalid input: email or password invalid\n", op)
		return nil, status.Error(codes.InvalidArgument, "Invalid input: email or password invalid")
	}

	// validate email and password format
	if !utils.IsValidEmail(req.GetEmail()) {
		log.Printf("%s Invalid email format\n", op)
		return nil, status.Error(codes.InvalidArgument, "Invalid input: email invalid")
	}

	// validate password format
	if !utils.IsValidPassword(req.GetPassword()) {
		log.Printf("%s Invalid password format\n", op)
		return nil, status.Error(codes.InvalidArgument, "Invalid input: password invalid")
	}

	// create a new login request message
	userLogin := &models.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	// call the authService to login and get the token
	tokenString, err := h.authService.Login(ctx, userLogin)
	if err != nil {
		log.Printf("%s Login failed for email: %v\n ", op, req.GetEmail())
		return nil, status.Error(codes.Unauthenticated, "Login failed")
	}

	// create a login response message
	response := &pb.LoginResponse{
		Message: "Login Successful",
		Token:   tokenString.Token,
	}
	log.Printf("%s Login successful for : %v\n", op, req.GetEmail())
	return response, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	op := "auth.GetUser"
	log.Printf("%s Get user request received: %v\n", op, req)

	// validate email format
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	userID := &models.GetUserRequest{
		ID: uuid.MustParse(req.GetId()),
	}

	// get the user from the service layer
	user, err := h.authService.GetUser(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// convert the user to a protobuf message
	response := &pb.User{
		Id:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	//
	log.Printf("%s user request processed was successful for: %v", op, user.Email)
	return response, nil

}

func (h *AuthHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	op := "authHandler.UpdateUserProfile"
	log.Printf("recieve update user profile request from client: %s", req.GetId())

	// validate user input
	if req.GetName() == "" {
		log.Printf("%s name was empty", op)
		return nil, status.Error(codes.InvalidArgument, "name was empty")
	}

	// convert the protobuf message to a models.User struct
	user := &models.UpdateUserProfileRequest{
		ID:   uuid.MustParse(req.GetId()),
		Name: req.GetName(),
	}

	// update the user profile
	res, err := h.authService.UpdateUserProfile(ctx, user)
	if err != nil {
		log.Printf("%s failed to update user profile", op)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// convert the user to a protobuf message
	userConvert := &pb.User{
		Id:        res.ID.String(),
		Name:      res.Name,
		Email:     res.Email,
		IsActive:  res.IsActive,
		CreatedAt: timestamppb.New(res.CreatedAt),
		UpdatedAt: timestamppb.New(res.UpdatedAt),
	}

	// return the updated user as a protobuf message
	response := &pb.UpdateUserProfileResponse{
		Message: "update user profile successful for: " + req.GetName(),
		User:    userConvert,
	}

	return response, nil
}

func (h *AuthHandler) ChangeUserEmail(ctx context.Context, req *pb.ChangeUserEmailRequest) (*pb.ChangeUserEmailResponse, error) {
	op := "authHandler.ChangeUserEmail"
	log.Printf("recieve change user email request from client: %s", req.GetId())

	// validate user input
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// validate email format
	if !utils.IsValidEmail(req.GetEmail()) {
		log.Printf("%s invalid email provided: %s", op, req.GetEmail())
		return nil, status.Error(codes.InvalidArgument, "email was invalid")
	}

	// create a new change user email request with the provided
	userEmail := &models.ChangeUserEmailRequest{
		ID:    uuid.MustParse(req.GetId()),
		Email: req.GetEmail(),
	}

	// call the service to change user email
	err := h.authService.ChangeUserEmail(ctx, userEmail)
	if err != nil {
		log.Printf("%s failed to change user email: %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert the response to protobuf format
	res := &pb.ChangeUserEmailResponse{
		Message: "change user email successful for: " + req.GetId(),
	}

	return res, err
}

func (h *AuthHandler) ChangeUserPassword(ctx context.Context, req *pb.ChangeUserPasswordRequest) (*pb.ChangeUserPasswordResponse, error) {
	op := "authHandler.ChangeUserPassword"
	log.Printf("recieve change user password request from client: %s", req.GetId())

	// validate user input
	if req.GetNewpassword() == "" || req.GetCurrentpassword() == "" {
		log.Printf("%s failed to change user password due to empty fields", op)
		return nil, status.Error(codes.InvalidArgument, "new password and current password cannot be empty")
	}

	// validate password length
	if !utils.IsValidPassword(req.GetNewpassword()) {
		log.Printf("%s failed to change user password due invalid password", op)
		return nil, status.Error(codes.InvalidArgument, "pasword minal have 8 character")
	}

	// create a new instance of ChangeUserPasswordRequest with	the provided fields
	userPassword := &models.ChangeUserPasswordRequest{
		ID:              uuid.MustParse(req.GetId()),
		CurrentPassword: req.GetCurrentpassword(),
		NewPassword:     req.GetNewpassword(),
	}

	// call the ChangeUserPassword method of authService
	err := h.authService.ChangeUserPassword(ctx, userPassword)
	if err != nil {
		log.Printf("%s failed to change user password due to error: %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pb.ChangeUserPasswordResponse{
		Message: "User password changed successfully",
	}

	return res, nil
}

func (h *AuthHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	op := "authHandler.DeleteUser"
	log.Printf("recieve Delete user request from client: %s", req.GetId())

	userId := &models.DeleteUserRequest{
		ID: uuid.MustParse(req.GetId()),
	}

	// call the DeleteUser method of authService
	err := h.authService.DeleteUser(ctx, userId)
	if err != nil {
		log.Printf("%s failed to delete user due to error: %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pb.DeleteUserResponse{
		Message: "User deleted successfully by id " + req.GetId(),
	}

	return res, nil
}
