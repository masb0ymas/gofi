package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"gofi/internal/config"
)

func connectRedis(cfg *config.ConfigRedis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
