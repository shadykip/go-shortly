package limiter

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewIPLimiter(rate.Every(100*time.Millisecond), 2) // burst=2 for test
	r := gin.New()
	r.Use(RateLimitMiddleware(rl))
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	// First 2 requests: allowed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	}

	// 3rd request: denied
	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 429, resp.Code)
	assert.Contains(t, resp.Body.String(), "Too Many Requests")
}
