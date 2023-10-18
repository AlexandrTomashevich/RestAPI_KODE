package database

import (
	"RestAPI_KODE/config"
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func _(db *sql.DB) error {

	err := createTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

func NewConnection(dbConfig config.Database) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("postgres", dbConfig.ConnectionString())
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func createTablesIfNotExists() error {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	)`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		userId INTEGER REFERENCES users(id),
		content TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	return err
}
