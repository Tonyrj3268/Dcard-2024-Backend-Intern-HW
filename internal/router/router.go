package router

import (
	"advertisement-api/internal/controller"

	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()
	// Health check
	r.GET("/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	var group = r.Group("/api/v1/ad") 
	{
		group.GET("/", controller.GetAd)
		group.POST("/", controller.CreateAd)
	}

	return r
}