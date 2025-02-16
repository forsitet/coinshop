package handlers

import (
	"coin/database"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBuyItemsTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("Не удалось инициализировать БД:", err)
	}

	database.DB.AutoMigrate(&database.User{}, &database.InventoryItem{})

	database.ItemPrices["book"] = 20
}

func TestBuyItemHandler_Unauthorized(t *testing.T) {
	setupBuyItemsTestDB(t)

	r := gin.Default()
	r.GET("/api/buy/:item", BuyItemHandler)

	req, _ := http.NewRequest("GET", "/api/buy/book", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBuyItemHandler_MissingItem(t *testing.T) {
	setupBuyItemsTestDB(t)

	r := gin.Default()
	r.GET("/api/buy/:item", func(c *gin.Context) {
		c.Set("username", "player1")
		BuyItemHandler(c)
	})

	req, _ := http.NewRequest("GET", "/api/buy/item", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyItemHandler_ItemNotFound(t *testing.T) {
	setupBuyItemsTestDB(t)

	r := gin.Default()
	r.GET("/api/buy/:item", func(c *gin.Context) {
		c.Set("username", "player1")
		BuyItemHandler(c)
	})

	req, _ := http.NewRequest("GET", "/api/buy/shield", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyItemHandler_Success_NewItem(t *testing.T) {
	setupBuyItemsTestDB(t)

	user := database.User{Username: "player1", Balance: 200}
	database.DB.Create(&user)

	r := gin.Default()
	r.GET("/api/buy/:item", func(c *gin.Context) {
		c.Set("username", "player1")
		BuyItemHandler(c)
	})

	req, _ := http.NewRequest("GET", "/api/buy/book", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedUser database.User
	database.DB.First(&updatedUser, "username = ?", "player1")
	assert.Equal(t, 180, updatedUser.Balance)

	var item database.InventoryItem
	database.DB.First(&item, "user_id = ? AND item_type = ?", updatedUser.ID, "book")
	assert.Equal(t, 1, item.Quantity)
}
