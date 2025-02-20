package storage

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

type RDS struct {
	cli            *redis.Client
	activeRoomsKey string
}

func NewRDS(cli *redis.Client) *RDS {
	return &RDS{
		cli:            cli,
		activeRoomsKey: "active-rooms",
	}
}

func (r *RDS) CreateRoom(ctx context.Context, roomCode string) error {
	return r.cli.SAdd(ctx, r.activeRoomsKey, roomCode).Err()
}

func (r *RDS) DeleteRoom(ctx context.Context, roomCode string) error {
	return r.cli.SRem(ctx, r.activeRoomsKey, roomCode).Err()
}

func (r *RDS) IsRoomActive(ctx context.Context, roomCode string) (bool, error) {
	return r.cli.SIsMember(ctx, r.activeRoomsKey, roomCode).Result()
}

func (r *RDS) CanUserJoinRoom(ctx context.Context, roomCode string, name string) error {
	roomIsActive, err := r.IsRoomActive(ctx, roomCode)
	if err != nil {
		return err
	}
	if !roomIsActive {
		return fmt.Errorf("room %s does not exist in active set", roomCode)
	}
	roomOccupancy, err := r.GetRoomOccupancy(ctx, roomCode)
	if err != nil {
		return err
	}
	if roomOccupancy < 2 {
		// check if user already exists in room
		userExists, err := r.cli.HExists(ctx, roomCode, name).Result()
		if err != nil {
			return err
		}
		if userExists {
			return fmt.Errorf("user %s already exists in room %s", name, roomCode)
		}
		return nil
	}
	return fmt.Errorf("room %s is at capacity", roomCode)
}

func (r *RDS) GetRoomOccupancy(ctx context.Context, roomCode string) (int, error) {
	// Using HLen to get the number of fields in the hash
	occupancy, err := r.cli.HLen(ctx, roomCode).Result()
	if err != nil {
		return 0, err
	}
	return int(occupancy), nil
}

func (r *RDS) AddUserToRoom(ctx context.Context, roomCode, name, value string) error {
	return r.cli.HSet(ctx, roomCode, name, value).Err()
}

func (r *RDS) RemoveUserFromRoom(ctx context.Context, roomCode, name string) error {
	return r.cli.HDel(ctx, roomCode, name).Err()
}

func (r *RDS) GetUserConnectionDetails(ctx context.Context, roomCode, name string) (string, error) {
	// Using HGet to retrieve a field from the hash
	connDetails, err := r.cli.HGet(ctx, roomCode, name).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("user %s not found in room %s", name, roomCode)
	}
	if err != nil {
		return "", err
	}
	return connDetails, nil
}

func (r *RDS) GetUserNamesFromRoom(ctx context.Context, roomCode string) ([]string, error) {
	// Using HKeys to get all the fields in the hash
	return r.cli.HKeys(ctx, roomCode).Result()
}
