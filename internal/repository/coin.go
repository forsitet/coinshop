package repository

import (
	"coin/domain"
	"coin/internal/database/postgres"
	"coin/internal/database/redis"
)

type CoinRepository struct {
	Postgres *postgres.PostgresStore
	Redis    *redis.RedisClient
}

func NewCoinRepository(pg *postgres.PostgresStore, rds *redis.RedisClient) *CoinRepository {
	return &CoinRepository{
		Postgres: pg,
		Redis:    rds,
	}
}

func (r *CoinRepository) CreateUser(username string, balance int) error {
	return r.Postgres.CreateUser(username, balance)
}

func (r *CoinRepository) GetUser(username string) (domain.User, error) {
	return r.Postgres.GetUser(username)
}

func (r *CoinRepository) SendCoinTransaction(senderUsername, recipientUsername string, amount int) error {
	return r.Postgres.SendCoinTransaction(senderUsername, recipientUsername, amount)
}

func (r *CoinRepository) PostBuyItem(userID uint, itemName string) error {
	price, err := r.Redis.GetItemPrice(itemName)
	if err != nil {
		return err
	}
	return r.Postgres.PostBuyItem(userID, itemName, price)
}

func (r *CoinRepository) GetOperations(username string) ([]domain.Operations, error) {
	return r.Postgres.GetOperations(username)
}

func (r *CoinRepository) GetItems() ([]domain.Item, error) {
	return r.Redis.GetItems()
}

func (r *CoinRepository) GetItemPrice(itemName string) (int, error) {
	return r.Redis.GetItemPrice(itemName)
}
