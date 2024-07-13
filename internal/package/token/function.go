package token

import (
	"api/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateJWT(userid, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
