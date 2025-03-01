package service

import (
	"coin/domain"
	"coin/internal/database"
	"log"
)

type CoinService struct {
	repo database.Object
}

func NewCoinService(repo database.Object) *CoinService {
	return &CoinService{
		repo: repo,
	}
}

const DefaultBalance = 1000

func (s *CoinService) GetUserByUsername(username string) (domain.User, error) {
	const op = "service.coin.GetUserByUsername"
	user, err := s.repo.GetUser(username)
	if err != nil {
		return domain.User{}, err
	}
	if user.ID == 0 && user.Username == "" && user.Balance == 0 && len(user.Inventory) == 0 {
		err = s.СreateUser(username)
		if err != nil {
			log.Println(op, err)
			return domain.User{}, err
		}
		user, err = s.repo.GetUser(username)
		if err != nil {
			log.Println(op, err)
			return domain.User{}, err
		}
	}
	return user, nil
}

func (s *CoinService) СreateUser(username string) error {
	const op = "service.coin.createUser"
	err := s.repo.CreateUser(username, DefaultBalance)
	if err != nil {
		log.Println(op, err)
		return err
	}
	_, err = s.GetUserByUsername(username)
	if err != nil {
		log.Println(op, err)
		return err
	}
	return nil
}

func (s *CoinService) BuyItem(username string, itemName string) (domain.User, error) {
	const op = "service.coin.BuyItem"
	price, err := s.repo.GetItemPrice(itemName)
	if err != nil {
		return domain.User{}, ErrItemNotFound
	}

	user, err := s.GetUserByUsername(username)
	if err != nil {
		log.Println(op, err)
		return domain.User{}, err
	}
	if user.Balance < price {
		return domain.User{}, ErrCoinNotEnough
	}
	user.Balance -= price
	err = s.repo.PostBuyItem(user.ID, itemName)
	if err != nil {
		log.Println(op, err)
		return domain.User{}, err
	}

	return user, nil
}

func (s *CoinService) GetItem() []domain.Item {
	return s.repo.GetItem()
}

func (s *CoinService) SendCoin(senderUsername string, recipientUsername string, amount int) (int, error) {
	const op = "service.coin.SendCoin"
	sender, err := s.GetUserByUsername(senderUsername)
	if err != nil {
		log.Println(op, err)
		return sender.Balance, err
	}
	log.Println(senderUsername, recipientUsername)
	recipient, err := s.GetUserByUsername(recipientUsername)
	if err != nil {
		log.Println(op, err)
		return sender.Balance, err
	}

	if sender.Balance <= 0 {
		return sender.Balance, ErrCoinNotEnough
	}

	if sender.Balance < amount {
		return sender.Balance, ErrCoinNotEnough
	}
	err = s.repo.SendCoinTransaction(sender.Username, recipient.Username, amount)
	if err != nil {
		log.Println(op, err)
		return sender.Balance, err
	}
	return sender.Balance - amount, nil
}

func (s *CoinService) GetOperations(username string) ([]domain.Operations, error) {
	const op = "service.coin.GetOperations"
	user, err := s.GetUserByUsername(username)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	operations, err := s.repo.GetOperations(user.Username)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return operations, nil
}
