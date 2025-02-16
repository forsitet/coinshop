package handlers

import (
	"coin/database"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSendCoinTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("Не удалось инициализировать БД:", err)
	}

	database.DB.AutoMigrate(&database.User{}, &database.Transaction{})
}

func TestSendCoinHandler_Unauthorized(t *testing.T) {
	setupSendCoinTestDB(t)

	r := gin.Default()
	r.POST("/api/sendCoin", SendCoinHandler)

	req, _ := http.NewRequest("POST", "/api/sendCoin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSendCoinHandler_InvalidJSON(t *testing.T) {
	setupSendCoinTestDB(t)

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "receiver", "amount": "abc"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendCoinHandler_InvalidAmount(t *testing.T) {
	setupSendCoinTestDB(t)

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "receiver", "amount": 0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendCoinHandler_SenderNotFound(t *testing.T) {
	setupSendCoinTestDB(t)

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "nonexistent_sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "receiver", "amount": 100}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSendCoinHandler_NotEnoughCoins(t *testing.T) {
	setupSendCoinTestDB(t)

	database.DB.Create(&database.User{Username: "sender", Balance: 50})
	database.DB.Create(&database.User{Username: "receiver", Balance: 100})

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "receiver", "amount": 100}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendCoinHandler_RecipientNotFound(t *testing.T) {
	setupSendCoinTestDB(t)

	database.DB.Create(&database.User{Username: "sender", Balance: 100})

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "nonexistent_user", "amount": 50}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSendCoinHandler_SelfTransfer(t *testing.T) {
	setupSendCoinTestDB(t)

	database.DB.Create(&database.User{Username: "sender", Balance: 100})

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "sender", "amount": 50}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSendCoinHandler_Success(t *testing.T) {
	setupSendCoinTestDB(t)

	sender := database.User{Username: "sender", Balance: 200}
	receiver := database.User{Username: "receiver", Balance: 100}

	database.DB.Create(&sender)
	database.DB.Create(&receiver)

	r := gin.Default()
	r.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("username", "sender")
		SendCoinHandler(c)
	})

	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(`{"toUser": "receiver", "amount": 50}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedSender, updatedReceiver database.User
	database.DB.First(&updatedSender, "username = ?", "sender")
	database.DB.First(&updatedReceiver, "username = ?", "receiver")

	assert.Equal(t, 150, updatedSender.Balance)
	assert.Equal(t, 150, updatedReceiver.Balance)

	var transaction database.Transaction
	database.DB.First(&transaction)

	assert.Equal(t, sender.ID, transaction.FromUser)
	assert.Equal(t, receiver.ID, transaction.ToUser)
	assert.Equal(t, 50, transaction.Amount)
}
