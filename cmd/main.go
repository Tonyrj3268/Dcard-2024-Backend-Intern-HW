package main

import (
	"advertisement-api/internal/model"
	"advertisement-api/internal/router"
	"advertisement-api/internal/utils"
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
func main() {

	if err := godotenv.Load(); err != nil {
	    panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{PrepareStmt: true})
	if err != nil {
	    panic(err)
	}
	
	if err := db.AutoMigrate(&model.Advertisement{}); err != nil {
        panic(err)
    }
	DB, _ := db.DB()
	DB.SetMaxIdleConns(100)
	DB.SetMaxOpenConns(1000)
	
	redis := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
		PoolSize: 128,
    })
	if _, err := redis.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	app := router.SetupRouter(db, redis)
    
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

    go utils.StartBackgroundTask(redis, db, ctx)
    
	go func() {
		if err := app.Run(":8080"); err != nil {
			panic (err)
		}
	}()

	<-ctx.Done()
    fmt.Println("server closing...")

    time.Sleep(1 * time.Second)

}