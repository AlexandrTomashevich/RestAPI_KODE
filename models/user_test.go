package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"testing"
)

var testDB *sql.DB

func init() {
	//var err error
	//connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgresSQL"
	//testDB, err = sql.Open("postgres", connStr)
	//if err != nil {
	//	panic(err)
	//}
	var err error // Определение err здесь

	// Используем только одну строку подключения
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgreSQL sslmode=disable connect_timeout=10"
	testDB, err = sql.Open("postgres", connStr)

	if err != nil {
		fmt.Printf("Error when trying to open connection: %v\n", err)
		panic(err)
	}

	err = testDB.Ping()
	if err != nil {
		fmt.Printf("Error when trying to ping database: %v\n", err)
		panic(err)
	}
}

func TestAuthUser(t *testing.T) {
	// Example test
	user, err := AuthUser(testDB, "non-existent-user", "wrong-password")
	if err == nil || user.ID != 0 {
		t.Error("Expected error for non-existent user, got none")
	}
}

func TestCreateUser(t *testing.T) {
	// Example test
	userID, err := CreateUser(testDB, "test-user", "test-password", "test-role")
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}
	if userID == 0 {
		t.Error("Expected non-zero user ID, got zero")
	}
}
