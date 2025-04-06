package database

import (
	"fmt"
	"gofi/config"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/masb0ymas/go-utils/pkg"
)

var DB *sqlx.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// NewDBConfig creates a new database configuration from environment variables
func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     config.Env("DB_HOST", "localhost"),
		Port:     config.Env("DB_PORT", "5432"),
		User:     config.Env("DB_USERNAME", "postgres"),
		Password: config.Env("DB_PASSWORD", "postgres"),
		Name:     config.Env("DB_DATABASE", "example"),
		SSLMode:  config.Env("DB_SSL_MODE", "disable"),
	}
}

// Connect establishes a connection to the database
func Connect() error {
	config := NewDBConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging the database: %v", err)
	}

	DB = db

	msg := pkg.Println("Sqlx", fmt.Sprintf("Successfully connected to database: %s", config.Name))
	log.Println(msg)

	return nil
}

func GetDB() *sqlx.DB {
	return DB
}
