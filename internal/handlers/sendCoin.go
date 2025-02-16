package handlers

import (
	db "coin/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SendCoinHandler(c *gin.Context) {
	senderUsername, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизайии"})
		return
	}

	var req struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Количетсво койнов должно быть больше 0"})
		return
	}

	var sender, recipient db.User

	if err := db.DB.Where("username = ?", senderUsername).First(&sender).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
		return
	}

	if sender.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нехватает койнов :("})
		return
	}

	if err := db.DB.Where("username = ?", req.ToUser).First(&recipient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Получатель не найден"})
		return
	}

	if sender.Username == recipient.Username {
		c.JSON(http.StatusNotFound, gin.H{"error": "Нельзя отправлять койны самому себе"})
		return
	}

	sender.Balance -= req.Amount
	recipient.Balance += req.Amount

	tx := db.DB.Begin()

	if err := tx.Save(&sender).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправления"})
		return
	}

	if err := tx.Save(&recipient).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправления"})
		return
	}

	transaction := db.Transaction{
		FromUser:  sender.ID,
		ToUser:    recipient.ID,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания записи"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"sender_new_balance": sender.Balance,
		"recipient":          recipient.Username,
		"amount":             req.Amount,
	})
}
