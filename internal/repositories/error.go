package repositories

import "errors"

var (
	ErrInsertDuplicate = errors.New("insert duplicate")
	ErrEditConflict    = errors.New("edit conflict")
	ErrNotFound        = errors.New("record not found")
)
