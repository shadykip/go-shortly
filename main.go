//go:build !test
// +build !test

package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shadykip/go-shortly/internal/cache"
	"github.com/shadykip/go-shortly/internal/handlers"
	"github.com/shadykip/go-shortly/internal/limiter"
	"github.com/shadykip/go-shortly/internal/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title URL Shortener API
// @version 1.0
// @description A high-performance URL shortening service
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=dev password=linspace dbname=taskflow port=5432 sslmode=disable"
	}
	// DB setup (local fallback)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("DB connection failed")
	}
	db.AutoMigrate(&models.URL{})

	// 🧠 Redis setup (local fallback)
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	} else {
		// Railway Redis URL is like: redis://user:pass@host:port
		// Parse host:port
		u, err := url.Parse(redisURL)
		if err != nil {
			log.Fatal("Invalid REDIS_URL")
		}
		redisURL = u.Host // e.g., "host:port"
	}
	redisCache := cache.NewRedisCache(redisURL)

	rl := limiter.NewIPLimiter(rate.Every(100*time.Millisecond), 10)

	r := gin.Default()
	r.Use(limiter.RateLimitMiddleware(rl))
	r.POST("/shorten", handlers.ShortenHandler(db))
	r.GET("/:code", handlers.RedirectHandler(db, redisCache))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
