package sendtest

import (
	"bytes"
	"coin/domain"
	hndl "coin/internal/api/http"
	"coin/internal/api/http/middleware"
	auth "coin/internal/auth/jwt"
	"coin/internal/database/postgres/mocks"
	"coin/service"
	"encoding/json"
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
		"bob":   {ID: 2, Username: "bob", Balance: 500}}
	s := service.NewCoinService(repo)
	handler := hndl.NewCoinHandler(*s)

	r := gin.Default()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.POST("/sendCoin", handler.SendCoin)
	token, _ := auth.GenerateToken("alice")

	payload := map[string]interface{}{
		"to_user": "bob",
		"amount":  100,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 900, repo.Users["alice"].Balance)
	assert.Equal(t, 600, repo.Users["bob"].Balance)
}
