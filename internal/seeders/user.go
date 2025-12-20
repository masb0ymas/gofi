package seeders

import (
	"database/sql"
	"time"

	"gofi/internal/lib"
	"gofi/internal/lib/constant"
	"gofi/internal/models"
	"gofi/internal/repositories"

	"github.com/google/uuid"
)

type UserSeeder struct {
	DB *sql.DB
}

func (s UserSeeder) Name() string {
	return "user"
}

var defaultPassword = "password"

func (s UserSeeder) Seed() {
	users := []*models.User{
		{
			Base: models.Base{
				ID: uuid.Must(uuid.NewV7()),
			},
			FirstName: "Admin",
			LastName:  lib.StringPtr("System"),
			Email:     "admin@localhost.test",
			Phone:     nil,
			Password:  lib.StringPtr(defaultPassword),
			ActiveAt:  lib.TimePtr(time.Now()),
			RoleID:    uuid.MustParse(constant.RoleAdmin),
			UploadID:  nil,
		},
		{
			Base: models.Base{
				ID: uuid.Must(uuid.NewV7()),
			},
			FirstName: "User",
			LastName:  nil,
			Email:     "user@localhost.test",
			Phone:     nil,
			Password:  lib.StringPtr(defaultPassword),
			ActiveAt:  lib.TimePtr(time.Now()),
			RoleID:    uuid.MustParse(constant.RoleUser),
			UploadID:  nil,
		},
	}

	userRepo := repositories.UserRepository{
		BaseRepository: repositories.BaseRepository{
			DB:        s.DB,
			TableName: "users",
		},
	}
	err := userRepo.Insert(users...)
	if err != nil {
		panic(NewErrSeedingFailed(err))
	}
}
