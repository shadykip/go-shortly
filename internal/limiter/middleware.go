package limiter

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(limiter *IPLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP (handle proxies)
		ip := c.ClientIP()

		// 10 requests/sec, burst 10
		r := rate.Every(100 * time.Millisecond) // 10/sec
		b := 10

		if !limiter.Allow(ip, r, b) {
			retryAfter := limiter.RetryAfter(ip, r, b)
			c.Header("Retry-After", strconv.FormatInt(int64(retryAfter.Seconds()), 10))
			c.AbortWithStatusJSON(429, gin.H{
				"error":               "Too Many Requests",
				"retry_after_seconds": int64(retryAfter.Seconds()),
			})
			return
		}

		c.Next()
	}
}
