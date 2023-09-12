package middleware

import (
	"RestAPI_KODE/utils"
	"fmt"
	"strings"
)

// Проверка токена
func CheckToken(token string) string {
	_, role, err := utils.ParseToken(strings.TrimSpace(token))
	if err != nil {
		fmt.Println("Error parsing token:", err) // добавьте эту строку
		return ""
	}
	return role
}
