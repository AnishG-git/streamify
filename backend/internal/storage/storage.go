package storage

import (
	"context"
)

type Storage interface {
	// Room Management
	CreateRoom(ctx context.Context, roomCode string) error
	DeleteRoom(ctx context.Context, roomCode string) error
	IsRoomActive(ctx context.Context, roomCode string) (bool, error)
	GetRoomOccupancy(ctx context.Context, roomCode string) (int, error)

	// User Management
	AddUserToRoom(ctx context.Context, roomCode, username, connID string) error
	RemoveUserFromRoom(ctx context.Context, roomCode, username string) error
	GetUserNamesFromRoom(ctx context.Context, roomCode string) ([]string, error)
	GetUserConnectionDetails(ctx context.Context, roomCode, username string) (string, error)

	// User can join room if and only if the returned error is nil
	CanUserJoinRoom(ctx context.Context, roomCode, name string) error
}
