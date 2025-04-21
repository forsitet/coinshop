package repository

import "coin/domain"

type Object interface {
	GetUser(username string) (domain.User, error)
	CreateUser(username string, balance int) error
	PostBuyItem(userID uint, itemType string) error
	GetItems() ([]domain.Item, error)
	GetItemPrice(itemName string) (int, error)
	GetOperations(username string) ([]domain.Operations, error)
	SendCoinTransaction(senderUsername, recipientUsername string, amount int) error
}
