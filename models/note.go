package models

import (
	"RestAPI_KODE/database"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func CreateNote(userID int, content string) (int, error) {
	var noteID int
	err := database.DB.QueryRow(`INSERT INTO notes (userId, content) VALUES ($1, $2) RETURNING id`, userID, content).Scan(&noteID)
	if err != nil {
		return 0, err
	}
	return noteID, nil
}
