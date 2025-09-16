package ratelimit

import "time"

// RateLimiter defines the interface for rate limiting operations
type RateLimiter interface {
	// CheckLimit checks if the key has exceeded the rate limit
	// Returns true if allowed, false if rate limited
	CheckLimit(key string, limit int, window time.Duration) (bool, error)
	
	// GetRemainingAttempts returns the number of remaining attempts for the key
	GetRemainingAttempts(key string, limit int, window time.Duration) (int, error)
	
	// GetResetTime returns the time when the rate limit will reset for the key
	GetResetTime(key string, window time.Duration) (time.Time, error)
}

// OTPRateLimitConfig holds configuration for OTP rate limiting
type OTPRateLimitConfig struct {
	MaxAttempts int           // Maximum attempts allowed
	Window      time.Duration // Time window for rate limiting
}

// DefaultOTPConfig returns default OTP rate limiting configuration
func DefaultOTPConfig() OTPRateLimitConfig {
	return OTPRateLimitConfig{
		MaxAttempts: 3,
		Window:      10 * time.Minute,
	}
}
