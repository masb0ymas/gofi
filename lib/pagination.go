package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/masb0ymas/go-utils/pkg"
)

type Pagination struct {
	Page     int32        `json:"page"`
	PageSize int32        `json:"page_size"`
	Filtered []FilterItem `json:"filtered"`
	Sorted   []SortItem   `json:"sorted"`
}

type FilterItem struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
}

type SortItem struct {
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func calculatePageSize(page_size int32, limit int32) int32 {
	minLimit := int32(10)

	if page_size > 0 {
		return int32(math.Min(float64(page_size), float64(limit)))
	}

	return minLimit
}

// Apply pagination to the query
//
//	example:
//		query := "SELECT * FROM table"
//		page := 1
//		page_size := 10
//		new_query := applyPagination(query, page, page_size)
func applyPagination(query string, page int32, page_size int32) string {
	var new_page = page

	if page == 0 {
		new_page = 1
	}

	if page_size <= 0 {
		page_size = 10
	}

	limit := calculatePageSize(page_size, 100)
	offset := (new_page - 1) * limit

	return fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
}

// Apply filtering to the query
//
//	example:
//		query := "SELECT * FROM table"
//		filtered := []FilterItem{
//			{ID: "column1", Value: "value1"},
//			{ID: "column2", Value: "value2"},
//		}
//		new_query := applyFilters(query, filtered)
func applyFilters(query string, filtered []FilterItem) string {
	var new_query = query

	for _, filter := range filtered {
		switch v := filter.Value.(type) {
		case int32:
			new_query = fmt.Sprintf("%s AND %s = :%d", new_query, filter.ID, v)

		case bool:
			new_query = fmt.Sprintf("%s AND %s = %t", new_query, filter.ID, v)

		case string:
			// Try to parse UUID if the string is in UUID format
			if uuidStr, ok := filter.Value.(string); ok {
				if _, err := uuid.Parse(uuidStr); err == nil {
					new_query = fmt.Sprintf("%s AND %s = '%s'", new_query, filter.ID, uuidStr)
					continue
				}
			}

			// If not UUID, treat as regular string with ILIKE
			new_query = fmt.Sprintf("%s AND %s ILIKE '%%%s%%'", new_query, filter.ID, v)

		default:
			// Handle other types or skip
			continue
		}
	}

	return new_query
}

// Apply sorting to the query
//
//	example:
//		query := "SELECT * FROM table"
//		sorted := []SortItem{
//			{Sort: "column1", Order: "ASC"},
//			{Sort: "column2", Order: "DESC"},
//		}
//		new_query := applySorting(query, sorted)
func applySorting(query string, sorted []SortItem) string {
	var new_query = query

	for _, sort := range sorted {
		new_query = fmt.Sprintf("%s ORDER BY %s %s", new_query, sort.Sort, sort.Order)
	}

	return new_query
}

// Query Builder
//
//	example:
//		query := "SELECT * FROM table"
//		filtered := []FilterItem{
//			{ID: "column1", Value: "value1"},
//			{ID: "column2", Value: "value2"},
//		}
//		sorted := []SortItem{
//			{Sort: "column1", Order: "ASC"},
//			{Sort: "column2", Order: "DESC"},
//		}
//		page := 1
//		page_size := 10
//		new_query := QueryBuilder(query, filtered, sorted, page, page_size)
func QueryBuilder(query string, filtered []FilterItem, sorted []SortItem, page int32, page_size int32) string {
	new_query := applyFilters(query, filtered)
	new_query = applySorting(new_query, sorted)
	new_query = applyPagination(new_query, page, page_size)

	message := pkg.Println("Query:", new_query)
	log.Println(message)

	return new_query
}

// Query Builder
//
//	example:
//		query := "SELECT * FROM table"
//		filtered := []FilterItem{
//			{ID: "column1", Value: "value1"},
//			{ID: "column2", Value: "value2"},
//		}
//		sorted := []SortItem{
//			{Sort: "column1", Order: "ASC"},
//			{Sort: "column2", Order: "DESC"},
//		}
func QueryBuilderForCount(query string, filtered []FilterItem) string {
	new_query := applyFilters(query, filtered)

	message := pkg.Println("Query Count:", new_query)
	log.Println(message)

	return new_query
}

// ParseFilterItems parses a JSON string into a slice of FilterItem
func ParseFilterItems(jsonStr string) []FilterItem {
	var filtered []FilterItem
	err := json.Unmarshal([]byte(jsonStr), &filtered)
	if err != nil {
		return []FilterItem{}
	}
	return filtered
}

// ParseSortItems parses a JSON string into a slice of SortItem
func ParseSortItems(jsonStr string) []SortItem {
	var sorted []SortItem
	err := json.Unmarshal([]byte(jsonStr), &sorted)
	if err != nil {
		return []SortItem{}
	}
	return sorted
}

// Paginate creates a paginated response
func Paginate(page int32, pageSize int32, total int64) map[string]interface{} {
	var totalPages int32
	if pageSize > 0 {
		totalPages = int32(math.Ceil(float64(total) / float64(pageSize)))
	} else {
		totalPages = 1
	}

	return map[string]interface{}{
		"total":        total,
		"total_pages":  totalPages,
		"current_page": page,
		"page_size":    pageSize,
	}
}
