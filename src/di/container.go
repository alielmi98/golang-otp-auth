package di

import (
	contractAuth "github.com/alielmi98/golang-otp-auth/internal/user/domain/auth"
	contractAuthRepo "github.com/alielmi98/golang-otp-auth/internal/user/domain/repository"

	infraAuth "github.com/alielmi98/golang-otp-auth/internal/user/infra/auth"
	infraAuthRepo "github.com/alielmi98/golang-otp-auth/internal/user/infra/repository"

	"github.com/alielmi98/golang-otp-auth/pkg/config"
)

// midedlewares
func GetTokenProvider(cfg *config.Config) contractAuth.TokenProvider {
	return infraAuth.NewJwtProvider(cfg)
}

func GetUserRepository(cfg *config.Config) contractAuthRepo.UserRepository {
	return infraAuthRepo.NewUserPgRepo()
}

func GetOtpProvider(cfg *config.Config) contractAuth.OtpProvider {
	return infraAuth.NewOtpProvider(cfg)
}
