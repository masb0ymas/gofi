package routes

import (
	"goarif-api/handler"
	"goarif-api/repository"
	"goarif-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func authRoute(db *sqlx.DB, route fiber.Router) {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(route)
}

func roleRoute(db *sqlx.DB, route fiber.Router) {
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)
	roleHandler.RegisterRoutes(route)
}

func v1Route(db *sqlx.DB, app *fiber.App) {
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	authRoute(db, v1)
	roleRoute(db, v1)
}
