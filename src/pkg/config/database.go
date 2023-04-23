package config

import (
	"database/sql"
	"fmt"
	"gofi/src/database/entities"
	"gofi/src/database/migrations"
	"gofi/src/database/seeds"
	"gofi/src/pkg/helpers"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// Connect to Database
func ConnectDB() {
	DB_CONNECTION := Env("DB_CONNECTION", "postgres")
	DB_HOST := Env("DB_HOST", "127.0.0.1")
	DB_PORT := Env("DB_PORT", "5432")
	DB_DATABASE := Env("DB_DATABASE", "dev_dbgofi")
	DB_USERNAME := Env("DB_USERNAME", "postgres")
	DB_PASSWORD := Env("DB_PASSWORD", "postgres")
	DB_SSL := Env("DB_SSL", "disable")

	var err error
	port, err := strconv.Atoi(DB_PORT)

	if err != nil {
		logErrMessage := helpers.PrintLog("Sqlx", "Error database port", helpers.Options{Label: "error"})
		fmt.Println(logErrMessage)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", DB_HOST, port, DB_USERNAME, DB_PASSWORD, DB_DATABASE, DB_SSL)

	// Open Connection
	db, err = sqlx.Open(DB_CONNECTION, dsn)

	db.SetMaxOpenConns(1000) // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(0) // 0, connections are reused forever.

	if err != nil {
		logErrMessage := helpers.PrintLog("Sqlx", "Failed to connect database", helpers.Options{Label: "error"})
		fmt.Println(logErrMessage)
	}

	// Run Migrate & Seeder
	migrate()
	seeder()

	// Verify Connection
	logMessage := helpers.PrintLog("Sqlx", "Connection "+DB_DATABASE+" has been established successfully.")
	fmt.Println(logMessage)
}

// Declare Global Database Instance
func GetDB() *sqlx.DB {
	return db
}

// Run Migration
func migrate() {
	collectSchema := []string{
		migrations.BaseMigration(),
		migrations.RoleMigration(),
		migrations.UserMigration(),
		migrations.SessionMigration(),
	}
	schema := strings.Join(collectSchema, ` `)

	result := db.MustExec(schema)

	logMigrate := helpers.PrintLog("Sqlx", "Migration Successfully")
	fmt.Println(logMigrate, result)
}

func seeder() {

	// initialize seeder
	roleSeeder()
	userSeeder()

	logMigrate := helpers.PrintLog("Sqlx", "Seeders Successfully")
	fmt.Println(logMigrate)
}

func roleSeeder() {
	var err error
	var result sql.Result

	role := []entities.RoleEntity{}
	err = db.Select(&role, `SELECT * FROM "role"`)

	if err != nil {
		fmt.Println(err)
	}

	if len(role) == 0 {
		// role seeder
		result, err = db.NamedExec(`INSERT INTO "role" (id, name)
						VALUES (:id, :name)`, seeds.RoleSeeds())
	}

	if err != nil {
		logErrMessage := helpers.PrintLog("Sqlx", "Failed to run role seeder", helpers.Options{Label: "error"})
		fmt.Println(logErrMessage, err)
	}

	logMessage := helpers.PrintLog("Sqlx", "Role Seeds")
	fmt.Println(logMessage, result)
}

func userSeeder() {
	var err error
	var result sql.Result

	data := []entities.UserEntity{}
	err = db.Select(&data, `SELECT * FROM "user"`)

	if err != nil {
		fmt.Println(err)
	}

	if len(data) == 0 {
		// user seeder
		result, err = db.NamedExec(`INSERT INTO "user" (fullname, email, password, phone, token_verify, address, is_active, is_blocked, role_id)
						VALUES (:fullname, :email, :password, :phone, :token_verify, :address, :is_active, :is_blocked, :role_id)`, seeds.UserSeeds())
	}

	if err != nil {
		logErrMessage := helpers.PrintLog("Sqlx", "Failed to run user seeder", helpers.Options{Label: "error"})
		fmt.Println(logErrMessage, err)
	}

	logMessage := helpers.PrintLog("Sqlx", "User Seeds")
	fmt.Println(logMessage, result)
}
