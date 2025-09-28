package routes

import (
	"gofi/src/app/handler"
	"gofi/src/app/middleware"
	"gofi/src/app/repository"
	"gofi/src/app/service"

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

func userRoute(db *sqlx.DB, route fiber.Router) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userHandler.RegisterRoutes(route)
}

func uploadRoute(db *sqlx.DB, route fiber.Router) {
	uploadRepo := repository.NewUploadRepository(db)
	uploadService := service.NewUploadService(uploadRepo)
	uploadHandler := handler.NewUploadHandler(uploadService)
	uploadHandler.RegisterRoutes(route)
}

func v1Route(db *sqlx.DB, app *fiber.App) {
	v1 := app.Group("/v1", middleware.SentryMiddleware(), func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	// Routes
	authRoute(db, v1)
	roleRoute(db, v1)
	userRoute(db, v1)
	uploadRoute(db, v1)
}
