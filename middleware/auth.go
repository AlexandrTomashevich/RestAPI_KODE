package middleware

import (
	"RestAPI_KODE/models"
	"RestAPI_KODE/utils"
)

func GenerateTokenForUser(user models.User) (string, error) {

	token, err := utils.CreateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func AuthenticateUser(username, passwordHash string) (*models.User, error) {
	return models.AuthUser(username, passwordHash)
}

func CreateUser(username, passwordHash, role string) (int, error) {
	return models.CreateUser(username, passwordHash, role)
}
