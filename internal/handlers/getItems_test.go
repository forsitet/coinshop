package handlers

import (
	"coin/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Ошибка подключения к БД")
	}
	database.DB.AutoMigrate(&database.Item{})
}

func TestGetItemsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	database.DB.Create(&database.Item{Name: "cup", Price: 100})
	database.DB.Create(&database.Item{Name: "book", Price: 150})

	r := gin.Default()
	r.GET("/api/items", GetItemsHandler)

	req, _ := http.NewRequest("GET", "/api/items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var items []database.Item
	err := json.Unmarshal(w.Body.Bytes(), &items)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "cup", items[0].Name)
	assert.Equal(t, 100, items[0].Price)
	assert.Equal(t, "book", items[1].Name)
	assert.Equal(t, 150, items[1].Price)
}
