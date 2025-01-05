package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	mainLog := log.Default()
	db, err := loadDB("streamify", false)
	if err != nil {
		mainLog.Fatalf("Failed to load database: %s", err)
	}
	mainLog.Print("Connected to database")

	server := newServer(":8080", db, mainLog)
	errCh := make(chan error)

	go func() {
		mainLog.Print("Starting up server on port 8080")
		// Allow CORS for local development
		cors := handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		)
		if err := http.ListenAndServe(server.Address, cors(server.Router)); err != nil {
			mainLog.Printf("Server failed: %s", err)
			errCh <- err
		}
	}()

	<-errCh

}
