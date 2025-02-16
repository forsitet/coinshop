package handlers

import (
	db "coin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuyItemHandler(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя товара обязательно"})
		return
	}

	price, ok := db.ItemPrices[itemName]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Такого товара нет"})
		return
	}

	var user db.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if user.Balance < price {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не достаточно койнов"})
		return
	}

	user.Balance -= price
	db.DB.Save(&user)

	var item db.InventoryItem
	if err := db.DB.Where("user_id = ? AND item_type = ?", user.ID, itemName).First(&item).Error; err != nil {
		item = db.InventoryItem{UserID: user.ID, ItemType: itemName, Quantity: 1}
		db.DB.Create(&item)
	} else {
		item.Quantity++
		db.DB.Save(&item)
	}

	c.JSON(http.StatusOK, gin.H{"balance": user.Balance, "inventory": item})
}
