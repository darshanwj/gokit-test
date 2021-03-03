package internal

import (
	"context"
	"database/sql"
	"errors"
	"local/gokit-test/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AuthService interface {
	Authenticate(username, password string) (bool, error)
	GetUser(ctx context.Context, id uint) (*models.User, error)
	GetUsers(ctx context.Context) (models.UserSlice, error)
	CreateUser(ctx context.Context, user CreateUserRequest) (*models.User, error)
}

type authService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) AuthService {
	return authService{db: db}
}

func (authService) Authenticate(username, password string) (bool, error) {
	if username == "darsh" && password == "pass" {
		return true, nil
	}

	return false, errors.New("invalid credentials")
}

var ErrNotFound = errors.New("not found")

func (a authService) GetUser(ctx context.Context, id uint) (*models.User, error) {
	user, err := models.FindUser(ctx, a.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNotFound
		}
		return user, err
	}

	// Fetch comments
	_ = user.L.LoadComments(ctx, a.db, true, user, nil)

	return user, nil
}

func (a authService) GetUsers(ctx context.Context) (models.UserSlice, error) {
	u, err := models.Users(qm.Load(models.UserRels.Comments)).All(ctx, a.db)

	return u, err
}

var ErrInvalidArgument = errors.New("invalid argument")

func (a authService) CreateUser(ctx context.Context, cur CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Name:  cur.Name,
		Phone: null.StringFromPtr(cur.Phone),
	}

	err := user.Insert(ctx, a.db, boil.Infer())

	return user, err
}
