package repository

import (
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

func (r *Repository[T]) FindAll() ([]T, error) {
	var entities []T
	query := "SELECT * FROM " + r.tableName + " WHERE deleted_at IS NULL"
	err := r.db.Select(&entities, query)
	return entities, err
}

func (r *Repository[T]) FindById(id uuid.UUID) (*T, error) {
	var entity T
	query := "SELECT * FROM " + r.tableName + " WHERE id = $1 AND deleted_at IS NULL"
	err := r.db.Get(&entity, query, id)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) Create(entity *T) error {
	query := "INSERT INTO " + r.tableName + " VALUES (:entity) RETURNING *"
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"entity": entity,
	})
	return err
}

func (r *Repository[T]) Update(id uuid.UUID, entity *T) error {
	query := "UPDATE " + r.tableName + " SET entity = :entity WHERE id = :id AND deleted_at IS NULL"
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":     id,
		"entity": entity,
	})
	return err
}

func (r *Repository[T]) Restore(id uuid.UUID) error {
	query := "UPDATE " + r.tableName + " SET deleted_at = NULL WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository[T]) SoftDelete(id uuid.UUID) error {
	query := "UPDATE " + r.tableName + " SET deleted_at = NOW() WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository[T]) ForceDelete(id uuid.UUID) error {
	query := "DELETE FROM " + r.tableName + " WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
