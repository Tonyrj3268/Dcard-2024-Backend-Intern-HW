package router

import (
	"advertisement-api/internal/controller"
	"advertisement-api/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	// Health check
	r.GET("/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	adRepo := repository.NewAdRepository(db)
	adController := controller.NewAdController(adRepo)
	var group = r.Group("/api/v1/ad") 
	{
		group.GET("/", adController.GetAd)
		group.POST("/", adController.CreateAd)
	}

	return r
}