package repositories

import "errors"

var (
	ErrInsertDuplicate = errors.New("insert duplicate")
	ErrEditConflict    = errors.New("edit conflict")
	ErrRecordNotFound  = errors.New("record not found")
)
