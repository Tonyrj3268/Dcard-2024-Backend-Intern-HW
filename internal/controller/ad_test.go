package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"advertisement-api/internal/dto"
	"advertisement-api/internal/model"

	"encoding/json"

	"os"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAdRepository struct {
    mock.Mock
}

func (m *MockAdRepository) GetActiveAdvertisements(now time.Time,adReq dto.AdGetRequest) ([]dto.AdGetResponse, error) {
    args := m.Called(now,adReq)
    return args.Get(0).([]dto.AdGetResponse), args.Error(1)
}

func (m *MockAdRepository) CreateAdvertisement(ad *model.Advertisement) error {
    args := m.Called(ad)
    return args.Error(0)
}
var (
    mr           *miniredis.Miniredis
    router       *gin.Engine
    mockAdRepo   *MockAdRepository
    adController *AdController
)
func TestMain(m *testing.M) {
    var err error
    mr, err = miniredis.Run()
    if err != nil {
        panic(err)
    }
    defer mr.Close()
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	mockAdRepo = new(MockAdRepository)
    redis :=  redis.NewClient(&redis.Options{
        Addr: mr.Addr(),
    })
    adController = NewAdController(mockAdRepo, redis)
	router.GET("/ad", adController.GetAd)
	router.POST("/ad", adController.CreateAd)

	os.Exit(m.Run())
}
// GetAd
func TestGetAdMultipleCases(t *testing.T) {
    mr.FlushAll()
	testCases := []struct{
		query    string
		expectedStatus int    // status code
	}{
		{"offset=0&limit=10&age=25&gender=M&country=US&platform=android", http.StatusOK},
		{"offset=0&limit=1&age=1&gender=F&country=CA&platform=ios", http.StatusOK},
		{"offset=0&limit=100&age=100&gender=M&country=JP&platform=web", http.StatusOK},
		{"offset=0&limit=101&age=25&gender=M&country=US&platform=android", http.StatusBadRequest},
		{"offset=0&limit=10&age=0&gender=F&country=GB&platform=ios", http.StatusBadRequest},
		{"offset=0&limit=10&age=25&gender=Unknown&country=US&platform=android", http.StatusBadRequest},
		{"offset=0&limit=10&age=25&gender=M&country=XX&platform=android", http.StatusBadRequest},
		{"offset=0&limit=10&age=25&gender=M&country=US&platform=alien", http.StatusBadRequest},
		{"", http.StatusBadRequest},
		{"off=0&lim=10&ag=25&gen=M&cou=US&plat=android", http.StatusBadRequest},
		{"offset=&limit=&age=&gender=&country=&platform=", http.StatusBadRequest},
		{"offset=0&limit=50&age=25&gender=M&country=US&platform=android&unknownParam=value", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Query: %s", tc.query), func(t *testing.T) {
            if tc.expectedStatus != http.StatusBadRequest {
                mockAdRepo.On("GetActiveAdvertisements",mock.AnythingOfType("time.Time"), mock.AnythingOfType("dto.AdGetRequest")).Return([]dto.AdGetResponse{}, nil)
            }
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/ad?%s", tc.query), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
            if tc.expectedStatus == http.StatusOK {
                mockAdRepo.AssertExpectations(t)
            }
		})
	}
}

func TestCreateAd(t *testing.T) {
    mr.FlushAll()
    testCases := []struct{
        name           string
        adCreate       dto.AdCreationRequest
        mockReturn     error
        expectedStatus int
    }{
    {
        name: "should return 400 if title is missing",
        adCreate: dto.AdCreationRequest{
            StartAt: time.Now(),
            EndAt:   time.Now().Add(time.Hour),
            Conditions: dto.AdCondition{
                Gender:   &[]string{"M", "F"},
                Country:  &[]string{"US", "GB"},
                Platform: &[]string{"ios", "android"},
            },
        },
        mockReturn:     nil,
        expectedStatus: http.StatusBadRequest,
    },
    {
        name: "should return 200 if conditions are missing",
        adCreate: dto.AdCreationRequest{
            Title:   "Test Ad",
            StartAt: time.Now(),
            EndAt:   time.Now().Add(time.Hour),
        },
        mockReturn:     nil,
        expectedStatus: http.StatusOK,
    },
    {
        name: "should return 200 if all fields are valid",
        adCreate: dto.AdCreationRequest{
            Title:   "Test Ad",
            StartAt: time.Now(),
            EndAt:   time.Now().Add(time.Hour),
            Conditions: dto.AdCondition{
                Gender:   &[]string{"M", "F"},
                Country:  &[]string{"US", "GB"},
                Platform: &[]string{"ios", "android"},
            },
        },
        mockReturn:     nil,
        expectedStatus: http.StatusOK,
    },
    {
        name: "should return 400 if start time is after end time",
        adCreate: dto.AdCreationRequest{
            Title:   "Test Ad",
            StartAt: time.Now().Add(time.Hour),
            EndAt:   time.Now(),
            Conditions:  dto.AdCondition{
                Gender:   &[]string{"M", "F"},
                Country:  &[]string{"US", "GB"},
                Platform: &[]string{"ios", "android"},
            },
        },
        mockReturn:     nil,
        expectedStatus: http.StatusBadRequest,
    },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.expectedStatus != http.StatusBadRequest {
                mockAdRepo.On("CreateAdvertisement", mock.AnythingOfType("*model.Advertisement")).Return(nil)
            }
            reqBody, _ := json.Marshal(tc.adCreate)
            req, _ := http.NewRequest(http.MethodPost, "/ad", bytes.NewBuffer(reqBody))
            rec := httptest.NewRecorder()
            router.ServeHTTP(rec, req)
            assert.Equal(t, tc.expectedStatus, rec.Code)
            if tc.expectedStatus == http.StatusOK {
                mockAdRepo.AssertExpectations(t)
                mockAdRepo.AssertCalled(t, "CreateAdvertisement", mock.Anything)
            }
        })
    }
}