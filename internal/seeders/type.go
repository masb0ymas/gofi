package seeders

// Seeder loads one data set. Seed is idempotent and applies completely or
// not at all.
type Seeder interface {
	Name() string
	Seed() error
}
