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

type userHandler struct {
	app *app.Application
}

func (h *userHandler) Index(c *fiber.Ctx) error {
	var dto dto.UserPagination

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

	users, meta, err := h.app.Repositories.User.List(opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseMultiData[*models.User]{
			Message: "list data has been retrieved successfully",
			Data:    users,
			Meta: fiber.Map{
				"total": meta.Total,
			},
		})
}

func (h *userHandler) Show(c *fiber.Ctx) error {
	userID, err := lib.ContextParamUUID(c, "userID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid role id must be uuid format",
			"error":   err.Error(),
		})
	}

	user, err := h.app.Repositories.User.Get(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.User]{
			Message: "get data has been retrieved successfully",
			Data:    user,
		})
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	var dto dto.UserCreate

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

	userID, err := uuid.NewV7()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := &models.User{
		Base: models.Base{
			ID: userID,
		},
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Phone:     dto.Phone,
		Password:  dto.Password,
		RoleID:    dto.RoleID,
		UploadID:  dto.UploadID,
	}

	err = h.app.Repositories.User.Insert(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.User]{
			Message: "data has been created successfully",
			Data:    user,
		})
}

func (h *userHandler) Update(c *fiber.Ctx) error {
	userID, err := lib.ContextParamUUID(c, "userID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id must be uuid format",
			"error":   err.Error(),
		})
	}

	var dto dto.UserUpdate

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

	user, err := h.app.Repositories.User.Get(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}

	if dto.LastName != nil {
		user.LastName = dto.LastName
	}

	if dto.Phone != nil {
		user.Phone = dto.Phone
	}

	if dto.UploadID != nil {
		user.UploadID = dto.UploadID
	}

	err = h.app.Repositories.User.Update(userID, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.User]{
			Message: "data has been updated successfully",
			Data:    user,
		})
}

func (h *userHandler) Delete(c *fiber.Ctx) error {
	userID, err := lib.ContextParamUUID(c, "userID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id must be uuid format",
			"error":   err.Error(),
		})
	}

	err = h.app.Repositories.User.Delete(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.User]{
			Message: "data has been deleted successfully",
		})
}

func (h *userHandler) SoftDelete(c *fiber.Ctx) error {
	userID, err := lib.ContextParamUUID(c, "userID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id must be uuid format",
			"error":   err.Error(),
		})
	}

	err = h.app.Repositories.User.SoftDelete(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.User]{
			Message: "data has been soft deleted successfully",
		})
}

func (h *userHandler) Restore(c *fiber.Ctx) error {
	userID, err := lib.ContextParamUUID(c, "userID")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id must be uuid format",
			"error":   err.Error(),
		})
	}

	err = h.app.Repositories.User.Restore(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseSingleData[*models.User]{
			Message: "data has been restored successfully",
		})
}
