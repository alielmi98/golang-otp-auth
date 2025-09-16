package usecase

import (
	"context"

	"github.com/alielmi98/golang-otp-auth/internal/user/api/dto"
	"github.com/alielmi98/golang-otp-auth/internal/user/domain/auth"
	model "github.com/alielmi98/golang-otp-auth/internal/user/domain/models"
	"github.com/alielmi98/golang-otp-auth/internal/user/domain/repository"
	"github.com/alielmi98/golang-otp-auth/internal/user/entity"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
)

type UserUsecase struct {
	cfg         *config.Config
	repo        repository.UserRepository
	token       auth.TokenProvider
	otpProvider auth.OtpProvider
}

func NewUserUsecase(cfg *config.Config, repository repository.UserRepository, token auth.TokenProvider, otpProvider auth.OtpProvider) *UserUsecase {
	return &UserUsecase{
		cfg:         cfg,
		repo:        repository,
		token:       token,
		otpProvider: otpProvider,
	}
}

// Register/login by mobile number
func (u *UserUsecase) RegisterAndLoginByMobileNumber(ctx context.Context, mobileNumber string, otp string) (*dto.TokenDetail, error) {
	err := u.otpProvider.ValidateOtp(mobileNumber, otp)
	if err != nil {
		return nil, err
	}
	exists, err := u.repo.ExistsMobileNumber(ctx, mobileNumber)
	if err != nil {
		return nil, err
	}

	user := model.User{MobileNumber: mobileNumber}

	if exists {
		user, err = u.repo.FetchUserInfo(ctx, user.MobileNumber)
		if err != nil {
			return nil, err
		}

		token, err := u.generateToken(&user)
		if err != nil {
			return nil, err
		}
		return token, nil
	}

	// Register and login
	user.RegisteredAt = user.CreatedAt
	user, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user, err = u.repo.FetchUserInfo(ctx, user.MobileNumber)
	token, err := u.generateToken(&user)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (s *UserUsecase) GetUserByMobileNumber(ctx context.Context, mobileNumber string) (dto.UserInfo, error) {
	user, err := s.repo.GetUserByMobileNumber(ctx, mobileNumber)
	if err != nil {
		return dto.UserInfo{}, err
	}

	// Map domain model to response DTO
	return dto.UserInfo{ID: user.Id, MobileNumber: user.MobileNumber, RegisteredAt: user.RegisteredAt}, nil
}

func (s *UserUsecase) RefreshToken(refreshToken string) (*dto.TokenDetail, error) {
	tokenDetail, err := s.token.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func (s *UserUsecase) generateToken(user *model.User) (*dto.TokenDetail, error) {
	tokenDto := entity.TokenPayload{UserId: user.Id, MobileNumber: user.MobileNumber}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := s.token.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *UserUsecase) GetAllUsers(ctx context.Context, page, pageSize int, mobileNumber string) (dto.UserList, error) {
	users, total, err := s.repo.GetAllUsers(ctx, page, pageSize, mobileNumber)
	if err != nil {
		return dto.UserList{}, err
	}
	userInfos := make([]dto.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = dto.UserInfo{
			ID:           user.Id,
			MobileNumber: user.MobileNumber,
			RegisteredAt: user.RegisteredAt,
		}
	}

	return dto.UserList{
		Users:    userInfos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
