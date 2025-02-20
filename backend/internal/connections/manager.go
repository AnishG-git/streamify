package connections

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AnishG-git/streamify/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	rdsModels "github.com/AnishG-git/streamify/internal/storage/models"
)

type ConnManager interface {
	RemoveConnectionFromRoom(ctx context.Context, logger *log.Logger, roomCode string, name string)
	BroadcastToRoom(ctx context.Context, logger *log.Logger, roomCode string, senderName string, message map[string]interface{}) (string, error)
	SetConnection(conn *websocket.Conn) *rdsModels.ConnectionDetails
	storage.Storage
}

type Manager struct {
	rds         storage.Storage
	mu          *sync.Mutex
	connections map[string]*websocket.Conn
	managerID   string
}

func NewManager(rds storage.Storage, mu *sync.Mutex, conns map[string]*websocket.Conn, managerID string) *Manager {
	return &Manager{
		rds:         rds,
		mu:          mu,
		connections: conns,
		managerID:   managerID,
	}
}

func (c *Manager) RemoveConnectionFromRoom(ctx context.Context, logger *log.Logger, roomCode string, name string) {
	storage := c.rds

	// removing connection details from redis
	connDetailsStr, err := storage.GetUserConnectionDetails(ctx, roomCode, name)
	if err != nil {
		logger.Printf("Failed to get connection details for %s in room %s", name, roomCode)
		return
	}

	var connDetails rdsModels.ConnectionDetails
	err = json.Unmarshal([]byte(connDetailsStr), &connDetails)
	if err != nil {
		logger.Printf("Failed to unmarshal connection details for %s in room %s", name, roomCode)
		return
	}

	err = storage.RemoveUserFromRoom(ctx, roomCode, name)
	if err != nil {
		logger.Printf("Failed to remove user %s from room %s", name, roomCode)
		return
	}

	// Delete the closed connection from in-memory map
	c.mu.Lock()
	delete(c.connections, connDetails.ConnectionID)
	c.mu.Unlock()

	// return if the room is not empty
	roomOccupancy, err := storage.GetRoomOccupancy(ctx, roomCode)
	if err != nil {
		logger.Printf("Failed to get room occupancy for room %s", roomCode)
		return
	}
	if roomOccupancy > 0 {
		return
	}

	const maxAttempts = 10
	const sleepTime = 500
	for i := 0; i < maxAttempts; i++ {
		roomOccupancy, err = storage.GetRoomOccupancy(ctx, roomCode)
		if err != nil {
			logger.Printf("Failed to get room occupancy for room %s in sleep", roomCode)
		}
		if roomOccupancy > 0 {
			return
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	// If the room is still empty, delete it
	err = storage.DeleteRoom(ctx, roomCode)
	if err != nil {
		logger.Printf("Failed to remove room %s", roomCode)
	}
}

func (m *Manager) BroadcastToRoom(ctx context.Context, logger *log.Logger, roomCode string, senderName string, message map[string]interface{}) (string, error) {
	storage := m.rds

	names, err := storage.GetUserNamesFromRoom(ctx, roomCode)
	if err != nil {
		logger.Printf("Failed to get user names from room %s", roomCode)
		return "", err
	}

	for _, name := range names {
		connDetailsStr, err := storage.GetUserConnectionDetails(ctx, roomCode, name)
		if err != nil {
			return "", err
		}

		var connDetails rdsModels.ConnectionDetails
		err = json.Unmarshal([]byte(connDetailsStr), &connDetails)
		if err != nil {
			logger.Printf("Failed to unmarshal connection details for %s in room %s", name, roomCode)
		}

		conn, ok := m.connections[connDetails.ConnectionID]
		if !ok {
			return "", fmt.Errorf("connection not found for user %s in room %s", name, roomCode)
		}
		if name != senderName {
			err := conn.WriteJSON(message)
			if err != nil {
				return name, err
			}
		}
	}
	return "", nil
}

func (m *Manager) SetConnection(conn *websocket.Conn) *rdsModels.ConnectionDetails {
	// checks have passed, adding connection to room
	connID := uuid.NewString()

	// adding connection to in-memory map
	m.mu.Lock()
	m.connections[connID] = conn
	m.mu.Unlock()

	return &rdsModels.ConnectionDetails{
		ManagerID:    m.managerID,
		ConnectionID: connID,
	}
}

func (m *Manager) CreateRoom(ctx context.Context, roomCode string) error {
	return m.rds.CreateRoom(ctx, roomCode)
}

func (m *Manager) DeleteRoom(ctx context.Context, roomCode string) error {
	return m.rds.DeleteRoom(ctx, roomCode)
}

func (m *Manager) IsRoomActive(ctx context.Context, roomCode string) (bool, error) {
	return m.rds.IsRoomActive(ctx, roomCode)
}

func (m *Manager) GetRoomOccupancy(ctx context.Context, roomCode string) (int, error) {
	return m.rds.GetRoomOccupancy(ctx, roomCode)
}

func (m *Manager) AddUserToRoom(ctx context.Context, roomCode, username, connID string) error {
	return m.rds.AddUserToRoom(ctx, roomCode, username, connID)
}

func (m *Manager) RemoveUserFromRoom(ctx context.Context, roomCode, username string) error {
	return m.rds.RemoveUserFromRoom(ctx, roomCode, username)
}

func (m *Manager) GetUserNamesFromRoom(ctx context.Context, roomCode string) ([]string, error) {
	return m.rds.GetUserNamesFromRoom(ctx, roomCode)
}

func (m *Manager) GetUserConnectionDetails(ctx context.Context, roomCode, username string) (string, error) {
	return m.rds.GetUserConnectionDetails(ctx, roomCode, username)
}

func (m *Manager) CanUserJoinRoom(ctx context.Context, roomCode, name string) error {
	return m.rds.CanUserJoinRoom(ctx, roomCode, name)
}
