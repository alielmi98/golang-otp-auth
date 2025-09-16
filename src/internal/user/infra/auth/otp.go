package auth

import (
	"fmt"
	"time"

	"github.com/alielmi98/golang-otp-auth/pkg/cache"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/alielmi98/golang-otp-auth/pkg/constants"
	"github.com/alielmi98/golang-otp-auth/pkg/service_errors"
	"github.com/go-redis/redis/v7"
)

type OtpProvider struct {
	cfg         *config.Config
	redisClient *redis.Client
}
type otpDto struct {
	Value string
	Used  bool
}

func NewOtpProvider(cfg *config.Config) *OtpProvider {
	redis := cache.GetRedis()
	return &OtpProvider{
		cfg:         cfg,
		redisClient: redis,
	}
}

func (s *OtpProvider) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	val := &otpDto{
		Value: otp,
		Used:  false,
	}

	res, err := cache.Get[otpDto](s.redisClient, key)
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OptExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	err = cache.Set(s.redisClient, key, val, s.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (s *OtpProvider) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[otpDto](s.redisClient, key)
	if err != nil {
		return err
	} else if res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpNotValid}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(s.redisClient, key, res, s.cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
