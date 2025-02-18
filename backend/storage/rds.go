package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnishG-git/streamify/models"
	redis "github.com/redis/go-redis/v9"
)

type RDS struct {
	cli *redis.Client
}

func NewRDS(cli *redis.Client) *RDS {
	return &RDS{cli: cli}
}

func (r *RDS) AddUserToRoom(ctx context.Context, roomCode, name, value string) error {
	// Using HSet to add a field to the hash (if it doesn't exist) or update it
	return r.cli.HSet(ctx, roomCode, name, value).Err()
}

func (r *RDS) GetConnDetails(ctx context.Context, roomCode, name string, unmarshal bool) (*models.ConnectionDetails, string, error) {
	// Using HGet to retrieve a field from the hash
	connDetails, err := r.cli.HGet(ctx, roomCode, name).Result()
	if err == redis.Nil {
		return nil, "", fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, "", err
	}

	if !unmarshal {
		return nil, connDetails, nil
	}

	var unmarshalledConnDetails models.ConnectionDetails
	err = json.Unmarshal([]byte(connDetails), &unmarshalledConnDetails)
	if err != nil {
		return nil, "", err
	}

	return &unmarshalledConnDetails, "", nil
}

func (r *RDS) GetRoomOccupancy(ctx context.Context, roomCode string) (int, error) {
	// Using HLen to get the number of fields in the hash
	occupancy, err := r.cli.HLen(ctx, roomCode).Result()
	if err != nil {
		return 0, err
	}
	return int(occupancy), nil
}

func (r *RDS) CheckForUser (ctx context.Context, roomCode, name string) (bool, error) {
	// Using HExists to check if a field exists in the hash
	exists, err := r.cli.HExists(ctx, roomCode, name).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *RDS) CheckForRoom(ctx context.Context, roomCode string) (bool, error) {
	// Using Exists to check if the key exists
	exists, err := r.cli.Exists(ctx, roomCode).Result()
	if err != nil {
		return false, err
	}
	if exists <= 0 {
		return false, nil
	}
	return true, nil
}

func (r *RDS) RemoveUserDetails(ctx context.Context, roomCode, name string) error {
	// Using HDel to delete a field from the hash
	return r.cli.HDel(ctx, roomCode, name).Err()
}

func (r *RDS) RemoveRoom(ctx context.Context, roomCode string) error {
	// Using Del to delete a key
	return r.cli.Del(ctx, roomCode).Err()
}

func (r *RDS) GetUserNamesFromRoom(ctx context.Context, roomCode string) ([]string, error) {
	// Using HKeys to get all the fields in the hash
	return r.cli.HKeys(ctx, roomCode).Result()
}