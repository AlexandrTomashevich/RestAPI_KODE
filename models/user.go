package models

import (
	"RestAPI_KODE/lib"
	"RestAPI_KODE/utils"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

func ValidateUser(db *sql.DB, username, password string) (User, bool) {
	user := User{}
	err := db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		lib.Errorf("Failed to validate user: %v", err)
		return User{}, false
	}
	if utils.CheckPassword(user.PasswordHash, password) {
		return user, true
	}
	return User{}, false
}

func AuthUser(db *sql.DB, username, password string) (*User, error) {
	fmt.Println("Entered AuthUser")
	user, valid := ValidateUser(db, username, password)
	if valid {
		return &user, nil
	}
	fmt.Println("Exiting AuthUser")
	return &User{}, errors.New("invalid credentials")
}

func CreateUser(db *sql.DB, username, password, role string) (int, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		lib.Errorf("Failed to hash password: %v", err)
		return 0, err
	}
	var userID int
	err = db.QueryRow("INSERT INTO users(username, password_hash, role) VALUES($1, $2, $3) RETURNING id", username, hashedPassword, role).Scan(&userID)
	if err != nil {
		lib.Errorf("Failed to create user: %v", err)
		return 0, err
	}
	return userID, nil
}
