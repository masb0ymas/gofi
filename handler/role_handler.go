package handler

import (
	"goarif-api/lib"
	"goarif-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

func (h *RoleHandler) RegisterRoutes(route fiber.Router) {
	role := route.Group("/role")
	role.Get("/", h.GetAllRoles)
	role.Get("/:id", h.GetRoleById)
	role.Post("/", h.CreateRole)
	role.Put("/:id", h.UpdateRole)
	role.Delete("/:id", h.SoftDeleteRole)
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been received", roles)
}

func (h *RoleHandler) GetRoleById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "Invalid role ID")
	}

	role, err := h.service.GetRoleById(id)
	if err != nil {
		return lib.SendNotFoundResponse(c, "Role not found")
	}

	return lib.SendSuccessResponse(c, "data has been received", role)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req service.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "Invalid request body")
	}

	role, err := h.service.CreateRole(&req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendCreatedResponse(c, "data has been added", role)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "Invalid role ID")
	}

	var req service.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "Invalid request body")
	}

	role, err := h.service.UpdateRole(id, &req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been updated", role)
}

func (h *RoleHandler) RestoreRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "Invalid role ID")
	}

	err = h.service.RestoreRole(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been restored", nil)
}

func (h *RoleHandler) SoftDeleteRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "Invalid role ID")
	}

	err = h.service.SoftDeleteRole(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been deleted", nil)
}

func (h *RoleHandler) ForceDeleteRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "Invalid role ID")
	}

	err = h.service.ForceDeleteRole(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been permanently deleted from the system", nil)
}
