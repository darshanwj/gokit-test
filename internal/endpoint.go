package internal

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"

	"local/gokit-test/models"
)

type errorer interface {
	error() error
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type HomeRequest struct{}

type HomeResponse struct {
	Message string `json:"msg"`
}

type GetUserRequest struct {
	Id uint
}

type GetUserResponse struct {
	User *models.User `json:"user"`
	Err  error        `json:"error,omitempty"`
}

func (r GetUserResponse) error() error {
	return r.Err
}

type GetUsersRequest struct{}

type GetUsersResponse struct {
	Users models.UserSlice `json:"users"`
	Err   error            `json:"error,omitempty"`
}

func (r GetUsersResponse) error() error {
	return r.Err
}

type CreateUserRequest struct {
	Name  string  `json:"name" validate:"required"`
	Phone *string `json:"phone"`
}

type ValidationErrorResponse struct {
	Errors string `json:"errors"`
}

func MakeAuthenticateEndpoint(svc AuthService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthRequest)
		success, err := svc.Authenticate(req.Username, req.Password)
		if err != nil {
			return AuthResponse{Success: false, Error: err.Error()}, nil
		}
		return AuthResponse{Success: success}, nil
	}
}

func MakeHomeEndpoint() endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return HomeResponse{Message: "gokit test service"}, nil
	}
}

func MakeGetUserEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := svc.GetUser(ctx, req.Id)
		return GetUserResponse{User: user, Err: err}, nil
	}
}

func MakeGetUsersEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := svc.GetUsers(ctx)
		return GetUsersResponse{Users: users, Err: err}, nil
	}
}

func MakeCreateUserEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)

		validate := validator.New()
		err := validate.Struct(req)
		if err != nil {
			return ValidationErrorResponse{Errors: err.Error()}, nil
		}

		user, err := svc.CreateUser(ctx, req)

		return GetUserResponse{User: user, Err: err}, nil
	}
}
