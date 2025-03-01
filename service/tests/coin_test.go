package cointests

import (
	"coin/domain"
	"coin/internal/database/postgres/mocks"
	"coin/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByUsername_NewUserCreated(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	s := service.NewCoinService(repo)

	user, err := s.GetUserByUsername("alice")
	assert.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, service.DefaultBalance, user.Balance)
}
func TestBuyItem_Success(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	repo.Users["bob"] = domain.User{ID: 1, Username: "bob", Balance: 200}
	s := service.NewCoinService(repo)

	user, err := s.BuyItem("bob", "sword")
	assert.NoError(t, err)
	assert.Equal(t, 100, user.Balance)
}

func TestBuyItem_NotEnoughCoins(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	repo.Users["carol"] = domain.User{ID: 2, Username: "carol", Balance: 50}
	s := service.NewCoinService(repo)

	_, err := s.BuyItem("carol", "shield")
	assert.ErrorIs(t, err, service.ErrCoinNotEnough)
}

func TestSendCoin_Success(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	repo.Users["dave"] = domain.User{Username: "dave", Balance: 100}
	repo.Users["eva"] = domain.User{Username: "eva", Balance: 50}
	s := service.NewCoinService(repo)

	_, err := s.SendCoin("dave", "eva", 70)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 30, repo.Users["dave"].Balance)
	assert.Equal(t, 120, repo.Users["eva"].Balance)
}

func TestGetItem(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	s := service.NewCoinService(repo)

	items := s.GetItem()
	assert.Len(t, items, 2)
}

func TestGetOperations_Success(t *testing.T) {
	repo := mocks.NewInMemoryRepo()
	repo.Users["alice"] = domain.User{Username: "alice"}
	repo.Ops["alice"] = []domain.Operations{{Amount: 42}}

	s := service.NewCoinService(repo)

	ops, err := s.GetOperations("alice")
	assert.NoError(t, err)
	assert.Len(t, ops, 1)
	assert.Equal(t, 42, ops[0].Amount)
}
