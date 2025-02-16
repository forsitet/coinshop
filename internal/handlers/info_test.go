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

func setupInfoTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("Не удалось инициализировать БД:", err)
	}

	database.DB.AutoMigrate(&database.User{}, &database.InventoryItem{})
}

func TestInfoHandler_Unauthorized(t *testing.T) {
	setupInfoTestDB(t)

	r := gin.Default()
	r.GET("/api/info", InfoHandler)

	req, _ := http.NewRequest("GET", "/api/info", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestInfoHandler_Success(t *testing.T) {
	setupInfoTestDB(t)

	testUser := database.User{
		Username: "test_user",
		Balance:  500,
		Inventory: []database.InventoryItem{
			{ItemType: "sword", Quantity: 1},
			{ItemType: "shield", Quantity: 2},
		},
	}
	database.DB.Create(&testUser)

	r := gin.Default()
	r.GET("/api/info", func(c *gin.Context) {
		c.Set("username", "test_user")
		InfoHandler(c)
	})

	req, _ := http.NewRequest("GET", "/api/info", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "test_user", response["username"])
	assert.Equal(t, float64(500), response["balance"])
	assert.Len(t, response["inventory"].([]interface{}), 2)
}
