package auth

import (
	"go-book-api/src/config"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateToken(userID uint) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userId"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte(config.SecretKey))
}
