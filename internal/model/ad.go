package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)


type Advertisement struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	StartAt   time.Time `json:"startAt" gorm:"index:idx_ad_start_end_age, priority:2"`
	EndAt     time.Time `json:"endAt" gorm:"index:idx_ad_start_end_age, priority:1"`
	AgeStart  *int      `json:"ageStart,omitempty" gorm:"index"`
	AgeEnd    *int      `json:"ageEnd,omitempty" gorm:"index"`
	Gender    *string   `json:"gender,omitempty"`
	Country   *string    `json:"country,omitempty"`
	Platform  *string    `json:"platform,omitempty"`
	Active    bool      `json:"active" gorm:"default:true"`
}

func (m *Advertisement) BeforeUpdate(db *gorm.DB) (err error) {
	now := time.Now()
	if m.EndAt.Before(now) {
		m.Active = false
	}
	return nil
}

func StringArrayToString(array *[]string) *string {
    if array == nil {
        return nil
    }
	result := strings.Join(*array, ",")
	return &result
}