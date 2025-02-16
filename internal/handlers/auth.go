package handlers

import (
	"coin/auth"
	db "coin/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHandler(c *gin.Context) {
	var user db.User

	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	var existingUser db.User

	db.DB.Where("username = ?", user.Username).First(&existingUser)

	if existingUser.ID == 0 {
		user.Balance = 1000
		db.DB.Create(&user)
	} else {
		user = existingUser
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("jwt_token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token})
}
