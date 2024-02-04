package model

import (
	"time"

	"github.com/lib/pq"
)
type Advertisement struct {
    ID        uint           `json:"id" gorm:"primary_key"`
    Title     string         `json:"title"`
    StartAt   time.Time      `json:"startAt"`
    EndAt     time.Time      `json:"endAt"`
    AgeStart  *int           `json:"ageStart,omitempty"`
    AgeEnd    *int           `json:"ageEnd,omitempty"`
    Gender    pq.StringArray `json:"gender,omitempty" gorm:"type:varchar(10)[]"`
    Country   pq.StringArray `json:"country,omitempty" gorm:"type:varchar(10)[]"`
    Platform  pq.StringArray `json:"platform,omitempty" gorm:"type:varchar(10)[]"`
    Active    bool           `json:"active" gorm:"default:true"`
}
