package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Env(key string, fallback string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
