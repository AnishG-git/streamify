package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/AnishG-git/streamify/models"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)

func generateRoomCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 5
	var sb strings.Builder
	sb.Grow(length)
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func (s *server) removeConnectionFromRoom(ctx context.Context, roomCode string, name string) {
	storage := s.RDS

	// removing connection details from redis
	connDetailsStr, err := storage.GetUserConnectionDetails(ctx, roomCode, name)
	if err != nil {
		s.Logger.Printf("Failed to get connection details for %s in room %s", name, roomCode)
		return
	}

	var connDetails models.ConnectionDetails
	err = json.Unmarshal([]byte(connDetailsStr), &connDetails)
	if err != nil {
		s.Logger.Printf("Failed to unmarshal connection details for %s in room %s", name, roomCode)
		return
	}

	err = storage.RemoveUserFromRoom(ctx, roomCode, name)
	if err != nil {
		s.Logger.Printf("Failed to remove user %s from room %s", name, roomCode)
		return
	}

	// Delete the closed connection from in-memory map
	s.mu.Lock()
	delete(s.connections, connDetails.ConnectionID)
	s.mu.Unlock()

	// return if the room is not empty
	roomOccupancy, err := storage.GetRoomOccupancy(ctx, roomCode)
	if err != nil {
		s.Logger.Printf("Failed to get room occupancy for room %s", roomCode)
		return
	}
	if roomOccupancy > 0 {
		return
	}

	const maxAttempts = 10
	const sleepTime = 500
	for i := 0; i < maxAttempts; i++ {
		s.Logger.Printf("Room %s is empty, waiting for %v ms before removing", roomCode, sleepTime)
		roomOccupancy, err = storage.GetRoomOccupancy(ctx, roomCode)
		if err != nil {
			s.Logger.Printf("Failed to get room occupancy for room %s in sleep", roomCode)
		}
		if roomOccupancy > 0 {
			return
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	// If the room is still empty, delete it
	err = storage.DeleteRoom(ctx, roomCode)
	if err != nil {
		s.Logger.Printf("Failed to remove room %s", roomCode)
	}
}

// the returned connection is faulty if and only if the returned error is not nil
func (s *server) broadcastToRoom(ctx context.Context, roomCode string, sender string, message map[string]interface{}) (string, error) {
	// Iterate over all connections in the room
	storage := s.RDS
	names, err := storage.GetUserNamesFromRoom(ctx, roomCode)
	if err != nil {
		s.Logger.Printf("Failed to get user names from room %s", roomCode)
		return "", err
	}

	for _, name := range names {
		connDetailsStr, err := storage.GetUserConnectionDetails(ctx, roomCode, name)
		if err != nil {
			return "", err
		}

		var connDetails models.ConnectionDetails
		err = json.Unmarshal([]byte(connDetailsStr), &connDetails)
		if err != nil {
			s.Logger.Printf("Failed to unmarshal connection details for %s in room %s", name, roomCode)
		}

		conn, ok := s.connections[connDetails.ConnectionID]
		if !ok {
			return "", fmt.Errorf("connection not found for user %s in room %s", name, roomCode)
		}
		if name != sender {
			err := conn.WriteJSON(message)
			if err != nil {
				return name, err
			}
		}
	}
	return "", nil
}

func (s *server) sendErrorMessage(conn *websocket.Conn, message string) {
	s.Logger.Print(message)
	conn.WriteJSON(map[string]string{
		"type":  "error",
		"error": message,
	})
}
