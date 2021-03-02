package internal

import (
	"database/sql"
	"errors"
	"local/gokit-test/internal/model"
)

type AuthService interface {
	Authenticate(username, password string) (bool, error)
	GetUser(id uint) (model.User, error)
	GetUsers() []model.User
	CreateUser(user CreateUserRequest) model.User
}

type authService struct {
	repo model.UserRepository
}

func NewAuthService(repo model.UserRepository) AuthService {
	return authService{repo: repo}
}

func (authService) Authenticate(username, password string) (bool, error) {
	if username == "darsh" && password == "pass" {
		return true, nil
	}

	return false, errors.New("invalid credentials")
}

var ErrNotFound = errors.New("not found")

func (a authService) GetUser(id uint) (model.User, error) {
	user, err := a.repo.FindOneById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (a authService) GetUsers() []model.User {
	return a.repo.FindMany()
}

var ErrInvalidArgument = errors.New("invalid argument")

func (a authService) CreateUser(user CreateUserRequest) model.User {
	id := a.repo.Insert(user.Name, user.Phone)
	u, _ := a.GetUser(uint(id))

	return u
}
