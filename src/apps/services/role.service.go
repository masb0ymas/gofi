package services

import (
	"context"
	"database/sql"
	"fmt"
	"gofi/src/database/entities"
	"gofi/src/pkg/helpers"
	"strconv"

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

func (s *RoleService) FindAll(c *fiber.Ctx) ([]entities.RoleEntity, error) {
	var data []entities.RoleEntity

	queryPage := c.Query("page")
	queryPageSize := c.Query("pageSize")

	if queryPage == "" {
		queryPage = "1"
	}

	if queryPageSize == "" {
		queryPageSize = "10"
	}

	page, _ := strconv.Atoi(queryPage)
	pageSize, _ := strconv.Atoi(queryPageSize)

	skip := (page - 1) * pageSize

	sqlf.SetDialect(sqlf.PostgreSQL)
	ctx := context.Background()

	var record entities.RoleEntity

	// query builder
	query := sqlf.Select("*").From("role").
		OrderBy("created_at DESC").
		Offset(skip).
		Limit(pageSize)

	queryLog := helpers.PrintLog("Sqlf", query.String())
	fmt.Println(queryLog)

	err := query.QueryAndClose(ctx, s.db, func(rows *sql.Rows) {
		// Scan Record Rows
		rows.Scan(&record.Id, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt.Time, &record.Name)
		data = append(data, record)
	})

	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}
