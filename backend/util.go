package main

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// func (s *server) addUserToRoom(ctx context.Context, roomCode, name, value string) error {
// 	// Using HSet to add a field to the hash (if it doesn't exist) or update it
// 	return s.RDS.HSet(ctx, roomCode, name, value).Err()
// }

// func (s *server) getConnDetails(ctx context.Context, roomCode, name string, unmarshal bool) (*ConnectionDetails, string, error) {
// 	// Using HGet to retrieve a field from the hash
// 	connDetails, err := s.RDS.HGet(ctx, roomCode, name).Result()
// 	if err == redis.Nil {
// 		return nil, "", fmt.Errorf("user not found")
// 	}
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	if !unmarshal {
// 		return nil, connDetails, nil
// 	}

// 	var unmarshalledConnDetails ConnectionDetails
// 	err = json.Unmarshal([]byte(connDetails), &unmarshalledConnDetails)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return &unmarshalledConnDetails, "", nil
// }

// func (s *server) getRoomOccupancy(ctx context.Context, roomCode string) (int, error) {
// 	// Using HLen to get the number of fields in the hash
// 	occupancy, err := s.RDS.HLen(ctx, roomCode).Result()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int(occupancy), nil
// }

// func (s *server) checkForRoom(ctx context.Context, roomCode string) (bool, error) {
// 	// Using Exists to check if the key exists
// 	exists, err := s.RDS.Exists(ctx, roomCode).Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	if exists <= 0 {
// 		return false, nil
// 	}
// 	return true, nil
// }

// func (s *server) removeUserDetails(ctx context.Context, roomCode, name string) error {
// 	// Using HDel to delete a field from the hash
// 	return s.RDS.HDel(ctx, roomCode, name).Err()
// }

// func (s *server) removeRoom(ctx context.Context, roomCode string) error {
// 	// Using Del to delete a key
// 	return s.RDS.Del(ctx, roomCode).Err()
// }

func (s *server) removeConnectionFromRoom(ctx context.Context, roomCode string, name string) {
	storage := s.RDS
	
	// log and return if room code does not exist
	exists, err := storage.CheckForRoom(ctx, roomCode)
	if err != nil {
		s.Logger.Printf("Failed to check existence of room %s: %v", roomCode, err)
		return
	}
	if !exists {
		s.Logger.Printf("Room %s does not exist, connection may have been removed already", roomCode)
		return
	}

	// removing connection details from redis
	connDetails, _, err := storage.GetConnDetails(ctx, roomCode, name, true)
	if err != nil {
		s.Logger.Printf("Failed to get connection details for %s in room %s", name, roomCode)
		return
	}

	err = storage.RemoveUserDetails(ctx, roomCode, name)
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
		s.Logger.Printf("Room %s is empty, waiting for %v ms before removing", sleepTime)
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
	err = storage.RemoveRoom(ctx, roomCode)
	if err != nil {
		s.Logger.Printf("Failed to remove room %s", roomCode)
	}
}

// func (s *server) getUserNamesFromRoom(ctx context.Context, roomCode string) ([]string, error) {
// 	// Using HKeys to get all the fields in the hash
// 	return s.RDS.HKeys(ctx, roomCode).Result()
// }

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
		connDetails, _, err := storage.GetConnDetails(ctx, roomCode, name, true)
		if err != nil {
			return "", err
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
	s.Logger.Printf(message)
	conn.WriteJSON(map[string]string{
		"type":  "error",
		"error": message,
	})
}
