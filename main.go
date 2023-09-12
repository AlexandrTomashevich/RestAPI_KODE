package main

import (
	"RestAPI_KODE/api"
	"RestAPI_KODE/config"
	"RestAPI_KODE/database"
	"RestAPI_KODE/lib"
	"RestAPI_KODE/middleware"
	"RestAPI_KODE/models"
	"database/sql"
	"log"
	"net/http"
)

func main() {
	log.Println("Started app")
	var err error
	var db *sql.DB

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		lib.Fatalf("Failed to load config: %s", err)
	}

	// Инициализация подключения к базе данных с использованием строки подключения из конфигурационного файла
	db, err = database.InitializeDB(cfg.Database.ConnectionString)
	if err != nil {
		lib.Fatalf("Failed to initialized db: %s", err)
	}

	models.SetUserDB(db)
	//api.SetDatabase(db)

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			api.AuthUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			api.CreateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.AuthorizationMiddleware(http.HandlerFunc(api.AddNote)).ServeHTTP(w, r)
		} else if r.Method == "GET" {
			middleware.AuthorizationMiddleware(http.HandlerFunc(api.GetNotes)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		lib.Errorf("Error starting the server: %s", err)
	}
}
