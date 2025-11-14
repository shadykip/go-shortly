package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shadykip/go-shortly/internal/models"
	"gorm.io/gorm"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// @Summary Shorten a URL
// @Description Create a short URL for a given long URL
// @Tags urls
// @Accept json
// @Produce json
// @Param url body ShortenRequest true "URL to shorten"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /shorten [post]
func ShortenHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid URL", "details": err.Error()})
			return
		}

		// Generate unique short code
		var url models.URL
		maxAttempts := 5
		for i := 0; i < maxAttempts; i++ {
			code := models.GenerateShortCode()
			url = models.URL{Original: req.URL, ShortCode: code}
			result := db.Create(&url)
			if result.Error == nil {
				// Success!
				domain := os.Getenv("APP_DOMAIN")
				if domain == "" {
					domain = "http://localhost:8080" // local fallback
				}
				c.JSON(201, gin.H{
					"short_url": domain + "/" + code,
				})
				return
			}
			// If duplicate key error, retry
		}

		c.JSON(500, gin.H{"error": "Failed to generate unique short code"})
	}
}
