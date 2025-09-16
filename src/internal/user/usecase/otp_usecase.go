package usecase

import (
	"github.com/alielmi98/golang-otp-auth/internal/user/domain/auth"
	"github.com/alielmi98/golang-otp-auth/pkg/cache"
	"github.com/alielmi98/golang-otp-auth/pkg/common"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/alielmi98/golang-otp-auth/pkg/ratelimit"
	"github.com/go-redis/redis/v7"
)

type OtpUsecase struct {
	cfg              *config.Config
	redisClient      *redis.Client
	otpProvider      auth.OtpProvider
	rateLimitService *ratelimit.OTPRateLimitService
}

func NewOtpUsecase(cfg *config.Config, otpProvider auth.OtpProvider, rateLimitService *ratelimit.OTPRateLimitService) *OtpUsecase {
	redis := cache.GetRedis()
	return &OtpUsecase{
		cfg:              cfg,
		redisClient:      redis,
		otpProvider:      otpProvider,
		rateLimitService: rateLimitService,
	}
}

func (u *OtpUsecase) SendOtp(mobileNumber string) error {
	// Check rate limit before sending OTP
	err := u.rateLimitService.CheckOTPRateLimit(mobileNumber)
	if err != nil {
		return err
	}

	// Generate and send OTP
	otp := common.GenerateOtp()
	print(otp) // TODO: send otp to user by sms
	err = u.otpProvider.SetOtp(mobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

// GetOTPRateLimitInfo returns rate limit information for a mobile number
func (u *OtpUsecase) GetOTPRateLimitInfo(mobileNumber string) (*ratelimit.OTPRateLimitInfo, error) {
	return u.rateLimitService.GetRateLimitInfo(mobileNumber)
}
