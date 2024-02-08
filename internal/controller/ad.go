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

// GetAd godoc
// @Summary Get advertisements
// @Description get advertisements by params and conditions
// @Tags advertisements
// @Accept json
// @Produce json
// @Param adGetRequest body dto.AdGetRequest true "Enter Advertisement Request Conditions"
// @Success 200 {dto.AdGetRequest} AdGetResponse "success"
// @Failure 400 {dto.AdGetRequest} json "{"error": "params error"}"
// @Failure 500 {dto.AdGetRequest} json "{"error": "server error"}"
// @Router /ad [get]
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

// CreateAd godoc
// @Summary Create advertisement
// @Description create a new advertisement
// @Tags advertisements
// @Accept json
// @Produce json
// @Param adCreationRequest body dto.AdCreationRequest true "Create Advertisement"
// @Success 200 {dto.AdCreationRequest} json "{"message": "success"}"
// @Failure 400 {dto.AdCreationRequest} json "{"error": "params error"}"
// @Failure 500 {dto.AdCreationRequest} json "{"error": "server error"}"
// @Router /ad [post]
func(a *AdController) CreateAd(c *gin.Context) {
    var adCreate dto.AdCreationRequest
    err := c.ShouldBindJSON(&adCreate)
    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    ad := model.Advertisement{
        Title:     adCreate.Title,
        StartAt:   adCreate.StartAt,
        EndAt:     adCreate.EndAt,
        AgeStart:  adCreate.Conditions.AgeStart,
        AgeEnd:    adCreate.Conditions.AgeEnd,
        // 避免 nil pointer dereference
        Gender:    assignConditionValue(adCreate.Conditions.Gender),  
        Country:   assignConditionValue(adCreate.Conditions.Country),
        Platform:  assignConditionValue(adCreate.Conditions.Platform),
    }
    if err := a.adRepository.CreateAdvertisement(&ad); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func assignConditionValue(condition *[]string) []string {
    if condition != nil {
        return *condition
    }
    return []string{}
}