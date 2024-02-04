package controller

import (
	"advertisement-api/internal/dto"
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

func(a *AdController) GetAd(c *gin.Context) {
	var adReq dto.AdGetRequest
    err := c.ShouldBind(&adReq)
	if err != nil {
        fmt.Println("err")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
    ads, err := a.adRepository.GetActiveAdvertisements(time.Now(), adReq)
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, ads)
}

func(a *AdController) CreateAd(c *gin.Context) {
    var adCreate dto.AdCreationRequest
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
        AgeStart:  adCreate.Conditions.AgeStart,
        AgeEnd:    adCreate.Conditions.AgeEnd,
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