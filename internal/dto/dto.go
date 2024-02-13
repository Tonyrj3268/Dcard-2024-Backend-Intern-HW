package dto

import (
	"fmt"
	"time"
)

type AdGetRequest struct {
    Offset   int    `form:"offset"`
    Limit    int    `form:"limit" binding:"gte=1,lte=100"`
    Age      *int    `form:"age" binding:"omitempty,gte=1,lte=100"`
    Gender   *string `form:"gender" binding:"omitempty,oneof=M F"`
    Country  *string `form:"country" binding:"omitempty,iso3166_1_alpha2"`
    Platform *string `form:"platform" binding:"omitempty,oneof=android ios web"`
}
func (a *AdGetRequest) GetParams() string {
    params := fmt.Sprintf("offset=%d&limit=%d", a.Offset, a.Limit)
    if a.Age != nil {
        params = fmt.Sprintf("%s&age=%d", params, *a.Age)
    }
    if a.Gender != nil{
        params = fmt.Sprintf("%s&gender=%s", params,*a.Gender)
    }
    if a.Country != nil{
        params = fmt.Sprintf("%s&country=%s", params,*a.Country)
    }
    if a.Platform != nil{
        params = fmt.Sprintf("%s&platform=%s", params,*a.Platform)
    }
    return params
}
type AdGetResponse struct {
    Title string    `json:"title"`
    EndAt time.Time `json:"endAt"`
}
type AdCondition struct {
	AgeStart      *int    `json:"ageStart" binding:"omitempty,gte=1,lte=100,ltefield=AgeEnd"`   
    AgeEnd        *int    `json:"ageEnd" binding:"omitempty,gte=1,lte=100"`   
	Gender   *[]string `json:"gender" binding:"omitempty,dive,oneof=M F"` 
	Country  *[]string `json:"country" binding:"omitempty,dive,iso3166_1_alpha2"`
	Platform *[]string `json:"platform" binding:"omitempty,dive,oneof=android ios web"`
}

type AdCreationRequest struct {
    Title      string       `json:"title" binding:"required"`
    StartAt    time.Time    `json:"startAt" binding:"required,ltefield=EndAt" time_format:"2024-12-10T03:00:00.000Z"`
    EndAt      time.Time    `json:"endAt" binding:"required" time_format:"2024-12-10T03:00:00.000Z"`
    Conditions AdCondition  `json:"conditions" binding:"omitempty"`
}