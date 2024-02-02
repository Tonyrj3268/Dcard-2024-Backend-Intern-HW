package repository

import (
	"advertisement-api/internal/model"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AdRepository interface {
	CreateAdvertisement(*model.Advertisement) error
	GetActiveAdvertisements(*int, *string, *string, *string, int, int) ([]struct{Title string; EndAt time.Time}, error)
}

type adRepository struct {
	db *gorm.DB
}
func NewAdRepository(db *gorm.DB) *adRepository {
    return &adRepository{db: db}
}

func (r adRepository)CreateAdvertisement(ad *model.Advertisement) error {
    return r.db.Create(ad).Error
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