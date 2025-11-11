package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shadykip/go-shortly/internal/cache"
	"github.com/shadykip/go-shortly/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost user=dev password=linspace dbname=taskflow port=5432 sslmode=disable"), &gorm.Config{})
	assert.NoError(t, err)
	db.Exec("DELETE FROM urls")
	return db
}

func TestRedirectHandler_CacheMissThenHit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Exec("DELETE FROM urls")

	// Setup Redis test instance (or use real one on :6379)
	redisCache := cache.NewRedisCache("localhost:6379")
	redisCache.Get("test") // clear if needed

	r := gin.New()
	r.GET("/:code", RedirectHandler(db, redisCache))

	// Seed DB
	db.Create(&models.URL{Original: "https://test.com", ShortCode: "test123"})

	// First request: cache miss → DB hit
	req1 := httptest.NewRequest("GET", "/test123", nil)
	resp1 := httptest.NewRecorder()
	r.ServeHTTP(resp1, req1)
	assert.Equal(t, 301, resp1.Code)
	assert.Equal(t, "https://test.com", resp1.Header().Get("Location"))

	// Second request: cache hit
	req2 := httptest.NewRequest("GET", "/test123", nil)
	resp2 := httptest.NewRecorder()
	r.ServeHTTP(resp2, req2)
	assert.Equal(t, 301, resp2.Code)

	// Verify cached
	val, _ := redisCache.Get("test123")
	assert.Equal(t, "https://test.com", val)
}
