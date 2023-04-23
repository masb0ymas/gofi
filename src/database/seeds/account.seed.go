package seeds

import (
	"fmt"
	"gofi/src/database/entities"
	"gofi/src/pkg/constants"
	"gofi/src/pkg/helpers"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// Role Seeds
func RoleSeeds() []entities.RoleEntity {
	roleStruct := []entities.RoleEntity{
		{
			BaseEntity: entities.BaseEntity{
				Id: uuid.MustParse(constants.ROLE_SUPER_ADMIN),
			},
			Name: "Super Admin",
		},
		{
			BaseEntity: entities.BaseEntity{
				Id: uuid.MustParse(constants.ROLE_ADMIN),
			},
			Name: "Admin",
		},
		{
			BaseEntity: entities.BaseEntity{
				Id: uuid.MustParse(constants.ROLE_USER),
			},
			Name: "User",
		},
	}

	return roleStruct
}

// User Seeds
func UserSeeds() []entities.UserEntity {
	defaultPassword := "Padang123"

	hashedPassword, err := helpers.HashPassword(defaultPassword)

	if err != nil {
		fmt.Println(err)
	}

	userStruct := []entities.UserEntity{
		{
			Fullname:    "Super Admin",
			Email:       "super.admin@mail.com",
			Password:    string(hashedPassword),
			Phone:       null.String{},
			TokenVerify: null.String{},
			Address:     null.String{},
			IsActive:    true,
			IsBlocked:   false,
			RoleId:      uuid.MustParse(constants.ROLE_SUPER_ADMIN),
		},
		{
			Fullname:    "Admin",
			Email:       "admin@mail.com",
			Password:    string(hashedPassword),
			Phone:       null.String{},
			TokenVerify: null.String{},
			Address:     null.String{},
			IsActive:    true,
			IsBlocked:   false,
			RoleId:      uuid.MustParse(constants.ROLE_ADMIN),
		},
		{
			Fullname:    "User",
			Email:       "user@mail.com",
			Password:    string(hashedPassword),
			Phone:       null.String{},
			TokenVerify: null.String{},
			Address:     null.String{},
			IsActive:    true,
			IsBlocked:   false,
			RoleId:      uuid.MustParse(constants.ROLE_USER),
		},
	}

	return userStruct
}
