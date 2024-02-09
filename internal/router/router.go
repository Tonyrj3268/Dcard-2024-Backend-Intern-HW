package router

import (
	"advertisement-api/internal/controller"
	"advertisement-api/internal/repository"

	"advertisement-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	adRepo := repository.NewAdRepository(db)
	adController := controller.NewAdController(adRepo)
	
	var group = r.Group("/api/v1") 
	{
		group.GET("hc", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		group.GET("ad", adController.GetAd)
		group.POST("ad", adController.CreateAd)
	}
	setDocs()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
func setDocs(){
	docs.SwaggerInfo.Title = "Dcard Advertisements API"
	docs.SwaggerInfo.Description = "This is a simple server for Dcard advertisments HW."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
}