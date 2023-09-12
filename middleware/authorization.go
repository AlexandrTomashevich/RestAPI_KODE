package middleware

import (
	"RestAPI_KODE/lib"
	"RestAPI_KODE/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
)

var signingKey = "35FSFJlgh4353KSFjX"

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		token := r.Header.Get("Authorization")
		if token == "" {
			lib.Errorf("No token provided")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userRole := CheckToken(token)

		if userRole == "admin" {
			next.ServeHTTP(w, r)
			return
		} else if userRole == "" || userRole == "unknown" {
			lib.Errorf("Unknown or empty user role: %v", userRole)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		noteID := r.URL.Query().Get("note_id")
		if noteID == "" {
			lib.Errorf("No noteID provided in the request")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !userHasAccessToNote(token, noteID) {
			lib.Errorf("User does not have access to note: %v", noteID)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Проверка прав пользователя на доступ к заметке
func userHasAccessToNote(token, noteID string) bool {
	userID, err := getUserIDFromToken(token)
	if err != nil {
		lib.Errorf("Failed to get user ID from token: %v", err)
		return false
	}

	noteUserID, err := getUserIDFromNote(noteID)
	if err != nil {
		lib.Errorf("Failed to get user ID from note: %v", err)
		return false
	}

	return userID == noteUserID
}

func getUserIDFromToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		lib.Errorf("Error parsing token: %v", err)
		return 0, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.UserId, nil
	}
	return 0, fmt.Errorf("invalid token claims")
}

func getUserIDFromNote(noteID string) (int, error) {
	id, err := strconv.Atoi(noteID)
	if err != nil {
		lib.Errorf("Error converting noteID to integer: %v", err)
		return 0, err
	}

	note, err := models.GetNoteByID(id)
	if err != nil {
		lib.Errorf("Error getting note from DB: %v", err)
		return 0, err
	}

	return note.UserID, nil
}
