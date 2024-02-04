package main

import (
	"advertisement-api/internal/model"
	"advertisement-api/internal/router"

	"os"

	"github.com/go-redis/redis/v8"
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
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
	
	if err := db.AutoMigrate(&model.Advertisement{}); err != nil {
        panic(err)
    }
    rdb := redis.NewClient(&redis.Options{
        Addr:     "redis:6379", // Redis地址
        Password: "", // 如果設置了Redis密碼
        DB:       0,  // 默認數據庫編號
    })
	app := router.SetupRouter(db, rdb)
	if err := app.Run(":8080");err != nil {
		panic(err)
	}

}