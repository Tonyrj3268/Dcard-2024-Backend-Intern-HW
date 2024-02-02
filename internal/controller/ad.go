package controller

import (
	"advertisement-api/internal/model"
	"advertisement-api/internal/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdController struct {
    adRepository repository.AdRepository
}

func NewAdController(adRepo repository.AdRepository) *AdController {
    return &AdController{
        adRepository: adRepo,
    }
}

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

func(a *AdController) GetAd(c *gin.Context) {
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
    ads, err := a.adRepository.GetActiveAdvertisements(adReq.Age, adReq.Gender, adReq.Country, adReq.Platform, *adReq.Offset, *adReq.Limit)
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, ads)
}

func(a *AdController) CreateAd(c *gin.Context) {
    var adCreate AdCreationRequest
    err := c.ShouldBindJSON(&adCreate)
    if err != nil {
        fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
    // 避免 nil pointer dereference
    gender, country, platform := []string{}, []string{}, []string{}
    if adCreate.Conditions.Gender != nil {
        gender = *adCreate.Conditions.Gender
    }
    if adCreate.Conditions.Country != nil {
        country = *adCreate.Conditions.Country
    }
    if adCreate.Conditions.Platform != nil {
        platform = *adCreate.Conditions.Platform
    }

    ad := model.Advertisement{
        Title:     adCreate.Title,
        StartAt:   adCreate.StartAt,
        EndAt:     adCreate.EndAt,
        Gender:    gender,
        Country:   country,
        Platform:  platform,
    }
    if err := a.adRepository.CreateAdvertisement(&ad); err != nil {
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
        c.JSON(http.StatusOK, gin.H{"message": "success"})
}