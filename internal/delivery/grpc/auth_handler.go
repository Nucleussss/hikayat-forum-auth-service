package grpc

import (
	"context"

	"log"

	"github.com/Nucleussss/hikayat-forum/auth/internal/service"
	"github.com/Nucleussss/hikayat-forum/auth/pkg/utils"

	authpb "github.com/Nucleussss/hikayat-proto/gen/go/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles the register request and returns a response.
func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	op := "authhandler.register"
	log.Printf("Received register request received for : %v\n", req.GetEmail())

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

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

	// register the user
	res, err := h.authService.Register(ctx, req)
	if err != nil {
		log.Printf("%s Error registering user: %v\n", op, err)
		return nil, status.Error(codes.Internal, "Error registering user")
	}

	// log the registration success
	response := &authpb.RegisterResponse{
		Message: res.Message,
	}
	log.Printf("%s Registration successful for : %v\n", op, req.GetEmail())
	return response, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	op := "authHandler.Login"
	log.Printf("%s Received login request for email: %v ", op, req.GetEmail())

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

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

	// call the authService to login and get the token
	tokenString, err := h.authService.Login(ctx, req)
	if err != nil {
		log.Printf("%s Login failed for email: %v\n ", op, req.GetEmail())
		return nil, status.Error(codes.Unauthenticated, "Login failed")
	}

	// create a login response message
	response := &authpb.LoginResponse{
		Message: "Login Successful",
		Token:   tokenString.Token,
	}
	log.Printf("%s Login successful for : %v\n", op, req.GetEmail())
	return response, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.User, error) {
	op := "auth.GetUser"
	log.Printf("%s Get user request received: %v\n", op, req)

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

	// get the user ID from the context
	if err := utils.EnsureUserAuthorized(ctx, req.GetId()); err != nil {
		log.Printf("%s user was not autorized. %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// validate email format
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	// get the user from the service layer
	user, err := h.authService.GetUser(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// log the successful request processing
	log.Printf("%s user request processed was successful for: %v", op, user.Email)
	return user, nil

}

func (h *AuthHandler) UpdateUserProfile(ctx context.Context, req *authpb.UpdateUserProfileRequest) (*authpb.UpdateUserProfileResponse, error) {
	op := "authHandler.UpdateUserProfile"
	log.Printf("recieve update user profile request from client: %s", req.GetId())

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

	// get the user ID from the context
	if err := utils.EnsureUserAuthorized(ctx, req.GetId()); err != nil {
		log.Printf("%s user was not autorized. %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// validate user input
	if req.GetName() == "" {
		log.Printf("%s name was empty", op)
		return nil, status.Error(codes.InvalidArgument, "name was empty")
	}

	// update the user profile
	res, err := h.authService.UpdateUserProfile(ctx, req)
	if err != nil {
		log.Printf("%s failed to update user profile", op)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// return the updated user as a protobuf message
	response := &authpb.UpdateUserProfileResponse{
		Message: "update user profile successful for: " + req.GetName(),
		User:    res.User,
	}

	return response, nil
}

func (h *AuthHandler) ChangeUserEmail(ctx context.Context, req *authpb.ChangeUserEmailRequest) (*authpb.ChangeUserEmailResponse, error) {
	op := "authHandler.ChangeUserEmail"
	log.Printf("recieve change user email request from client: %s", req.GetId())

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

	// get the user ID from the context
	if err := utils.EnsureUserAuthorized(ctx, req.GetId()); err != nil {
		log.Printf("%s user was not autorized. %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// validate user input
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// validate email format
	if !utils.IsValidEmail(req.GetEmail()) {
		log.Printf("%s invalid email provided: %s", op, req.GetEmail())
		return nil, status.Error(codes.InvalidArgument, "email was invalid")
	}

	// call the service to change user email
	err := h.authService.ChangeUserEmail(ctx, req)
	if err != nil {
		log.Printf("%s failed to change user email: %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert the response to protobuf format
	res := &authpb.ChangeUserEmailResponse{
		Message: "change user email successful for: " + req.GetId(),
	}

	return res, err
}

func (h *AuthHandler) ChangeUserPassword(ctx context.Context, req *authpb.ChangeUserPasswordRequest) (*authpb.ChangeUserPasswordResponse, error) {
	op := "authHandler.ChangeUserPassword"
	log.Printf("recieve change user password request from client: %s", req.GetId())

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

	// get the user ID from the context
	if err := utils.EnsureUserAuthorized(ctx, req.GetId()); err != nil {
		log.Printf("%s user was not autorized. %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

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

	// call the ChangeUserPassword method of authService
	err := h.authService.ChangeUserPassword(ctx, req)
	if err != nil {
		log.Printf("%s failed to change user password due to error: %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &authpb.ChangeUserPasswordResponse{
		Message: "User password changed successfully",
	}

	return res, nil
}

func (h *AuthHandler) DeleteUser(ctx context.Context, req *authpb.DeleteUserRequest) (*authpb.DeleteUserResponse, error) {
	op := "authHandler.DeleteUser"
	log.Printf("recieve Delete user request from client: %s", req.GetId())

	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service not initialized")
	}

	// get the user ID from the context
	if err := utils.EnsureUserAuthorized(ctx, req.GetId()); err != nil {
		log.Printf("%s user was not autorized. %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// call the DeleteUser method of authService
	err := h.authService.DeleteUser(ctx, req)
	if err != nil {
		log.Printf("%s failed to delete user due to error: %v", op, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &authpb.DeleteUserResponse{
		Message: "User deleted successfully by id " + req.GetId(),
	}

	return res, nil
}
