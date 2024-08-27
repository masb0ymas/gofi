package main

import (
	"gofi/config"
	"gofi/database"
	"gofi/routes"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var (
	port         = config.Env("APP_PORT", "8080")
	dbname       = config.Env("DB_DATABASE", "db_example")
	ratelimit, _ = strconv.Atoi(config.Env("APP_RATE_LIMIT", "100"))
)

func main() {
	// database instance
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Printf("successfully connected to database %v", dbname)

	// fiber instance
	app := fiber.New()

	// use middleware
	app.Use(cors.New(config.Cors()))
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{Max: ratelimit}))
	app.Use(requestid.New())
	app.Use(recover.New())

	// static file
	app.Static("/", "./public")

	// initial routes
	routes.Routes(db.GetDB(), app)

	// listen app
	log.Fatal(app.Listen(":" + port))
}
