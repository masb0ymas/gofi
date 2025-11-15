package seeders

import (
	"errors"
	"fmt"
)

var ErrSeedingFailed = errors.New("seeding failed")

func NewErrSeedingFailed(err error) error {
	return fmt.Errorf("%w: %w", ErrSeedingFailed, err)
}
