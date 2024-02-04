package repository

import (
	"advertisement-api/internal/dto"
	"advertisement-api/internal/model"
	"time"

	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AdRepository interface {
	CreateAdvertisement(*model.Advertisement) error
	GetActiveAdvertisements(time.Time, dto.AdGetRequest) ([]dto.AdGetResponse, error)
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

func (r adRepository)GetActiveAdvertisements(now time.Time, adReq dto.AdGetRequest) ([]dto.AdGetResponse, error) {
    var ads []dto.AdGetResponse

    query := r.db.Model(&model.Advertisement{}).Where("start_at <= ? AND end_at >= ?", now, now)
    if adReq.Age != nil {
        query = query.Where("? BETWEEN age_start AND age_end", adReq.Age)
    }

    if adReq.Gender != nil{
        query = query.Where("gender @> ?", pq.Array([]string{*adReq.Gender}))
    }

    if adReq.Country!= nil {
        query = query.Where("country @> ?", pq.Array([]string{*adReq.Country}))
    }

    if adReq.Platform != nil{
        query = query.Where("platform @> ?", pq.Array([]string{*adReq.Platform}))
    }

    err := query.Select("title, end_at").Order("end_at ASC").Offset(adReq.Offset).Limit(adReq.Limit).Find(&ads).Error
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