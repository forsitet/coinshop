package buyitems

import (
	"coin/auth"
	db "coin/database"
	hndl "coin/internal/handlers"
	mwAuth "coin/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db.DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.DB.AutoMigrate(&db.User{}, &db.InventoryItem{})

	db.ItemPrices["tshirt"] = 500

	db.DB.Create(&db.User{Username: "testuser", Balance: 1000})
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api")
	api.GET("/buy/:item", mwAuth.AuthMiddleware, hndl.BuyItemHandler)
	return r
}

func TestBuyItemHandler(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	token, _ := auth.GenerateToken("testuser")

	req, _ := http.NewRequest("GET", "/api/buy/tshirt", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
