package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shadykip/go-shortly/internal/handlers"
	"github.com/shadykip/go-shortly/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// DB setup (local fallback)
	dsn := "host=localhost user=dev password=linspace dbname=taskflow port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("DB connection failed")
	}
	db.AutoMigrate(&models.URL{})

	r := gin.Default()
	r.POST("/shorten", handlers.ShortenHandler(db))

	r.Run(":8080")
}
