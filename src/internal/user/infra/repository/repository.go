package repository

import (
	"context"
	"log"

	model "github.com/alielmi98/golang-otp-auth/internal/user/domain/models"
	"github.com/alielmi98/golang-otp-auth/pkg/constants"
	"github.com/alielmi98/golang-otp-auth/pkg/db"
	"gorm.io/gorm"
)

const userFilterExp string = "mobile_number = ?"
const countFilterExp string = "count(*) > 0"

type PgRepo struct {
	db *gorm.DB
}

func NewUserPgRepo() *PgRepo {
	return &PgRepo{db: db.GetDb()}
}

func (r *PgRepo) CreateUser(ctx context.Context, u model.User) (model.User, error) {

	roleId, err := r.GetDefaultRole(ctx)
	if err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.DefaultRoleNotFound, err.Error())

		return u, err
	}
	tx := r.db.WithContext(ctx).Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())

		return u, err
	}
	err = tx.Create(&model.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return u, err
	}
	tx.Commit()
	return u, nil
}

func (r *PgRepo) Update(ctx context.Context, id int, user *model.User) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func (r *PgRepo) Delete(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func (r *PgRepo) GetUserByMobileNumber(ctx context.Context, mobileNumber string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Where(userFilterExp, mobileNumber).First(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *PgRepo) GetAllUsers(ctx context.Context, page, pageSize int, mobileNumber string) ([]model.User, int, error) {
	offset := (page - 1) * pageSize
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})
	if mobileNumber != "" {
		query = query.Where("mobile_number LIKE ?", "%"+mobileNumber+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (r *PgRepo) GetDefaultRole(ctx context.Context) (roleId int, err error) {

	if err = r.db.WithContext(ctx).Model(&model.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

func (r *PgRepo) FetchUserInfo(ctx context.Context, mobileNumber string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where(userFilterExp, mobileNumber).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *PgRepo) ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())

		return false, err
	}
	return exists, nil
}
