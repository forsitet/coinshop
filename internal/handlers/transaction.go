package handlers

import (
	db "coin/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TransactionHandler(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	var user db.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var transactions []struct {
		CreatedAt time.Time `json:"created_at"`
		Send      string    `json:"send"`
		Rec       string    `json:"rec"`
		Amount    int       `json:"amount"`
	}

	query := `
		SELECT t.amount, t.created_at, u1.username AS send, u2.username AS rec
		FROM transactions t
		JOIN users u1 ON t.from_user = u1.id
		JOIN users u2 ON t.to_user = u2.id
		WHERE t.from_user = ? OR t.to_user = ?
		ORDER BY t.created_at DESC
		`

	if err := db.DB.Raw(query, user.ID, user.ID).Scan(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
