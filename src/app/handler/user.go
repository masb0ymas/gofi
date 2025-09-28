package handler

import (
	"gofi/src/app/middleware"
	"gofi/src/app/service"
	"gofi/src/lib"
	"gofi/src/lib/constant"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/masb0ymas/go-utils/pkg"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(route fiber.Router) {
	// only admin can access
	adminOnly := []string{constant.ID_SUPER_ADMIN, constant.ID_ADMIN}

	new_route := route.Group("/user", middleware.Authorization(), middleware.PermissionAccess(adminOnly))
	new_route.Get("/", h.GetAllUsers)
	new_route.Get("/:id", h.GetUserById)
	new_route.Post("/", h.CreateUser)
	new_route.Put("/:id", h.UpdateUser)
	new_route.Put("/restore/:id", h.RestoreUser)
	new_route.Delete("/soft-delete/:id", h.SoftDeleteUser)
	new_route.Delete("/force-delete/:id", h.ForceDeleteUser)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	req := &lib.Pagination{
		Page:     pkg.StringToInt32(c.Query("page", "1")),
		PageSize: pkg.StringToInt32(c.Query("page_size", "10")),
		Filtered: lib.ParseFilterItems(c.Query("filtered", "[]")),
		Sorted:   lib.ParseSortItems(c.Query("sorted", "[]")),
	}

	records, total, err := h.service.GetAllUsers(req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendGetResponse(c, "data has been received", records, lib.Paginate(req.Page, req.PageSize, total))
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid user id")
	}

	record, err := h.service.GetUserById(id)
	if err != nil {
		return lib.SendNotFoundResponse(c, "user not found")
	}

	return lib.SendSuccessResponse(c, "data has been received", record)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	record, err := h.service.CreateUser(&req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendCreatedResponse(c, "data has been added", record)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid user id")
	}

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	record, err := h.service.UpdateUser(id, &req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been updated", record)
}

func (h *UserHandler) RestoreUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid user id")
	}

	err = h.service.RestoreUser(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been restored", nil)
}

func (h *UserHandler) SoftDeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid user id")
	}

	err = h.service.SoftDeleteUser(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been deleted", nil)
}

func (h *UserHandler) ForceDeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid user id")
	}

	err = h.service.ForceDeleteUser(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been permanently deleted from the system", nil)
}
