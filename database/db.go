package database

import (
	"RestAPI_KODE/lib"
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitializeDB(connectionString string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		lib.Logger.Fatalf("Failed to initialize the database: %v", err)
		return nil, err
	}
	return db, nil
}

func NewConnection() (*sql.DB, error) {
	connStr := "host=db dbname=notesapp user=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Создание таблиц, если их нет
	err = createTablesIfNotExists(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTablesIfNotExists(db *sql.DB) error {
	// Таблица пользователей
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	)`)
	if err != nil {
		return err
	}

	// Таблица заметок
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		userId INTEGER REFERENCES users(id),
		content TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}
