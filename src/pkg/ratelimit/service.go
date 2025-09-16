package ratelimit

import (
	"fmt"
	"time"

	"github.com/alielmi98/golang-otp-auth/pkg/service_errors"
)

// OTPRateLimitService provides OTP-specific rate limiting functionality
type OTPRateLimitService struct {
	rateLimiter RateLimiter
	config      OTPRateLimitConfig
}

// NewOTPRateLimitService creates a new OTP rate limiting service
func NewOTPRateLimitService(rateLimiter RateLimiter, config OTPRateLimitConfig) *OTPRateLimitService {
	return &OTPRateLimitService{
		rateLimiter: rateLimiter,
		config:      config,
	}
}

// CheckOTPRateLimit checks if OTP can be sent to the given mobile number
func (s *OTPRateLimitService) CheckOTPRateLimit(mobileNumber string) error {
	key := fmt.Sprintf("otp:%s", mobileNumber)
	
	allowed, err := s.rateLimiter.CheckLimit(key, s.config.MaxAttempts, s.config.Window)
	if err != nil {
		return &service_errors.ServiceError{
			EndUserMessage:   "Internal server error",
			TechnicalMessage: "Rate limit check failed",
			Err:              err,
		}
	}

	if !allowed {
		resetTime, _ := s.rateLimiter.GetResetTime(key, s.config.Window)
		return &service_errors.ServiceError{
			EndUserMessage:   fmt.Sprintf("OTP request limit exceeded. Try again after %s", resetTime.Format("15:04:05")),
			TechnicalMessage: "OTP rate limit exceeded",
			Err:              nil,
		}
	}

	return nil
}

// GetRemainingAttempts returns the number of remaining OTP attempts for a mobile number
func (s *OTPRateLimitService) GetRemainingAttempts(mobileNumber string) (int, error) {
	key := fmt.Sprintf("otp:%s", mobileNumber)
	
	remaining, err := s.rateLimiter.GetRemainingAttempts(key, s.config.MaxAttempts, s.config.Window)
	if err != nil {
		return 0, &service_errors.ServiceError{
			EndUserMessage:   "Internal server error",
			TechnicalMessage: "Failed to get remaining attempts",
			Err:              err,
		}
	}

	return remaining, nil
}

// GetResetTime returns when the rate limit will reset for a mobile number
func (s *OTPRateLimitService) GetResetTime(mobileNumber string) (time.Time, error) {
	key := fmt.Sprintf("otp:%s", mobileNumber)
	
	resetTime, err := s.rateLimiter.GetResetTime(key, s.config.Window)
	if err != nil {
		return time.Time{}, &service_errors.ServiceError{
			EndUserMessage:   "Internal server error",
			TechnicalMessage: "Failed to get reset time",
			Err:              err,
		}
	}

	return resetTime, nil
}

// GetRateLimitInfo returns comprehensive rate limit information for a mobile number
func (s *OTPRateLimitService) GetRateLimitInfo(mobileNumber string) (*OTPRateLimitInfo, error) {
	remaining, err := s.GetRemainingAttempts(mobileNumber)
	if err != nil {
		return nil, err
	}

	resetTime, err := s.GetResetTime(mobileNumber)
	if err != nil {
		return nil, err
	}

	return &OTPRateLimitInfo{
		MobileNumber:       mobileNumber,
		MaxAttempts:        s.config.MaxAttempts,
		RemainingAttempts:  remaining,
		WindowDuration:     s.config.Window,
		ResetTime:          resetTime,
		IsLimited:          remaining == 0,
	}, nil
}

// OTPRateLimitInfo contains comprehensive rate limit information
type OTPRateLimitInfo struct {
	MobileNumber       string        `json:"mobile_number"`
	MaxAttempts        int           `json:"max_attempts"`
	RemainingAttempts  int           `json:"remaining_attempts"`
	WindowDuration     time.Duration `json:"window_duration"`
	ResetTime          time.Time     `json:"reset_time"`
	IsLimited          bool          `json:"is_limited"`
}
