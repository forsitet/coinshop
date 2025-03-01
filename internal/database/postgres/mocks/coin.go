package mocks

import (
	"coin/domain"
	"errors"
	"time"
)

type InMemoryRepo struct {
	Users         map[string]domain.User
	ItemPrices    map[string]int
	Items         []domain.Item
	Ops           map[string][]domain.Operations
	autoIncUserID uint
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		Users:      make(map[string]domain.User),
		ItemPrices: map[string]int{"sword": 100, "shield": 150},
		Items: []domain.Item{
			{Name: "sword", Price: 100},
			{Name: "shield", Price: 150},
		},
		Ops:           make(map[string][]domain.Operations),
		autoIncUserID: 1,
	}
}

func (m *InMemoryRepo) GetUser(username string) (domain.User, error) {
	user, ok := m.Users[username]
	if !ok {
		return domain.User{}, nil
	}

	return user, nil
}

func (m *InMemoryRepo) CreateUser(username string, balance int) error {
	m.Users[username] = domain.User{
		ID:       m.autoIncUserID,
		Username: username,
		Balance:  balance,
	}
	m.autoIncUserID++
	return nil
}

func (m *InMemoryRepo) PostBuyItem(userID uint, itemName string) error {
	for username, user := range m.Users {
		if user.ID == userID {
			user.Inventory = append(user.Inventory, domain.InventoryItem{ItemType: itemName, Quantity: 1})
			user.Balance -= m.ItemPrices[itemName]
			m.Users[username] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *InMemoryRepo) GetItem() []domain.Item {
	return m.Items
}

func (m *InMemoryRepo) GetItemPrice(itemName string) (int, error) {
	return m.ItemPrices[itemName], nil
}

func (m *InMemoryRepo) GetOperations(username string) ([]domain.Operations, error) {
	return m.Ops[username], nil
}

func (m *InMemoryRepo) SendCoinTransaction(senderUsername, recipientUsername string, amount int) error {
	sender := m.Users[senderUsername]
	recipient := m.Users[recipientUsername]

	sender.Balance -= amount
	recipient.Balance += amount

	m.Users[senderUsername] = sender
	m.Users[recipientUsername] = recipient

	m.Ops[senderUsername] = append(m.Ops[senderUsername], domain.Operations{
		FromUser:  sender.Username,
		ToUser:    recipient.Username,
		Amount:    -amount,
		CreatedAt: time.Now(),
	})
	m.Ops[recipientUsername] = append(m.Ops[recipientUsername], domain.Operations{
		FromUser:  sender.Username,
		ToUser:    recipient.Username,
		Amount:    amount,
		CreatedAt: time.Now(),
	})
	return nil
}
