package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	mainLog := log.Default()
	rds, err := mustLoadRedis()
	if err != nil {
		mainLog.Fatalf("Failed to load database: %s", err)
	}
	mainLog.Print("Connected to database")

	const addr = ":8080"
	server := newServer(mainLog, addr, rds)
	errCh := make(chan error)

	go func() {
		mainLog.Print("Starting up server on port 8080")
		// Allow CORS for local development
		cors := handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		)
		if err := http.ListenAndServe(server.address, cors(server.router)); err != nil {
			mainLog.Printf("Server failed: %s", err)
			errCh <- err
		}
	}()

	<-errCh

}
