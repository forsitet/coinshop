package main

import (
	db "coin/database"
	hndl "coin/internal/handlers"
	mwAuth "coin/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static")
	})
	api := r.Group("/api")
	{
		api.POST("/auth", hndl.AuthHandler)
		api.GET("/info", mwAuth.AuthMiddleware, hndl.InfoHandler)
		api.POST("/sendCoin", mwAuth.AuthMiddleware, hndl.SendCoinHandler)
		api.GET("/buy/:item", mwAuth.AuthMiddleware, hndl.BuyItemHandler)
		api.GET("/items", hndl.GetItemsHandler)
		api.GET("/transactions", mwAuth.AuthMiddleware, hndl.TransactionHandler)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
