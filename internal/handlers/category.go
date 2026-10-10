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
	categories, err := repo.List(c.Context(), params)
	if err != nil {
		return err
	}
	total, err := repo.Count(c.Context(), params)
	if err != nil {
		return err
	}

	return dtos.List(c, categories, dtos.ListMeta(int(total), q.Offset, q.Limit))
}

func (h *categoryHandler) Get(c fiber.Ctx) error {
	id, err := pathID(c, "category")
	if err != nil {
		return err
	}

	category, err := h.app.Repositories.Category.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return dtos.Get(c, category)
}

func (h *categoryHandler) Create(c fiber.Ctx) error {
	var body dtos.CategoryBody
	if err := lib.ValidateRequestBody(c, &body); err != nil {
		return err
	}

	category := models.Category{Name: body.Name}
	var err error
	if category.ID, err = uuid.NewV7(); err != nil {
		return err
	}

	repo := &h.app.Repositories.Category
	if err := repo.Insert(c.Context(), category); err != nil {
		return err
	}

	// Re-read so the response carries the DB-set timestamps.
	category, err = repo.Get(c.Context(), category.ID)
	if err != nil {
		return err
	}
	return dtos.Created(c, category, "category created")
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

	category := models.Category{Name: body.Name}
	category.ID = id

	repo := &h.app.Repositories.Category
	if err := repo.Update(c.Context(), category); err != nil {
		return err
	}

	category, err = repo.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return dtos.OK(c, category)
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
