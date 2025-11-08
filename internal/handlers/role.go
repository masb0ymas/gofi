package handlers

import (
	"net/http"

	"gofi/internal/app"
	"gofi/internal/dto"
	"gofi/internal/lib"
	"gofi/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type roleHandler struct {
	app *app.Application
}

func (h *roleHandler) Index(c *fiber.Ctx) error {
	var dto dto.RolePagination

	if err := lib.ValidateRequestQuery(c, &dto); err != nil {
		switch err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
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

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "data has been retrieved successfully",
		"data":    roles,
		"meta":    meta,
	})
}
