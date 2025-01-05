package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	Router  *mux.Router
	Address string
	DB      *sql.DB
	Logger  *log.Logger
}

func newServer(addr string, db *sql.DB, logger *log.Logger) *server {
	router := mux.NewRouter()

	s := &server{
		Router:  router,
		Address: addr,
		DB:      db,
		Logger:  logger,
	}

	s.routes()
	return s
}

func (s *server) routes() {
	room := s.Router.PathPrefix("/room").Subrouter()
	room.HandleFunc("/generate", s.generateRoomHandler()).Methods("GET")
}

func (s *server) generateRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Print("/room/generate endpoint called")
		code := generateRoomCode()
		for roomCodeExists(s.DB, code) {
			code = generateRoomCode()
		}
		w.Header().Set("Content-Type", "application/json")
		if err := createRoom(s.DB, code); err != nil {
			s.Logger.Printf("Failed to create room: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"code": "1"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"code": code})
	}
}
