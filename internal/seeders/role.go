package seeders

import (
	"database/sql"

	"gofi/internal/lib/constant"
	"gofi/internal/models"
	"gofi/internal/repositories"

	"github.com/google/uuid"
)

type RoleSeeder struct {
	DB *sql.DB
}

func (s RoleSeeder) Name() string {
	return "role"
}

func (s RoleSeeder) Seed() {
	roles := []*models.Role{
		{
			Base: models.Base{
				ID: uuid.MustParse(constant.RoleAdmin),
			},
			Name: "admin",
		},
		{
			Base: models.Base{
				ID: uuid.MustParse(constant.RoleUser),
			},
			Name: "user",
		},
	}

	roleRepo := repositories.RoleRepository{DB: s.DB}
	err := roleRepo.Insert(roles...)
	if err != nil {
		panic(NewErrSeedingFailed(err))
	}
}
