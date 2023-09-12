package models

import (
	"database/sql"
	"errors"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var db *sql.DB

func SetUserDB(database *sql.DB) {
	db = database
}

func GetNoteByID(noteID int) (*Note, error) {
	if db == nil {
		return nil, errors.New("Database connection not initialized")
	}

	var note Note
	query := "SELECT id, userId, content, timestamp FROM notes WHERE id = $1"
	err := db.QueryRow(query, noteID).Scan(&note.ID, &note.UserID, &note.Content, &note.Timestamp)
	if err != nil {
		return nil, err
	}
	return &note, nil
}
