package main

import (
	"database/sql"
	"strings"

	"golang.org/x/exp/rand"
)

func generateRoomCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 5
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// checks if a room code is valid
func roomCodeExists(db *sql.DB, code string) bool {
	if len(code) != 5 {
		return false
	}
	const query = `SELECT EXISTS(SELECT 1 FROM room WHERE code = $1)`
	var exists bool
	err := db.QueryRow(query, code).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// creates a new room in the database
func createRoom(db *sql.DB, code string) error {
	const query = `INSERT INTO room (code, participants) VALUES ($1, $2)`
	_, err := db.Exec(query, code, 1)
	return err
}