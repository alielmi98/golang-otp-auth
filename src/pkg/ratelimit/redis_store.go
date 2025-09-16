package ratelimit

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

// RedisRateLimiter implements RateLimiter using Redis as storage
type RedisRateLimiter struct {
	client *redis.Client
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(client *redis.Client) RateLimiter {
	return &RedisRateLimiter{
		client: client,
	}
}

// CheckLimit implements sliding window rate limiting using Redis
func (r *RedisRateLimiter) CheckLimit(key string, limit int, window time.Duration) (bool, error) {
	countKey := fmt.Sprintf("rate_limit_count:%s", key)
	startKey := fmt.Sprintf("rate_limit_start:%s", key)

	// Get current window start time
	startTs, err := r.client.Get(startKey).Int64()
	now := time.Now().Unix()
	
	// Handle Redis nil response (key doesn't exist)
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("failed to get start timestamp: %w", err)
	}

	// If no previous window or window has expired, start new window
	if err == redis.Nil || now-startTs >= int64(window.Seconds()) {
		pipe := r.client.TxPipeline()
		pipe.Set(startKey, now, window)
		pipe.Set(countKey, 1, window)
		_, err := pipe.Exec()
		if err != nil {
			return false, fmt.Errorf("failed to initialize new window: %w", err)
		}
		return true, nil
	}

	// Check current count in the window
	count, err := r.client.Get(countKey).Int()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("failed to get current count: %w", err)
	}

	// If count doesn't exist, initialize it
	if err == redis.Nil {
		count = 0
	}

	// Check if limit is exceeded
	if count >= limit {
		return false, nil
	}

	// Increment counter
	err = r.client.Incr(countKey).Err()
	if err != nil {
		return false, fmt.Errorf("failed to increment counter: %w", err)
	}

	return true, nil
}

// GetRemainingAttempts returns the number of remaining attempts
func (r *RedisRateLimiter) GetRemainingAttempts(key string, limit int, window time.Duration) (int, error) {
	countKey := fmt.Sprintf("rate_limit_count:%s", key)
	startKey := fmt.Sprintf("rate_limit_start:%s", key)

	// Get current window start time
	startTs, err := r.client.Get(startKey).Int64()
	now := time.Now().Unix()
	
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get start timestamp: %w", err)
	}

	// If no previous window or window has expired, return full limit
	if err == redis.Nil || now-startTs >= int64(window.Seconds()) {
		return limit, nil
	}

	// Get current count
	count, err := r.client.Get(countKey).Int()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get current count: %w", err)
	}

	if err == redis.Nil {
		count = 0
	}

	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// GetResetTime returns when the rate limit will reset
func (r *RedisRateLimiter) GetResetTime(key string, window time.Duration) (time.Time, error) {
	startKey := fmt.Sprintf("rate_limit_start:%s", key)

	// Get current window start time
	startTs, err := r.client.Get(startKey).Int64()
	if err != nil && err != redis.Nil {
		return time.Time{}, fmt.Errorf("failed to get start timestamp: %w", err)
	}

	// If no window exists, reset time is now
	if err == redis.Nil {
		return time.Now(), nil
	}

	resetTime := time.Unix(startTs, 0).Add(window)
	return resetTime, nil
}
