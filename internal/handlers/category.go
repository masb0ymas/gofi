package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"gofi/internal/app"
	"gofi/internal/dtos"
	"gofi/internal/lib"
	"gofi/internal/models"
)

type categoryHandler struct {
	app *app.Application
}

func (h *categoryHandler) List(c fiber.Ctx) error {
	var q dtos.ListQuery
	if err := lib.ValidateRequestQuery(c, &q); err != nil {
		return err
	}
	q.Clamp()

	params := listParams(q)

	repo := &h.app.Repositories.Category
	organizations, err := repo.List(c.Context(), params)
	if err != nil {
		return err
	}
	total, err := repo.Count(c.Context(), params)
	if err != nil {
		return err
	}

	return dtos.List(c, organizations, dtos.ListMeta(int(total), q.Offset, q.Limit))
}

func (h *categoryHandler) Get(c fiber.Ctx) error {
	id, err := pathID(c, "category")
	if err != nil {
		return err
	}

	org, err := h.app.Repositories.Category.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return dtos.Get(c, org)
}

func (h *categoryHandler) Create(c fiber.Ctx) error {
	var body dtos.CategoryBody
	if err := lib.ValidateRequestBody(c, &body); err != nil {
		return err
	}

	org := models.Category{Name: body.Name}
	var err error
	if org.ID, err = uuid.NewV7(); err != nil {
		return err
	}

	repo := &h.app.Repositories.Category
	if err := repo.Insert(c.Context(), org); err != nil {
		return err
	}

	// Re-read so the response carries the DB-set timestamps.
	org, err = repo.Get(c.Context(), org.ID)
	if err != nil {
		return err
	}
	return dtos.Created(c, org, "category created")
}

func (h *categoryHandler) Update(c fiber.Ctx) error {
	id, err := pathID(c, "category")
	if err != nil {
		return err
	}

	var body dtos.CategoryBody
	if err := lib.ValidateRequestBody(c, &body); err != nil {
		return err
	}

	org := models.Category{Name: body.Name}
	org.ID = id

	repo := &h.app.Repositories.Category
	if err := repo.Update(c.Context(), org); err != nil {
		return err
	}

	org, err = repo.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return dtos.OK(c, org)
}

func (h *categoryHandler) Delete(c fiber.Ctx) error {
	return mutateByID(c, "category", h.app.Repositories.Category.Delete, "category deleted")
}

func (h *categoryHandler) Restore(c fiber.Ctx) error {
	return mutateByID(c, "category", h.app.Repositories.Category.Restore, "category restored")
}

func (h *categoryHandler) DeletePermanently(c fiber.Ctx) error {
	return mutateByID(c, "category", h.app.Repositories.Category.ForceDelete, "category permanently deleted")
}
