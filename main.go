package main

import (
	"RestAPI_KODE/api"
	"RestAPI_KODE/config"
	"RestAPI_KODE/database"
	"RestAPI_KODE/lib"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	log.Println("Started app")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		lib.Fatalf("Failed to load config: %s", err)
	}

	database.DB, err = database.NewConnection(cfg.Database)
	if err != nil {
		lib.Fatalf("Failed to initialized db: %s", err)
	}

	if database.DB == nil {
		log.Fatal("Database not initialized!")
	}

	lib.Infof("Migrations ok")
	http.HandleFunc("/auth", authenticateHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/notes", notesHandler)

	server = &http.Server{Addr: ":8080", Handler: nil}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			lib.Fatalf("Error starting the server: %s", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals

	if err := server.Shutdown(nil); err != nil {
		lib.Fatalf("Error shutting down the server: %s", err)
	}
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		api.AuthUser(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		api.CreateUser(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		api.AddNote(w, r)
	} else if r.Method == "GET" {
		api.GetNotes(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
