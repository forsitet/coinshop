package handlers

import (
	db "coin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetItemsHandler(c *gin.Context) {
	var items []db.Item
	if err := db.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренная ошибка"})
		return
	}
	c.JSON(http.StatusOK, items)
}
