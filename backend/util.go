package main

import (
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

// checks if a room code is valid
// func roomCodeExists(db *sql.DB, code string) bool {
// 	if len(code) != 5 {
// 		return false
// 	}
// 	const query = `SELECT EXISTS(SELECT 1 FROM room WHERE code = $1)`
// 	var exists bool
// 	err := db.QueryRow(query, code).Scan(&exists)
// 	if err != nil {
// 		return false
// 	}
// 	return exists
// }

// creates a new room in the database
// func createRoom(db *sql.DB, code string) error {
// 	const query = `INSERT INTO room (code, participants) VALUES ($1, $2)`
// 	_, err := db.Exec(query, code, 1)
// 	return err
// }

func (s *server) removeConnectionFromRoom(roomCode string, conn *websocket.Conn) {
	s.mu.RLock()

	conns, exists := s.Rooms[roomCode]
	if !exists {
		s.mu.RUnlock()
		return
	}

	// Remove the closed connection
	for i, c := range conns {
		if c == conn {
			s.mu.RUnlock()
			s.mu.Lock()
			s.Rooms[roomCode] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	// If the room is not empty, return
	if len(s.Rooms[roomCode]) > 0 {
		s.mu.Unlock()
		return
	}

	s.mu.Unlock()

	for i := 0; i < 50; i++ {
		s.mu.RLock()
		if len(s.Rooms[roomCode]) > 0 {
			s.mu.RUnlock()
			return
		}
		s.mu.RUnlock()
		time.Sleep(100 * time.Millisecond)
	}

	// If the room is still empty, delete it
	s.mu.Lock()
	delete(s.Rooms, roomCode)
	s.mu.Unlock()
}
