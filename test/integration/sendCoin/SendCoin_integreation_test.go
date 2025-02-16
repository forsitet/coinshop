package sendcoin

import (
	"bytes"
	"coin/auth"
	db "coin/database"
	hndl "coin/internal/handlers"
	mwAuth "coin/internal/middleware"
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
	db.DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.DB.AutoMigrate(&db.User{}, &db.Transaction{})
	db.DB.Create(&db.User{Username: "sender", Balance: 1000})
	db.DB.Create(&db.User{Username: "recipient", Balance: 500})
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api")
	api.POST("/sendCoin", mwAuth.AuthMiddleware, hndl.SendCoinHandler)
	return r
}

func TestSendCoinHandler(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	token, _ := auth.GenerateToken("sender")

	payload := map[string]interface{}{
		"toUser": "recipient",
		"amount": 100,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
