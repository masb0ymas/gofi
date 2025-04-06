package handler

import (
	"gofi/lib"
	"gofi/lib/constant"
	"gofi/middleware"
	"gofi/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/masb0ymas/go-utils/pkg"
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
	// only admin can access
	adminOnly := []string{constant.ID_SUPER_ADMIN, constant.ID_ADMIN}

	new_route := route.Group("/role", middleware.Authorization(), middleware.PermissionAccess(adminOnly))
	new_route.Get("/", h.GetAllRoles)
	new_route.Get("/:id", h.GetRoleById)
	new_route.Post("/", h.CreateRole)
	new_route.Put("/:id", h.UpdateRole)
	new_route.Put("/restore/:id", h.RestoreRole)
	new_route.Delete("/soft-delete/:id", h.SoftDeleteRole)
	new_route.Delete("/force-delete/:id", h.ForceDeleteRole)
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	req := &lib.Pagination{
		Page:     pkg.StringToInt32(c.Query("page", "1")),
		PageSize: pkg.StringToInt32(c.Query("page_size", "10")),
		Filtered: lib.ParseFilterItems(c.Query("filtered", "[]")),
		Sorted:   lib.ParseSortItems(c.Query("sorted", "[]")),
	}

	records, total, err := h.service.GetAllRoles(req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendGetResponse(c, "data has been received", records, lib.Paginate(req.Page, req.PageSize, total))
}

func (h *RoleHandler) GetRoleById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid role id")
	}

	record, err := h.service.GetRoleById(id)
	if err != nil {
		return lib.SendNotFoundResponse(c, "role not found")
	}

	return lib.SendSuccessResponse(c, "data has been received", record)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req service.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	record, err := h.service.CreateRole(&req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendCreatedResponse(c, "data has been added", record)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid role id")
	}

	var req service.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	record, err := h.service.UpdateRole(id, &req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been updated", record)
}

func (h *RoleHandler) RestoreRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid role id")
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
		return lib.SendBadRequestResponse(c, "invalid role id")
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
		return lib.SendBadRequestResponse(c, "invalid role id")
	}

	err = h.service.ForceDeleteRole(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been permanently deleted from the system", nil)
}
