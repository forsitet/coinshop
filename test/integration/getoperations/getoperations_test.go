package operationstest

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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOperationsHandler_Success(t *testing.T) {

	repo := mocks.NewInMemoryRepo()
	repo.Users["alice"] = domain.User{ID: 1, Username: "alice", Balance: 100}
	repo.Users["bob"] = domain.User{ID: 2, Username: "bob", Balance: 100}

	repo.Ops["alice"] = []domain.Operations{
		{
			CreatedAt: time.Date(2025, 4, 16, 20, 0, 0, 0, time.UTC),
			ID:        1,
			FromUser:  "alice",
			ToUser:    "bob",
			Amount:    10,
		},
	}

	s := service.NewCoinService(repo)
	handler := hndl.NewCoinHandler(*s)

	r := gin.Default()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.GET("/transactions", handler.Operations)

	token, err := auth.GenerateToken("alice")
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/api/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var ops []domain.Operations
	err = json.Unmarshal(w.Body.Bytes(), &ops)
	assert.NoError(t, err)

	assert.Len(t, ops, 1)
	assert.Equal(t, "alice", ops[0].FromUser)
	assert.Equal(t, "bob", ops[0].ToUser)
	assert.Equal(t, 10, ops[0].Amount)
	assert.Equal(t, uint(1), ops[0].ID)
	assert.Equal(t, "2025-04-16T20:00:00Z", ops[0].CreatedAt.Format(time.RFC3339))
}
