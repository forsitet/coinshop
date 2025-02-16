package database_test

import (
	"coin/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func TestUserCreation(t *testing.T) {
	mockDB := new(MockDB)
	mockDB.On("Create", mock.Anything).Return(nil)

	user := database.User{Username: "testuser"}
	err := mockDB.Create(&user)
	assert.NoError(t, err, "Ошибка создания пользователя")
	mockDB.AssertExpectations(t)
}
