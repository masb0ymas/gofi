package handlers

import (
	"net/http"

	"gofi/internal/app"
	"gofi/internal/dto"
	"gofi/internal/lib"
	"gofi/internal/models"
	"gofi/internal/repositories"
	"gofi/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type roleHandler struct {
	app *app.Application
}

func (h *roleHandler) Index(c *fiber.Ctx) error {
	var dto dto.RolePagination

	if err := lib.ValidateRequestQuery(c, &dto); err != nil {
		switch e := err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(lib.WrapValidationError(e.MessageRecord))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	opts := &repositories.QueryOptions{
		Offset: dto.Offset,
		Limit:  dto.Limit,
	}

	roles, meta, err := h.app.Repositories.Role.List(opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseMultiData[*models.Role]{
			Message: "list data has been retrieved successfully",
			Data:    roles,
			Meta: fiber.Map{
				"total": meta.Total,
			},
		})
}

func (h *roleHandler) Show(c *fiber.Ctx) error {
	roleID, err := lib.ContextParamUUID(c, "roleID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid role id must be uuid format",
			"error":   err.Error(),
		})
	}

	role, err := h.app.Repositories.Role.Get(roleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.Role]{
			Message: "get data has been retrieved successfully",
			Data:    role,
		})
}

func (h *roleHandler) Create(c *fiber.Ctx) error {
	var dto dto.RoleCreate

	if err := lib.ValidateRequestBody(c, &dto); err != nil {
		switch e := err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(lib.WrapValidationError(e.MessageRecord))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	roleID, err := uuid.NewV7()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	role := &models.Role{
		Base: models.Base{
			ID: roleID,
		},
		Name: dto.Name,
	}

	err = h.app.Repositories.Role.Insert(role)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.Role]{
			Message: "data has been created successfully",
			Data:    role,
		})
}

func (h *roleHandler) Update(c *fiber.Ctx) error {
	roleID, err := lib.ContextParamUUID(c, "roleID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid role id must be uuid format",
			"error":   err.Error(),
		})
	}

	var dto dto.RoleUpdate

	if err := lib.ValidateRequestBody(c, &dto); err != nil {
		switch e := err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(lib.WrapValidationError(e.MessageRecord))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	role, err := h.app.Repositories.Role.Get(roleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if dto.Name != "" {
		role.Name = dto.Name
	}

	err = h.app.Repositories.Role.Update(roleID, role)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.Role]{
			Message: "data has been updated successfully",
			Data:    role,
		})
}

func (h *roleHandler) Delete(c *fiber.Ctx) error {
	roleID, err := lib.ContextParamUUID(c, "roleID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid role id must be uuid format",
			"error":   err.Error(),
		})
	}

	err = h.app.Repositories.Role.Delete(roleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.Role]{
			Message: "data has been deleted successfully",
		})
}

func (h *roleHandler) SoftDelete(c *fiber.Ctx) error {
	roleID, err := lib.ContextParamUUID(c, "roleID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid role id must be uuid format",
			"error":   err.Error(),
		})
	}

	err = h.app.Repositories.Role.SoftDelete(roleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.Role]{
			Message: "data has been soft deleted successfully",
		})
}

func (h *roleHandler) Restore(c *fiber.Ctx) error {
	roleID, err := lib.ContextParamUUID(c, "roleID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid role id must be uuid format",
			"error":   err.Error(),
		})
	}

	err = h.app.Repositories.Role.Restore(roleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.Role]{
			Message: "data has been restored successfully",
		})
}
