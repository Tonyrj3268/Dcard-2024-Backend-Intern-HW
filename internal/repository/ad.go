package repository

import (
	"advertisement-api/internal/dto"
	"advertisement-api/internal/model"
	"time"

	"gorm.io/gorm"
)

type AdRepository interface {
	CreateAdvertisement(*model.Advertisement) error
	GetActiveAdvertisements(time.Time, dto.AdGetRequest) ([]dto.AdGetResponse, error)
}

type adRepository struct {
	db *gorm.DB
}
func NewAdRepository(db *gorm.DB) *adRepository {
    return &adRepository{db: db}
}

func (r adRepository)CreateAdvertisement(ad *model.Advertisement) error {
    if err := r.db.Create(ad).Error; err != nil {
        return err
    }
    
    return nil
}
func UpdateActiveCount(db *gorm.DB) error {
    now := time.Now()
    count := int64(0)
    err := db.Model(&model.Advertisement{}).Where("active = ? AND start_at <= ? AND end_at >= ?", true, now, now).Count(&count).Error
    if err != nil {
        return err
    }
    return nil
}
func (r adRepository)GetActiveAdvertisements(now time.Time, adReq dto.AdGetRequest) ([]dto.AdGetResponse, error) {
    var ads []dto.AdGetResponse
    query := r.db.Model(&model.Advertisement{})
    query = query.Where("? BETWEEN start_at AND end_at", now)
    if adReq.Age != nil {
        query = query.Where("? BETWEEN age_start AND age_end", adReq.Age)
    }

    if adReq.Gender != nil{
        query = query.Where("FIND_IN_SET(?,gender)>0", *adReq.Gender)
    }

    if adReq.Country!= nil {
        query = query.Where("FIND_IN_SET(?,country)>0", *adReq.Country)
    }

    if adReq.Platform != nil{
        query = query.Where("FIND_IN_SET(?,platform)>0", *adReq.Platform)
    }
    
    err := query.Select("title, end_at").Order("end_at ASC").Offset(adReq.Offset).Limit(adReq.Limit).Find(&ads).Error
    return ads, err
}
