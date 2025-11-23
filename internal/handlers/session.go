package handlers

import (
	"gofi/internal/app"
	"gofi/internal/dto"
	"gofi/internal/lib"
	"gofi/internal/models"
	"gofi/internal/repositories"
	"gofi/internal/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type sessionHandler struct {
	app *app.Application
}

func (h *sessionHandler) Index(c *fiber.Ctx) error {
	var dto dto.SessionPagination

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

	sessions, meta, err := h.app.Repositories.Session.List(opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(
		types.ResponseMultiData[*models.Session]{
			Message: "list data has been retrieved successfully",
			Data:    sessions,
			Meta: fiber.Map{
				"total": meta.Total,
			},
		})
}
