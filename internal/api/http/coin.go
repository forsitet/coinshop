package http

import (
	"coin/domain"
	mwAuth "coin/internal/api/http/middleware"
	auth "coin/internal/auth/jwt"
	"coin/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CoinHandler struct {
	service service.CoinService
}

func NewCoinHandler(service service.CoinService) *CoinHandler {
	return &CoinHandler{
		service: service,
	}
}

func (h *CoinHandler) Auth(c *gin.Context) {
	op := "http.coin.Auth"
	var rawUser domain.User
	if err := c.ShouldBindJSON(&rawUser); err != nil {
		log.Println(op, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	user, err := h.service.GetUserByUsername(rawUser.Username)
	if err != nil {
		log.Println(op, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check user"})
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		log.Println(op, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("jwt_token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (h *CoinHandler) BuyItem(c *gin.Context) {
	op := "http.coin.BuyItem"
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		log.Println(op, "Error converting username")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка преобразования имени пользователя"})
		return
	}
	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя товара обязательно"})
		return
	}
	user, err := h.service.BuyItem(usernameStr, itemName)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Такого товара нет"})
			return
		}
		if errors.Is(err, service.ErrCoinNotEnough) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не достаточно койнов"})
			return
		}
		log.Println(op, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not buy item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": user.Balance, "inventory": user.Inventory})
}

func (h *CoinHandler) GetItems(c *gin.Context) {
	items := h.service.GetItem()
	c.JSON(http.StatusOK, items)
}

func (h *CoinHandler) Info(c *gin.Context) {
	op := "http.coin.Info"
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		log.Println(op, "Error converting username")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка преобразования имени пользователя"})
		return
	}
	user, err := h.service.GetUserByUsername(usernameStr)
	if err != nil {
		log.Println(op, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username":  username,
		"balance":   user.Balance,
		"inventory": user.Inventory,
	})
}

func (h *CoinHandler) SendCoin(c *gin.Context) {
	op := "http.coin.SendCoin"
	var req struct {
		ToUser string `json:"to_user"`
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

	senderUsername, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}
	senderUsernameStr, ok := senderUsername.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка преобразования имени пользователя"})
		return
	}

	if senderUsernameStr == req.ToUser {
		c.JSON(http.StatusNotFound, gin.H{"error": "Нельзя отправлять койны самому себе"})
		return
	}

	balanceSender, err := h.service.SendCoin(senderUsernameStr, req.ToUser, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrCoinNotEnough) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недостаточно койнов"})
			return
		}
		log.Println(op, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send coin"})
		return
	}
	log.Println(op, "Coin sent successfully")
	c.JSON(http.StatusOK, gin.H{
		"to_user":             req.ToUser,
		"amount":             req.Amount,
		"sender_new_balance": balanceSender,
	})
}

func (h *CoinHandler) Operations(c *gin.Context) {
	const op = "http.coin.Operations"
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		log.Println(op, "Error getting username from context")
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка преобразования имени пользователя"})
		return
	}
	operations, err := h.service.GetOperations(usernameStr)
	if err != nil {
		log.Println(op, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get operations"})
		return
	}
	c.JSON(http.StatusOK, operations)
	log.Println(op, "Operations retrieved successfully")

}

func (h *CoinHandler) WithObjectHandlers(r *gin.Engine) {
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static")
	})
	api := r.Group("/api")
	{
		api.POST("/auth", h.Auth)
		api.GET("/info", mwAuth.AuthMiddleware, h.Info)
		api.POST("/sendCoin", mwAuth.AuthMiddleware, h.SendCoin)
		api.GET("/buy/:item", mwAuth.AuthMiddleware, h.BuyItem)
		api.GET("/items", h.GetItems)
		api.GET("/transactions", mwAuth.AuthMiddleware, h.Operations)
	}

}
