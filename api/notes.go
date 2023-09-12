package api

import (
	"RestAPI_KODE/lib"
	models2 "RestAPI_KODE/models"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"net/http"
)

var db *sql.DB

func AddNote(w http.ResponseWriter, r *http.Request) {
	var note models2.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка через Яндекс.Спеллер здесь...
	err = db.QueryRow("INSERT INTO notes(content) VALUES($1) RETURNING id", note.Content).Scan(&note.ID)
	if err != nil {
		lib.Errorf("Error message here: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		lib.Errorf("Error encoding the note: %v", err)
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, content FROM notes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			lib.Errorf("Error closing rows: %v", cerr)
			err = cerr
		}
	}()

	var notes []models2.Note
	for rows.Next() {
		var note models2.Note
		if err := rows.Scan(&note.ID, &note.Content); err != nil {
			lib.Errorf("Error message here: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	json.NewEncoder(w).Encode(notes)
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		lib.Errorf("Error encoding the notes: %v", err)
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	var user models2.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	authenticatedUser, err := models2.AuthUser(db, user.Username, user.PasswordHash)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Здесь для упрощения возвращается объект пользователя.
	json.NewEncoder(w).Encode(authenticatedUser)
	if err := json.NewEncoder(w).Encode(authenticatedUser); err != nil {
		lib.Errorf("Error encoding the authenticatedUser: %v", err)
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models2.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID, err := models2.CreateUser(db, user.Username, user.PasswordHash, user.Role)
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}

	resp := map[string]int{
		"id": userID,
	}
	json.NewEncoder(w).Encode(resp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		lib.Errorf("Error encoding the resp: %v", err)
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}
}
