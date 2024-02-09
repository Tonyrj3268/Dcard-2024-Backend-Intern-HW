package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)
type Advertisement struct {
    ID        uint           `json:"id" gorm:"primary_key"`
    Title     string         `json:"title"`
    StartAt   time.Time      `json:"startAt" gorm:"index:idx_ad_start_end_age, priority:2"`
    EndAt     time.Time      `json:"endAt" gorm:"index:idx_ad_start_end_age, priority:1"`
    AgeStart  *int           `json:"ageStart,omitempty" gorm:"index"`
    AgeEnd    *int           `json:"ageEnd,omitempty" gorm:"index"`
    Gender    pq.StringArray `json:"gender,omitempty" gorm:"type:varchar(10)[]"`
    Country   pq.StringArray `json:"country,omitempty" gorm:"type:varchar(10)[]"`
    Platform  pq.StringArray `json:"platform,omitempty" gorm:"type:varchar(10)[]"`
    Active    bool           `json:"active" gorm:"default:true"`
}

func (m *Advertisement)BeforeUpdate(db *gorm.DB) (err error) {
    now := time.Now()
    if m.EndAt.Before(now) {
        m.Active = false
    }
    return nil
}