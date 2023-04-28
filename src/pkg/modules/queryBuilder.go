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

func QueryBuilder(table string, c *fiber.Ctx) (*sqlf.Stmt, int) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	ctx := context.Background()
	db := config.GetDB()

	var total int
	var err error

	qPage := c.Query("page")
	qPageSize := c.Query("pageSize")
	qFiltered := c.Query("filtered")
	qSorted := c.Query("sorted")

	if qPage == "" {
		qPage = "1"
	}

	if qPageSize == "" {
		qPageSize = "10"
	}

	page, _ := strconv.Atoi(qPage)
	pageSize, _ := strconv.Atoi(qPageSize)

	// query record to database
	qRecord := sqlf.From(table).Select("*").Paginate(page, pageSize)

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

			queryFilter := fmt.Sprintln(filterId + " ILIKE '%" + filterValue + "%' ")
			qRecord.Where(queryFilter)
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
