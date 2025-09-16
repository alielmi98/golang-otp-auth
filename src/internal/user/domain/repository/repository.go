package repository

import (
	"context"

	model "github.com/alielmi98/golang-otp-auth/internal/user/domain/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) (model.User, error)
	Update(ctx context.Context, id int, user *model.User) error
	Delete(ctx context.Context, id int) error
	GetUserByMobileNumber(ctx context.Context, mobileNumber string) (model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetDefaultRole(ctx context.Context) (roleId int, err error)
}
