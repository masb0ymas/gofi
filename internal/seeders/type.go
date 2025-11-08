package seeders

type Seeder interface {
	Name() string
	Seed()
}
