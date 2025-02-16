package auth_test

import (
	"coin/auth"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := auth.GenerateToken("testuser")
	assert.NoError(t, err, "Ошибка генерации токена")
	assert.NotEmpty(t, token, "Токен не должен быть пустым")
}

func TestValidateToken_ValidToken(t *testing.T) {
	token, err := auth.GenerateToken("testuser")
	assert.NoError(t, err)

	claims, err := auth.ValidateToken(token)
	assert.NoError(t, err, "Ошибка валидации токена")
	assert.Equal(t, "testuser", claims.Username, "Неверное имя пользователя")
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	claims := &auth.JwtClaim{
		Username: "testuser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("supersecretkey"))
	assert.NoError(t, err)

	_, err = auth.ValidateToken(signedToken)
	assert.Error(t, err, "Ожидалась ошибка при проверке просроченного токена")
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := auth.ValidateToken("invalidtoken")
	assert.Error(t, err, "Ожидалась ошибка при проверке неверного токена")
}
