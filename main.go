package main

import (
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shadykip/go-shortly/internal/cache"
	"github.com/shadykip/go-shortly/internal/handlers"
	"github.com/shadykip/go-shortly/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	r := gin.Default()
	r.POST("/shorten", handlers.ShortenHandler(db))
	r.GET("/:code", handlers.RedirectHandler(db, redisCache))

	r.Run(":8080")
}
