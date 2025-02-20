package main

import (
	"context"
	"time"

	"github.com/AnishG-git/streamify/internal/storage"
	"github.com/redis/go-redis/v9"
)

func mustLoadRedis() (storage.Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	rds := storage.NewRDS(client)
	return rds, nil
}
