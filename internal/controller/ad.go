package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdGetRequest struct {
    Offset   *int    `form:"offset" binding:"omitempty"`
    Limit    *int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
    Age      *int    `form:"age" binding:"omitempty,gte=1,lte=100"`
    Gender   *string `form:"gender" binding:"omitempty,oneof=M F"`
    Country  *string `form:"country" binding:"omitempty,iso3166_1_alpha2"`
    Platform *string `form:"platform" binding:"omitempty,oneof=android ios web"`
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
func GetAd(c *gin.Context) {
	var adReq AdGetRequest
    err := c.ShouldBind(&adReq)
	if err != nil {
        fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
    if adReq.Offset == nil {
        defaultOffset := 0
        adReq.Offset = &defaultOffset
    
    }
    if adReq.Limit == nil {
        defaultLimit := 5
        adReq.Limit = &defaultLimit
    }
	c.JSON(http.StatusOK, gin.H{"message": "get"})
}

func CreateAd(c *gin.Context) {
    var adCreate AdCreationRequest
    err := c.ShouldBindJSON(&adCreate)
    if err != nil {
        fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "create",
	})
}