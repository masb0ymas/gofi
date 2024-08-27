package handler

import (
	"gofi/database/repository"
	"gofi/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

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
