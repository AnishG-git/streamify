package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type server struct {
	Router  *mux.Router
	Address string
	DB      *sql.DB
	Logger  *log.Logger
	Rooms   map[string][]*websocket.Conn
	mu      sync.RWMutex
}

func newServer(addr string, db *sql.DB, logger *log.Logger) *server {
	router := mux.NewRouter()

	s := &server{
		Router:  router,
		Address: addr,
		DB:      db,
		Logger:  logger,
		mu:      sync.RWMutex{},
		Rooms:   make(map[string][]*websocket.Conn),
	}

	s.routes()
	return s
}

func (s *server) routes() {
	room := s.Router.PathPrefix("/room").Subrouter()
	room.HandleFunc("/generate", s.generateRoomHandler()).Methods("GET")
	room.HandleFunc("/connect/{code}", s.connectRoomHandler()).Methods("GET")
}

func (s *server) generateRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Print("/room/generate endpoint called")
		w.Header().Set("Content-Type", "application/json")

		code := generateRoomCode()

		// check if the room code already exists
		s.mu.Lock()
		for _, exists := s.Rooms[code]; exists; _, exists = s.Rooms[code] {
			code = generateRoomCode()
		}
		s.Rooms[code] = make([]*websocket.Conn, 0, 2)
		s.mu.Unlock()

		// DEBUGGING (DELETE LATER)
		s.Logger.Printf("Room %s created", code)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"code": code})
	}
}

func (s *server) connectRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomCode := vars["code"]

		// Attempt to upgrade to WebSocket
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity; restrict in production
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.Logger.Printf("Failed to upgrade to WebSocket: %v", err)
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		s.mu.Lock()
		if _, exists := s.Rooms[roomCode]; !exists {
			s.mu.Unlock()
			conn.WriteJSON(map[string]string{
				"type":  "error",
				"error": "invalid room code",
			})
			return
		}

		if len(s.Rooms[roomCode]) == 2 {
			s.mu.Unlock()
			conn.WriteJSON(map[string]string{
				"type":  "error",
				"error": "Room is full",
			})
			return
		}

		s.Rooms[roomCode] = append(s.Rooms[roomCode], conn)
		s.mu.Unlock()
		s.Logger.Printf("New connection to room: %s", roomCode)

		for {
			var message map[string]interface{}
			err := conn.ReadJSON(&message)
			if err != nil {
				// Handle normal WebSocket closure without logging an error
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					s.Logger.Printf("Unexpected WebSocket close error: %v", err)
				} else {
					s.Logger.Printf("WebSocket closed for room %s: %v", roomCode, err)
				}

				// Remove connection from room
				go s.removeConnectionFromRoom(roomCode, conn)
				break
			}

			s.Logger.Printf("Message from %s: %v", roomCode, message)
			// Broadcast message to other connections in the room (implementation needed)
		}
	}
}
