package utils

import (
	"context"
	"fmt"
	"time"

	"advertisement-api/internal/repository"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)
func nextTick(targetHour int) time.Duration {
	now := time.Now()
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), targetHour, 0, 0, 0, next.Location())
	diff := next.Sub(now)
	if diff < 0 {
		diff += 24 * time.Hour
	}
	return diff
}
func StartBackgroundTask(rds *redis.Client,db *gorm.DB,ctx context.Context) {
    for {
        select {
		// 1s 更新一次
		case <-time.After(1 * time.Second):
			UpdateActiveCount(db)
		// 每天 00:00 重置一次
		case <-time.After(nextTick(0)):
			doBackgroundTask(rds)
		case <-ctx.Done():
			fmt.Println("定時器關閉...")
			return
		}
    }
}

func doBackgroundTask(rds *redis.Client) {
	rds.Set(context.Background(), "CreatedAd", 0, 86400)
}

func UpdateActiveCount(db *gorm.DB) {
	repository.UpdateActiveCount(db)
}