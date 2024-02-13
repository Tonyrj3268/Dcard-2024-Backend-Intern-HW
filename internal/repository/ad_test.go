package repository

import (
	"database/sql"
	"testing"

	"regexp"
	"time"

	"os"

	"advertisement-api/internal/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var repo *adRepository
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
    var err error
    var db *sql.DB
    db, mock, err = sqlmock.New()
    if err != nil {
		panic("failed to create sqlmock")
	}
	defer db.Close()
    mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.31"))

    gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic("failed to open gorm db")
	}

	repo = NewAdRepository(gormDB)

	os.Exit(m.Run())
}

func TestGetActiveAdvertisementsBasic(t *testing.T) {

    now := time.Now()
    rows := sqlmock.NewRows([]string{"title", "start_at", "end_at"}).
        AddRow("Ad 1", now, now).
        AddRow("Ad 2", now.Add(-time.Minute * 2), now.Add(time.Minute * 2))

    mock.ExpectQuery(regexp.QuoteMeta(
        "SELECT title, end_at FROM `advertisements` WHERE ? BETWEEN start_at AND end_at ORDER BY end_at ASC LIMIT 10 OFFSET 20")).
        WithArgs(now).
        WillReturnRows(rows)
        
    adReq := dto.AdGetRequest{Offset: 20, Limit: 10}
    ads, err := repo.GetActiveAdvertisements(now, adReq)
    assert.NoError(t, err)
    assert.Len(t, ads, 2)
    assert.Equal(t, ads, []dto.AdGetResponse{{Title: "Ad 1", EndAt:now}, {Title:"Ad 2", EndAt:now.Add(time.Minute * 2)}})
}

func TestGetActiveAdvertisementsWithOptions(t *testing.T) {
    now := time.Now()

    rows := sqlmock.NewRows([]string{"title", "end_at"}).
        AddRow("Ad 1", now)
    mock.ExpectQuery(regexp.QuoteMeta("SELECT title, end_at FROM `advertisements` WHERE (? BETWEEN start_at AND end_at) AND (? BETWEEN age_start AND age_end) AND FIND_IN_SET(?,gender)>0 AND FIND_IN_SET(?,country)>0 AND FIND_IN_SET(?,platform)>0  ORDER BY end_at ASC LIMIT 1")). 
        WithArgs(now, 25, "M", "US", "web").
        WillReturnRows(rows)
    age := 25
    gender := "M"
    country := "US"
    platform := "web"
    adReq := dto.AdGetRequest{Offset: 0, Limit:1, Age: &age, Gender: &gender, Country: &country, Platform: &platform}
    ads, err := repo.GetActiveAdvertisements(now,adReq)
    assert.NoError(t, err)
    assert.Len(t, ads, 1)
}