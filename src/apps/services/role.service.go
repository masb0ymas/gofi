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

type RoleService struct {
	db *sqlx.DB
}

func NewRoleService(db *sqlx.DB) *RoleService {
	return &RoleService{db}
}

func (s *RoleService) FindAll(c *fiber.Ctx) ([]entities.RoleEntity, int, error) {
	var data []entities.RoleEntity
	var err error

	sqlf.SetDialect(sqlf.PostgreSQL)
	ctx := context.Background()

	// get query builder
	qRecord, total := modules.QueryBuilder("role", c)
	qRecordLog := helpers.PrintLog("Sqlf", qRecord.String())

	fmt.Println(qRecordLog)

	// check query record
	err = qRecord.QueryAndClose(ctx, s.db, func(rows *sql.Rows) {
		var record entities.RoleEntity

		// Scan Record
		err = rows.Scan(&record.Id, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt, &record.Name)

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
