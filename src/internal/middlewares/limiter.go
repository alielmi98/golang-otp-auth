package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alielmi98/golang-otp-auth/pkg/cache"
	"github.com/gin-gonic/gin"
)

const otpLimit = 3
const otpWindow = 10 * time.Minute

func OtpRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.Param("mobile_number")
		if mobile == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "mobile_number is required"})
			return
		}

		countKey := fmt.Sprintf("otp_limit_count:%s", mobile)
		startKey := fmt.Sprintf("otp_limit_start:%s", mobile)

		redis := cache.GetRedis()

		startTs, err := redis.Get(startKey).Int64()
		now := time.Now().Unix()
		if err != nil && err.Error() != "redis: nil" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if startTs == 0 || now-startTs >= int64(otpWindow.Seconds()) {
			pipe := redis.TxPipeline()
			pipe.Set(startKey, now, otpWindow)
			pipe.Set(countKey, 1, otpWindow)
			_, _ = pipe.Exec()
			c.Next()
			return
		}

		count, err := redis.Get(countKey).Int()
		if err != nil && err.Error() != "redis: nil" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if count >= otpLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "OTP request limit reached. Try again later."})
			return
		}

		redis.Incr(countKey)
		c.Next()
	}
}
