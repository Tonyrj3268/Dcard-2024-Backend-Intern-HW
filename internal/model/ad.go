package model

import (
	"time"

	"github.com/lib/pq"
)
type Advertisement struct {
    ID        uint           `json:"id" gorm:"primary_key"`
    Title     string         `json:"title"`
    StartAt   time.Time      `json:"start_at"`
    EndAt     time.Time      `json:"end_at"`
    AgeStart  *int           `json:"age_start,omitempty"`
    AgeEnd    *int           `json:"age_end,omitempty"`
    Gender    pq.StringArray `json:"gender,omitempty" gorm:"type:varchar(10)[]"`
    Country   pq.StringArray `json:"country,omitempty" gorm:"type:varchar(10)[]"`
    Platform  pq.StringArray `json:"platform,omitempty" gorm:"type:varchar(10)[]"`
}
