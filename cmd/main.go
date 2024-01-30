package main

import (
	"advertisement-api/internal/model"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    dbHost := os.Getenv("POSTGRES_HOSTNAME")
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbDatabase := os.Getenv("POSTGRES_DB")

    dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbDatabase + " sslmode=disable TimeZone=Asia/Taipei"
    db, connect_err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if connect_err != nil {
        panic(connect_err)
    }
	migrateErr := db.AutoMigrate(&model.Advertisement{})
	if migrateErr != nil {
		panic(migrateErr)
	}
	app := gin.Default()
	app.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(200, gin.H{
			"message": "hello " + name,
		})
	})
	server_err := app.Run(":8080")
	if server_err != nil {
		panic(server_err)
	}

}