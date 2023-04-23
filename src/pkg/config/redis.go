package config

import (
	"fmt"
	"gofi/src/pkg/helpers"

	"github.com/redis/go-redis/v9"
)

/*
Initialize Redis
*/
func NewRedisClient() *redis.Client {
	redisHost := Env("REDIS_HOST", "127.0.0.1")
	redisPort := Env("REDIS_PORT", "6379")
	redisPass := Env("REDIS_PASSWORD", "")

	// format string
	redisURL := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       0,
	})

	logMessage := helpers.PrintLog("Redis", "initialized successfully")
	fmt.Println(logMessage)

	return client
}
