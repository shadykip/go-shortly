package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shadykip/go-shortly/internal/cache"
	"github.com/shadykip/go-shortly/internal/models"
	"gorm.io/gorm"
)

func RedirectHandler(db *gorm.DB, redisCache *cache.RedisCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		if code == "" {
			c.Status(400)
			return
		}

		// 🔥 1. Check Redis cache first
		if originalURL, err := redisCache.Get(code); err == nil && originalURL != "" {
			c.Redirect(301, originalURL)
			return
		}

		// 🗃️ 2. If cache miss, hit database
		var url models.URL
		if err := db.Where("short_code = ?", code).First(&url).Error; err != nil {
			c.Status(404)
			return
		}

		// 📈 3. Increment click count (optional but useful)
		db.Model(&url).Update("clicks", url.Clicks+1)

		// 💾 4. Cache the result for next time
		redisCache.Set(code, url.Original)

		// 🚀 5. Redirect
		c.Redirect(301, url.Original)
	}
}
