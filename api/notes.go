package api

import (
	"RestAPI_KODE/database"
	"RestAPI_KODE/lib"
	"RestAPI_KODE/middleware"
	"RestAPI_KODE/models"
	"RestAPI_KODE/utils"
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
	"time"
)

func AddNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, _, err := utils.CheckToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if _, err := models.UserExists(userID); err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	correctedContent, err := utils.CheckSpelling(note.Content)
	if err != nil {
		http.Error(w, "Failed to check the content", http.StatusInternalServerError)
		return
	}
	note.Content = correctedContent

	noteID, err := models.CreateNote(userID, note.Content)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	note.ID = noteID
	note.UserId = userID
	note.Timestamp = time.Now()

	response := map[string]interface{}{
		"id":        note.ID,
		"userID":    note.UserId,
		"content":   note.Content,
		"timestamp": note.Timestamp,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
	}
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, role, err := utils.CheckToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	notes, err := retrieveNotes(userID, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
}

func retrieveNotes(userID int, role string) ([]models.Note, error) {
	var rows *sql.Rows
	var err error
	var query string
	if role == "admin" {
		query = "SELECT id, content FROM notes"
	} else if role == "user" {
		query = "SELECT id, content FROM notes WHERE userId = $1"
	} else {
		return nil, errors.New("Unauthorized")
	}

	rows, err = database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []models.Note{}
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Content); err != nil {
			lib.Errorf("Error scanning rows: %v", err)
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	authenticatedUser, err := middleware.AuthenticateUser(user.Username, user.PasswordHash)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateTokenForUser(*authenticatedUser)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":       authenticatedUser.ID,
		"username": authenticatedUser.Username,
		"role":     authenticatedUser.Role,
		"token":    token,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID, err := middleware.CreateUser(user.Username, user.PasswordHash, user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			status := http.StatusConflict
			message := "User already exists."
			if user.Role != "admin" {
				message = "User already exists. Please authenticate to access your notes."
			}
			http.Error(w, message, status)
			return
		}
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}

	token, err := middleware.GenerateTokenForUser(user)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id":    userID,
		"token": token,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
	}
}
