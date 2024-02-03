package router

import (
	"advertisement-api/internal/controller"
	"advertisement-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)


func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	// Health check
	r.GET("app/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	adRepo := repository.NewAdRepository(db,rdb)
	adController := controller.NewAdController(adRepo)
	var group = r.Group("app/api/v1/ad") 
	{
		group.GET("", adController.GetAd)
		group.POST("", adController.CreateAd)
	}

	return r
}