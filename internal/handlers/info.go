package handlers

import (
	db "coin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InfoHandler(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	var user db.User
	if err := db.DB.Preload("Inventory").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  username,
		"balance":   user.Balance,
		"inventory": user.Inventory,
	})
}
