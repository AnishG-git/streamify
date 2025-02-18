package storage

import (
	"context"
	"github.com/AnishG-git/streamify/models"
)

type Storage interface {
	AddUserToRoom(ctx context.Context, roomCode, name, value string) error
	GetConnDetails(ctx context.Context, roomCode, name string, unmarshal bool) (*models.ConnectionDetails, string, error)
	GetRoomOccupancy(ctx context.Context, roomCode string) (int, error)
	CheckForUser(ctx context.Context, roomCode, name string) (bool, error)
	CheckForRoom(ctx context.Context, roomCode string) (bool, error)
	RemoveUserDetails(ctx context.Context, roomCode, name string) error
	RemoveRoom(ctx context.Context, roomCode string) error
	GetUserNamesFromRoom(ctx context.Context, roomCode string) ([]string, error)
}
