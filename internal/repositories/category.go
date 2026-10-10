package repositories

import (
	"context"
	"fmt"

	"gofi/internal/models"

	"github.com/google/uuid"
)

// CategoryRepository manages categories.
type CategoryRepository struct {
	BaseRepository
}

const categoryColumns = `
	id, created_at, updated_at, deleted_at, name`

func scanCategory(row rowScanner) (models.Category, error) {
	var c models.Category
	err := row.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.Name)
	return c, translateNotFound("categories", err)
}

// Insert persists a new category.
func (r *CategoryRepository) Insert(ctx context.Context, c models.Category) error {
	_, err := r.execContext(ctx, `
		INSERT INTO categories (id, name)
		VALUES ($1, $2)`,
		c.ID, c.Name,
	)
	return err
}

// categorySortColumns whitelists the columns a `sorted` value may order by.
var categorySortColumns = map[string]string{
	"name":       "name",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

// categoryFilters builds the WHERE conditions shared by List and Count.
func categoryFilters(p ListParams) *whereBuilder {
	w := newWhere(`deleted_at IS NULL`)
	w.search(p.Search, "name")
	return w
}

// List returns one page of live categories, optionally filtered by a
// case-insensitive name search and ordered by p.Sorted. Callers clamp the
// offset/limit first (see dtos.ListQuery.Clamp).
func (r *CategoryRepository) List(ctx context.Context, p ListParams) ([]models.Category, error) {
	where, args, limit, offset := categoryFilters(p).page(p)

	query := fmt.Sprintf(`
		SELECT %s FROM categories
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		categoryColumns, where,
		resolveOrder(p.Sorted, categorySortColumns, "created_at DESC"), limit, offset)

	rows, err := r.queryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return collectRows(rows, scanCategory)
}

// Count returns how many live categories match the List filters, for the
// response meta block.
func (r *CategoryRepository) Count(ctx context.Context, p ListParams) (int64, error) {
	return r.countRows(ctx, "categories", categoryFilters(p))
}

// Get returns one live category by id.
func (r *CategoryRepository) Get(ctx context.Context, id uuid.UUID) (models.Category, error) {
	query := fmt.Sprintf(`SELECT %s FROM categories WHERE id = $1 AND deleted_at IS NULL`, categoryColumns)
	return scanCategory(r.queryRowContext(ctx, query, id))
}

// Update rewrites a category.
func (r *CategoryRepository) Update(ctx context.Context, c models.Category) error {
	res, err := r.execContext(ctx, `
		UPDATE categories
		SET name = $2, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		c.ID, c.Name,
	)
	if err != nil {
		return err
	}
	return requireAffected(res, "category")
}

// Delete soft-removes one category by stamping deleted_at.
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.execContext(ctx, `
		UPDATE categories
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	return requireAffected(res, "category")
}

// Restore clears deleted_at on a soft-deleted category, returning
// ErrNotFound when the row is absent or not trashed.
func (r *CategoryRepository) Restore(ctx context.Context, id uuid.UUID) error {
	res, err := r.execContext(ctx, `
		UPDATE categories
		SET deleted_at = NULL, updated_at = now()
		WHERE id = $1 AND deleted_at IS NOT NULL`,
		id,
	)
	if err != nil {
		return err
	}
	return requireAffected(res, "category")
}

// ForceDelete permanently removes the row, whether or not it is trashed.
// FK violations from referencing tables surface as a conflict.
func (r *CategoryRepository) ForceDelete(ctx context.Context, id uuid.UUID) error {
	res, err := r.execContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return translateForeignKey(err, "category is still referenced by other records")
	}
	return requireAffected(res, "category")
}
