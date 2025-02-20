package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnishG-git/streamify/internal/connections"
	"github.com/AnishG-git/streamify/internal/logic"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	// this should be set on/retrieved from the context
	logger *log.Logger
	// these should belong to a connection manager
	manager connections.ConnManager
}

func New(logger *log.Logger, manager connections.ConnManager) *Handlers {
	return &Handlers{
		logger:  logger,
		manager: manager,
	}
}

func (h *Handlers) GenerateRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h.logger.Print("/room/generate endpoint called")
		w.Header().Set("Content-Type", "application/json")

		roomCode, err := logic.GenerateRoomLogic(ctx, h.logger, h.manager)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"code": roomCode})
	}
}

func (h *Handlers) ConnectRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h.logger.Print("/connect endpoint called")

		// getting room code and name from URL
		vars := mux.Vars(r)
		roomCode := vars["code"]
		name := r.URL.Query().Get("name")

		// attempting to upgrade to WebSocket connection
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity; restrict in production
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			h.logger.Printf("Failed to upgrade to WebSocket: %v", err)
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		defer conn.Close()

		// Executing the logic to connect to the room
		errMsg, err := logic.ConnectToRoomLogic(ctx, h.logger, h.manager, roomCode, name, conn)
		if err != nil {
			h.logger.Printf("Failed at ConnectToRoomLogic: %v", err)
			conn.WriteJSON(map[string]string{
				"type":  "error",
				"error": errMsg,
			})
		}
	}
}
