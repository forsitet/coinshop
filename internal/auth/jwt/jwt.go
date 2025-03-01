package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("supersecretkey")

type JwtClaim struct {
	Username string
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	claims := &JwtClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(mySigningKey)
}

func ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(_ *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok || !token.Valid {
		return nil, errors.New("неверный токен")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("ожидался JWT")
	}

	return claims, nil

}
