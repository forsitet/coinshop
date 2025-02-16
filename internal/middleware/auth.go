package middleware

import (
	"coin/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходим токен"})
		c.Abort()
		return
	}
	tokenString = tokenString[7:]
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		c.Abort()
		return
	}

	c.Set("username", claims.Username)
	c.Next()
}
