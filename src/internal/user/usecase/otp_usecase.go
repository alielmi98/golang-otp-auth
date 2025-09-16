package usecase

import (
	"github.com/alielmi98/golang-otp-auth/internal/user/domain/auth"
	"github.com/alielmi98/golang-otp-auth/pkg/cache"
	"github.com/alielmi98/golang-otp-auth/pkg/common"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/go-redis/redis/v7"
)

type OtpUsecase struct {
	cfg         *config.Config
	redisClient *redis.Client
	otpProvider auth.OtpProvider
}

func NewOtpUsecase(cfg *config.Config, otpProvider auth.OtpProvider) *OtpUsecase {
	redis := cache.GetRedis()
	return &OtpUsecase{cfg: cfg, redisClient: redis, otpProvider: otpProvider}
}

func (u *OtpUsecase) SendOtp(mobileNumber string) error {
	otp := common.GenerateOtp()
	print(otp) // TODO: send otp to user by sms
	err := u.otpProvider.SetOtp(mobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}
