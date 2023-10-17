package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("35FSFJlgh4353KSFjX")

func CreateToken(userId int, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userId
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	return token.SignedString(jwtSecret)
}
func CheckToken(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if token == nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("failed to parse token claims")
	}

	userIDClaim, userOK := claims["userID"].(float64)
	roleClaim, roleOK := claims["role"].(string)

	if !userOK || !roleOK {
		return 0, "", errors.New("failed to extract user ID or role from token")
	}

	return int(userIDClaim), roleClaim, nil
}
