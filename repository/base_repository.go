package repository

import (
	"gofi/lib"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BaseRepository[T any] interface {
	FindAll() ([]T, error)
	FindById(id uuid.UUID) (*T, error)
	Create(entity *T) error
	Update(id uuid.UUID, entity *T) error
	Restore(id uuid.UUID) error
	SoftDelete(id uuid.UUID) error
	ForceDelete(id uuid.UUID) error
}

type Repository[T any] struct {
	db        *sqlx.DB
	tableName string
}

func NewRepository[T any](db *sqlx.DB, tableName string) *Repository[T] {
	return &Repository[T]{
		db:        db,
		tableName: tableName,
	}
}

// Find All
func (r *Repository[T]) FindAll() ([]T, int64, error) {
	var entities []T
	var total int64

	query := `
		SELECT * FROM "` + r.tableName + `" WHERE "deleted_at" IS NULL
	`
	err := r.db.Select(&entities, query)
	if err != nil {
		return nil, 0, err
	}

	queryTotal := `
		SELECT COUNT(*) FROM "` + r.tableName + `" WHERE "deleted_at" IS NULL
	`
	err = r.db.Get(&total, queryTotal)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Find All With Pagination
func (r *Repository[T]) FindAllWithPagination(req *lib.Pagination) ([]T, int64, error) {
	var entities []T
	var total int64

	query := `
		SELECT * FROM "` + r.tableName + `" WHERE "deleted_at" IS NULL
	`
	query = lib.QueryBuilder(query, req.Filtered, req.Sorted, req.Page, req.PageSize)
	err := r.db.Select(&entities, query)
	if err != nil {
		return nil, 0, err
	}

	queryTotal := `
		SELECT COUNT(*) FROM "` + r.tableName + `" WHERE "deleted_at" IS NULL
	`
	queryTotal = lib.QueryBuilderForCount(queryTotal, req.Filtered)
	err = r.db.Get(&total, queryTotal)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Find By ID
func (r *Repository[T]) FindById(id uuid.UUID) (*T, error) {
	var entity T
	query := `
		SELECT * FROM "` + r.tableName + `" WHERE "id" = $1 AND "deleted_at" IS NULL
	`
	err := r.db.Get(&entity, query, id)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Find Record By ID (without deleted_at)
func (r *Repository[T]) FindRecordById(id uuid.UUID) (*T, error) {
	var entity T
	query := `
		SELECT * FROM "` + r.tableName + `" WHERE "id" = $1
	`
	err := r.db.Get(&entity, query, id)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Create
func (r *Repository[T]) Create(entity *T) error {
	query := `
		INSERT INTO "` + r.tableName + `" VALUES (:entity) RETURNING *
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"entity": entity,
	})
	return err
}

// Update
func (r *Repository[T]) Update(id uuid.UUID, entity *T) error {
	query := `
		UPDATE "` + r.tableName + `" SET entity = :entity WHERE "id" = :id AND "deleted_at" IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":     id,
		"entity": entity,
	})
	return err
}

// Restore
func (r *Repository[T]) Restore(id uuid.UUID) error {
	query := `
		UPDATE "` + r.tableName + `" SET "deleted_at" = NULL WHERE "id" = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}

// Soft Delete
func (r *Repository[T]) SoftDelete(id uuid.UUID) error {
	query := `
		UPDATE "` + r.tableName + `" SET "deleted_at" = NOW() WHERE "id" = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}

// Force Delete
func (r *Repository[T]) ForceDelete(id uuid.UUID) error {
	query := `
		DELETE FROM "` + r.tableName + `" WHERE "id" = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}
