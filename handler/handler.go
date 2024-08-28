package handler

import (
	"gofi/database/repository"
	"gofi/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func toTimePtr(t time.Time) time.Time {
	return t
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func isDateValue(date string) bool {
	_, err := time.Parse("01/02/2006", date)
	return err == nil
}

func RoleHandler(db *sqlx.DB, route fiber.Router) {
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := NewRoleHandler(roleService)

	r := route.Group("/role")
	r.Get("/", roleHandler.listRoles)
	r.Post("/", roleHandler.createRole)

	r_id := r.Group("/:id")
	r_id.Get("/", roleHandler.getRole)
	r_id.Put("/", roleHandler.updateRole)
	r_id.Delete("/", roleHandler.deleteRole)
}

func SessionHandler(db *sqlx.DB, route fiber.Router) {
	sessionRepo := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo)
	sessionHandler := NewSessionHandler(sessionService)

	r := route.Group("/session")
	r.Get("/", sessionHandler.listSessions)
	r.Post("/", sessionHandler.createSession)

	r_id := r.Group("/:id")
	r_id.Get("/", sessionHandler.getSession)
	r_id.Put("/", sessionHandler.updateSession)
	r_id.Delete("/", sessionHandler.deleteSession)
}
