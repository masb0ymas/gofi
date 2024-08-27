package database

import (
	"fmt"
	"gofi/config"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

var (
	host     = config.Env("DB_HOST", "127.0.01")
	port, _  = strconv.Atoi(config.Env("DB_PORT", "5432"))
	dbname   = config.Env("DB_DATABASE", "db_example")
	username = config.Env("DB_USERNAME", "postgres")
	password = config.Env("DB_PASSWORD", "postgres")
)

func NewDatabase() (*Database, error) {
	// Connection to the database
	var connect string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)

	db, err := sqlx.Open("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
