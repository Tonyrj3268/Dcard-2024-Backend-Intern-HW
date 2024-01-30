package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)
type JSONStringArray []string
func (j JSONStringArray) Value() (driver.Value, error) {
    if len(j) == 0 {
        return nil, nil
    }
    return json.Marshal(j)
}

func (j *JSONStringArray) Scan(src interface{}) error {
    if src == nil {
        *j = nil
        return nil
    }
    bytes, ok := src.([]byte)
    if !ok {
        return errors.New("invalid data type for JSONStringArray")
    }
    return json.Unmarshal(bytes, &j)
}
type Advertisement struct {
    ID        uint           `json:"id" gorm:"primary_key"`
    Title     string         `json:"title"`
    StartAt   time.Time      `json:"start_at"`
    EndAt     time.Time      `json:"end_at"`
    AgeStart  *int           `json:"age_start,omitempty"`
    AgeEnd    *int           `json:"age_end,omitempty"`
    Gender    JSONStringArray `json:"gender,omitempty" gorm:"type:json"`
    Country   JSONStringArray `json:"country,omitempty" gorm:"type:json"`
    Platform  JSONStringArray `json:"platform,omitempty" gorm:"type:json"`
}

func (ad *Advertisement) Create(db *gorm.DB) error {
    return db.Create(ad).Error
}