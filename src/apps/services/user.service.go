package services

import (
	"context"
	"database/sql"
	"fmt"
	"gofi/src/database/entities"
	"gofi/src/pkg/helpers"
	"gofi/src/pkg/modules"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/leporo/sqlf"
)

type UserService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{db}
}

func (s *UserService) FindAll(c *fiber.Ctx) ([]entities.UserEntity, int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	ctx := context.Background()

	var data []entities.UserEntity
	var err error

	// get query builder
	qRecord, total := modules.QueryBuilder("user", c)
	qRecordLog := helpers.PrintLog("Sqlf", qRecord.String())

	fmt.Println(qRecordLog)

	// check query record
	err = qRecord.QueryAndClose(ctx, s.db, func(rows *sql.Rows) {
		var record entities.UserEntity

		// Scan Record
		err = rows.Scan(&record.Id, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt, &record.Fullname, &record.Email, &record.Password, &record.Phone, &record.TokenVerify, &record.Address, &record.IsActive, &record.IsBlocked, &record.RoleId)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		data = append(data, record)
	})

	if err != nil {
		fmt.Println(err)
		return data, total, err
	}

	return data, total, nil
}
