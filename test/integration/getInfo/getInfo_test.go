package getinfotest

import (
	"coin/domain"
	hndl "coin/internal/api/http"
	"coin/internal/api/http/middleware"
	auth "coin/internal/auth/jwt"
	"coin/internal/repository/mocks"
	"coin/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInfoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewInMemoryRepo()
	repo.Users = map[string]domain.User{
		"alice": {ID: 1, Username: "alice", Balance: 150, Inventory: []domain.InventoryItem{
			{ItemType: "book", Quantity: 2},
		}},
	}
	s := service.NewCoinService(repo)
	handler := hndl.NewCoinHandler(*s)

	r := gin.Default()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.GET("/info", handler.Info)

	token, err := auth.GenerateToken("alice")
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "alice", resp["username"])
	assert.EqualValues(t, 150, resp["balance"])

	inventory := resp["inventory"].([]interface{})
	assert.Len(t, inventory, 1)

	item := inventory[0].(map[string]interface{})
	assert.Equal(t, "book", item["item_type"])
	assert.EqualValues(t, 2, item["quantity"])
}
