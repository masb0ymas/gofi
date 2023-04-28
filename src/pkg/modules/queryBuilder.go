package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"gofi/src/pkg/config"
	"gofi/src/pkg/helpers"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leporo/sqlf"
)

type Filtered struct {
	ID    string
	Value string
}

type Sorted struct {
	ID    string
	Order string
}

type QueryBuilderOptions struct {
	Limit int
}

/*
Query Builder with Sqlf
*/
func QueryBuilder(table string, c *fiber.Ctx, options ...QueryBuilderOptions) (*sqlf.Stmt, int) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	ctx := context.Background()
	db := config.GetDB()

	var limitQuery int
	var total int
	var err error

	minLimitQuery := 10
	maxLimitQuery := 1000

	qPage := c.Query("page")
	qPageSize := c.Query("pageSize")
	qFiltered := c.Query("filtered")
	qSorted := c.Query("sorted")

	if qPage == "" {
		qPage = "1"
	}

	if qPageSize == "" {
		qPageSize = string(rune(minLimitQuery))
	}

	page, _ := strconv.Atoi(qPage)
	pageSize, _ := strconv.Atoi(qPageSize)

	for _, opt := range options {
		// increase max limit query
		if opt.Limit > 0 {
			maxLimitQuery = opt.Limit
		}
	}

	if pageSize <= 0 {
		limitQuery = minLimitQuery
	} else if pageSize > maxLimitQuery {
		limitQuery = maxLimitQuery
	} else {
		limitQuery = pageSize
	}

	// query record to database
	qRecord := sqlf.From(`public.`+table+``).Select("*").Paginate(page, limitQuery)

	// log query offset & limit
	logOffsetLimit := fmt.Sprintf("OFFSET %d LIMIT %d", page, limitQuery)
	fmt.Println(helpers.PrintLog("Sqlf", logOffsetLimit))

	// check query filtered not empty
	if qFiltered != "" {
		var filtered []Filtered

		err = json.Unmarshal([]byte(qFiltered), &filtered)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		// looping query filtered array
		for k := range filtered {
			filterId := filtered[k].ID
			filterValue := filtered[k].Value

			checkUUID := helpers.IsValidUUID(filterValue)
			checkNumber := helpers.IsDigit(filterValue)
			checkBool, _ := strconv.ParseBool(filterValue)

			if checkUUID {
				// check value is UUID
				queryFilter := fmt.Sprintf("%s = '%s'", filterId, filterValue)
				qRecord.Where(queryFilter)
			} else if !checkUUID && checkNumber {
				// check value is number
				queryFilter := fmt.Sprintf("%s = '%s'", filterId, filterValue)
				qRecord.Where(queryFilter)
			} else if checkBool {
				// check value is boolean
				queryFilter := fmt.Sprintf("%s = '%s'", filterId, filterValue)
				qRecord.Where(queryFilter)
			} else {
				// default query LIKE
				queryFilter := fmt.Sprintln(filterId + " ILIKE '%" + filterValue + "%'")
				qRecord.Where(queryFilter)
			}
		}
	}

	// check query sorted not empty
	if qSorted != "" {
		var sorted []Sorted

		err = json.Unmarshal([]byte(qSorted), &sorted)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		// looping query sorted array
		for k := range sorted {
			sortedId := sorted[k].ID
			sortedOrder := sorted[k].Order

			queryOrder := fmt.Sprintf(`%s %s`, sortedId, sortedOrder)
			qRecord.OrderBy(queryOrder)
		}
	} else {
		// default query sorted
		qRecord.OrderBy("created_at DESC")
	}

	// query count total record
	qTotal := sqlf.From(table).Select("COUNT(*)").To(&total)

	qTotalLog := helpers.PrintLog("Sqlf", qTotal.String())
	fmt.Println(qTotalLog)

	// check query count
	err = qTotal.QueryRowAndClose(ctx, db)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return qRecord, total
}
