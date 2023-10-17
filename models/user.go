package models

import (
	"RestAPI_KODE/database"
	"RestAPI_KODE/lib"
	"RestAPI_KODE/utils"
	"errors"
	"log"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

func UserExists(userID int) (bool, error) {
	if database.DB == nil {
		return false, errors.New("database connection not initialized")
	}

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ValidateUser(username, password string) (*User, bool) {
	user := &User{}
	err := database.DB.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username=$1",
		username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)

	if err != nil {
		log.Printf("Failed to validate user: %v", err)
		return nil, false
	}
	if utils.CheckPassword(user.PasswordHash, password) {
		return user, true
	}
	return nil, false
}

func AuthUser(username, password string) (*User, error) {
	lib.Logger.Printf("Entered AuthUser")
	user, valid := ValidateUser(username, password)
	switch {
	case user == nil:
		return nil, errors.New("user not found or invalid user received from ValidateUser")
	case !valid:
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func CreateUser(username, password, role string) (int, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		lib.Errorf("Failed to hash password: %v", err)
		return 0, err
	}
	var userID int
	err = database.DB.QueryRow("INSERT INTO users(username, password_hash, role) VALUES($1, $2, $3) RETURNING id",
		username, hashedPassword, role).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil

}
