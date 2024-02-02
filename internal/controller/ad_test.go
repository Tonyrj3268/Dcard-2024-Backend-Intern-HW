package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"advertisement-api/internal/model"

	"encoding/json"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAdRepository struct {
    mock.Mock
}

func (m *MockAdRepository) GetActiveAdvertisements(age *int, gender, country, platform *string, offset, limit int) ([]struct{Title string; EndAt time.Time}, error) {
    args := m.Called(age, gender, country, platform, offset, limit)
    return args.Get(0).([]struct{Title string; EndAt time.Time}), args.Error(1)
}

func (m *MockAdRepository) CreateAdvertisement(ad *model.Advertisement) error {
    args := m.Called(ad)
    return args.Error(0)
}

var router *gin.Engine
var mockAdRepo *MockAdRepository
var adController *AdController

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	mockAdRepo = new(MockAdRepository)
    adController = NewAdController(mockAdRepo)
	router.GET("/ad", adController.GetAd)
	router.POST("/ad", adController.CreateAd)

	os.Exit(m.Run())
}
// GetAd
func TestGetAdMultipleCases(t *testing.T) {

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
		{"", http.StatusOK},
		{"off=0&lim=10&ag=25&gen=M&cou=US&plat=android", http.StatusOK},
		{"offset=&limit=&age=&gender=&country=&platform=", http.StatusBadRequest},
		{"offset=0&limit=50&age=25&gender=M&country=US&platform=android&unknownParam=value", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Query: %s", tc.query), func(t *testing.T) {
            if tc.expectedStatus != http.StatusBadRequest {
                mockAdRepo.On("GetActiveAdvertisements", mock.AnythingOfType("*int"),mock.AnythingOfType("*string"),mock.AnythingOfType("*string"),mock.AnythingOfType("*string"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return([]struct{ Title string; EndAt time.Time }{}, nil)
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
    testCases := []struct{
        name           string
        adCreate       AdCreationRequest
        mockReturn     error
        expectedStatus int
    }{
    {
        name: "should return 400 if title is missing",
        adCreate: AdCreationRequest{
            StartAt: time.Now(),
            EndAt:   time.Now().Add(time.Hour),
            Conditions: AdCondition{
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
        adCreate: AdCreationRequest{
            Title:   "Test Ad",
            StartAt: time.Now(),
            EndAt:   time.Now().Add(time.Hour),
        },
        mockReturn:     nil,
        expectedStatus: http.StatusOK,
    },
    {
        name: "should return 200 if all fields are valid",
        adCreate: AdCreationRequest{
            Title:   "Test Ad",
            StartAt: time.Now(),
            EndAt:   time.Now().Add(time.Hour),
            Conditions: AdCondition{
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
        adCreate: AdCreationRequest{
            Title:   "Test Ad",
            StartAt: time.Now().Add(time.Hour),
            EndAt:   time.Now(),
            Conditions:  AdCondition{
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