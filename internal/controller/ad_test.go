package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// GetAd
func TestGetAdMultipleCases(t *testing.T) {
	testCases := []struct{
		query    string
		expected int    // status code
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

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/ad", GetAd)

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Query: %s", tc.query), func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/ad?%s", tc.query), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
func TestCreateAd(t *testing.T) {
    router := gin.Default()
    router.POST("/ad", CreateAd)

	age_18 := 18
	age_50 := 50
    testCases := []struct {
        payload     AdCreationRequest
        expectedCode int
    }{
        {
            payload: AdCreationRequest{
                Title: "Summer Sale",
                StartAt: time.Now(),
                EndAt: time.Now().Add(24 * time.Hour),
                Conditions: AdCondition{
                    AgeStart:      &age_18,
					AgeEnd:        &age_50,
                    Gender:   &[]string{"M", "F"},
                    Country:  &[]string{"TW", "JP"},
                    Platform: &[]string{"android", "ios"},
                },
            },
            expectedCode: http.StatusOK,
        },
        {
            payload: AdCreationRequest{
                StartAt: time.Now(),
                EndAt: time.Now().Add(24 * time.Hour),
            },
            expectedCode: http.StatusBadRequest,
        },
        {
            payload: AdCreationRequest{
                Title: "Winter Sale",
                StartAt: time.Now().Add(24 * time.Hour), // 開始比結束晚
                EndAt: time.Now(),
                Conditions: AdCondition{
                    AgeStart:      &age_18,
					AgeEnd:        &age_50,
                    Gender:   &[]string{"M"},
                    Country:  &[]string{"US"},
                    Platform: &[]string{"web"},
                },
            },
            expectedCode: http.StatusBadRequest,
        },
		{
            payload: AdCreationRequest{
                Title: "Summer Sale",
                StartAt: time.Now(),
                EndAt: time.Now().Add(24 * time.Hour),
                Conditions: AdCondition{
                    AgeStart:      &age_50, // 年齡開始比結束晚
					AgeEnd:        &age_18,
                    Gender:   &[]string{"M", "F"},
                    Country:  &[]string{"TW", "JP"},
                    Platform: &[]string{"android", "ios"},
                },
            },
            expectedCode: http.StatusBadRequest,
        },
    }

    for _, tc := range testCases {
        t.Run("",func(t *testing.T) {
            body, _ := json.Marshal(tc.payload)
            req := httptest.NewRequest("POST", "/ad", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedCode, w.Code)
        })
    }
}