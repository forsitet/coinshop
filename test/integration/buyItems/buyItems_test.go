package buytest

import (
	"coin/domain"
	hndl "coin/internal/api/http"
	"coin/internal/api/http/middleware"
	auth "coin/internal/auth/jwt"
	"coin/internal/repository/mocks"
	"coin/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSendCoinHandler(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	repo.Users = map[string]domain.User{
		"alice": {ID: 1, Username: "alice", Balance: 1000},
	}
	repo.ItemPrices = map[string]int{"tshirt": 500}
	s := service.NewCoinService(repo)
	handler := hndl.NewCoinHandler(*s)
	r := gin.Default()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.GET("/buy/:item", handler.BuyItem)
	token, _ := auth.GenerateToken("alice")

	req, _ := http.NewRequest("GET", "/api/buy/tshirt", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 500, repo.Users["alice"].Balance)
}
