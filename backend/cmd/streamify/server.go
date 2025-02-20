package main

import (
	"log"
	"sync"

	"github.com/AnishG-git/streamify/internal/connections"
	"github.com/AnishG-git/streamify/internal/handlers"
	"github.com/AnishG-git/streamify/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type server struct {
	router      *mux.Router
	address     string
	rds         storage.Storage
	connections map[string]*websocket.Conn
	mu          *sync.Mutex
	serverID    string
	logger      *log.Logger
}

func newServer(logger *log.Logger, addr string, storage storage.Storage) *server {
	router := mux.NewRouter()

	s := &server{
		router:      router,
		address:     addr,
		rds:         storage,
		connections: make(map[string]*websocket.Conn),
		mu:          &sync.Mutex{},
		serverID:    uuid.NewString(),
		logger:      logger,
	}

	manager := connections.NewManager(s.rds, s.mu, s.connections, s.serverID)
	h := handlers.New(s.logger, manager)

	s.routes(h)
	return s
}

func (s *server) routes(h *handlers.Handlers) {
	room := s.router.PathPrefix("/room").Subrouter()
	room.HandleFunc("/generate", h.GenerateRoomHandler()).Methods("GET")
	room.HandleFunc("/connect/{code}", h.ConnectRoomHandler()).Methods("GET")
}

// func (s *server) connectRoomHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		storage := s.RDS
// 		s.Logger.Print("/connect endpoint called")

// 		// getting room code and name from URL
// 		vars := mux.Vars(r)
// 		roomCode := vars["code"]
// 		name := r.URL.Query().Get("name")

// 		// attempting to upgrade to WebSocket connection
// 		upgrader := websocket.Upgrader{
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true // Allow all origins for simplicity; restrict in production
// 			},
// 		}

// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			s.Logger.Printf("Failed to upgrade to WebSocket: %v", err)
// 			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
// 			return
// 		}
// 		defer conn.Close()

// 		var errMsg string

// 		err = storage.CanUserJoinRoom(ctx, roomCode, name)
// 		if err != nil {
// 			errMsg = fmt.Sprintf("User cannot join room: %v", err)
// 			s.sendErrorMessage(conn, errMsg)
// 			return
// 		}

// 		// checks have passed, adding connection to room
// 		connID := uuid.NewString()

// 		// adding connection to in-memory map
// 		s.mu.Lock()
// 		s.connections[connID] = conn
// 		s.mu.Unlock()

// 		marshalledConnObj, err := json.Marshal(
// 			rdsModels.ConnectionDetails{
// 				ManagerID:    s.serverID,
// 				ConnectionID: connID,
// 			},
// 		)
// 		if err != nil {
// 			errMsg = fmt.Sprintf("Failed to marshal connection object: %v", err)
// 			s.sendErrorMessage(conn, errMsg)
// 			return
// 		}

// 		err = storage.AddUserToRoom(ctx, roomCode, name, string(marshalledConnObj))
// 		if err != nil {
// 			errMsg = fmt.Sprintf("Failed to add connection to room: %v", err)
// 			s.sendErrorMessage(conn, errMsg)
// 			return
// 		}

// 		s.Logger.Printf("User %s has joined room %s", name, roomCode)
// 		ctxWithoutCancel := context.WithoutCancel(ctx)
// 		for {
// 			var message map[string]interface{}
// 			err := conn.ReadJSON(&message)
// 			if err != nil {
// 				// Handle normal WebSocket closure without logging an error
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
// 					s.Logger.Printf("Unexpected WebSocket close error for user %s in room %s: %v", name, roomCode, err)
// 				} else {
// 					s.Logger.Printf("WebSocket closed for user %s in room %s: %v", name, roomCode, err)
// 				}

// 				// Remove connection from room
// 				go s.removeConnectionFromRoom(ctxWithoutCancel, roomCode, name)
// 				break
// 			}

// 			faultyReceiverName, err := s.broadcastToRoom(ctx, roomCode, name, message)
// 			if err != nil {
// 				s.Logger.Printf("Failed to send message to room %s: %v", roomCode, err)
// 				if faultyReceiverName != "" {
// 					go s.removeConnectionFromRoom(ctxWithoutCancel, roomCode, name) // Remove faulty connection
// 				}
// 			} else {
// 				s.Logger.Printf("Message from %s: %v", roomCode, message)
// 			}
// 		}
// 	}
// }
