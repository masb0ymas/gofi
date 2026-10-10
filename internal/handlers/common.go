package handlers

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"gofi/internal/dtos"
	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
	"gofi/internal/repositories"
)

// pathID parses the :id path parameter, naming the resource in the error so a
// malformed id reads as "invalid category id".
func pathID(c fiber.Ctx, resource string) (uuid.UUID, error) {
	id, err := lib.ContextParamUUID(c, "id")
	if err != nil {
		return uuid.Nil, apperr.New(apperr.KindBadRequest, "invalid %s id", resource)
	}
	return id, nil
}

// listParams maps the clamped pagination query onto the repository controls.
func listParams(q dtos.ListQuery) repositories.ListParams {
	return repositories.ListParams{Offset: q.Offset, Limit: q.Limit, Search: q.Search, Sorted: q.Sorted}
}

// mutateByID runs a by-id write (delete, restore, permanent delete) on the
// :id path parameter and answers with the bare {message} envelope.
func mutateByID(c fiber.Ctx, resource string, fn func(context.Context, uuid.UUID) error, message string) error {
	id, err := pathID(c, resource)
	if err != nil {
		return err
	}
	if err := fn(c.Context(), id); err != nil {
		return err
	}
	return dtos.Message(c, fiber.StatusOK, message)
}
