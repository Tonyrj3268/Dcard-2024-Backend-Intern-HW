package repository

import (
	"advertisement-api/internal/model"
	"time"

	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AdRepository interface {
	CreateAdvertisement(*model.Advertisement) error
	GetActiveAdvertisements(*int, *string, *string, *string, int, int) ([]struct{Title string; EndAt time.Time}, error)
}

type adRepository struct {
	db *gorm.DB
    redis *redis.Client
}
func NewAdRepository(db *gorm.DB, rdb *redis.Client) *adRepository {
    return &adRepository{db: db, redis: rdb}
}

func (r adRepository)CreateAdvertisement(ad *model.Advertisement) error {
    if err := r.db.Create(ad).Error; err != nil {
        
        return err
    }
    c := context.Background()
    count, err := r.redis.Incr(c, "active_ads_count").Result();
    if err != nil {
        return err
    }
    if count >= 1000 {
        return r.UpdateActiveCount(c)
    }
    return nil
}

func (r adRepository)GetActiveAdvertisements(age *int, gender, country, platform *string, offset, limit int) ([]struct{Title string; EndAt time.Time}, error) {
    var ads []struct{Title string; EndAt time.Time}
    now := time.Now()

    query := r.db.Model(&model.Advertisement{}).Where("start_at <= ? AND end_at >= ?", now, now)
    if age != nil {
        query = query.Where("? BETWEEN age_start AND age_end", *age)
    }

    if gender != nil{
        query = query.Where("gender @> ?", pq.Array([]string{*gender}))
    }

    if country!= nil {
        query = query.Where("country @> ?", pq.Array([]string{*country}))
    }

    if platform != nil{
        query = query.Where("platform @> ?", pq.Array([]string{*platform}))
    }

    err := query.Select("title, end_at").Order("end_at ASC").Offset(offset).Limit(limit).Find(&ads).Error
    return ads, err
}

func (r adRepository) UpdateActiveCount(c context.Context) error {
    now := time.Now()
    result := r.db.Model(&model.Advertisement{}).Where("active = ? AND end_at < ?", true, now).Update("active", false)
    if result.Error != nil {
        return result.Error
    }
    count := int64(0)
    err := r.db.Model(&model.Advertisement{}).Where("active = ? AND start_at <= ? AND end_at >= ?", true, now, now).Count(&count).Error
    if err != nil {
        return err
    }
    if count < 1000 {
        return r.redis.Set(c, "active_ads_count", count, 0).Err()
    }
    if err := r.db.Where("active = ? AND start_at <= ? AND end_at >= ?", true, now, now).Order("id ASC").Limit(int(count - 999)).Update("active", false).Error; err != nil {
        return err
    }
    return r.redis.Set(c, "active_ads_count", 999, 0).Err()
}