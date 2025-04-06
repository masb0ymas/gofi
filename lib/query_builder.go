package lib

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// QueryBuilderIn generates a parameterized IN clause for SQL queries
func QueryBuilderIn(ids []uuid.UUID) string {
	if len(ids) == 0 {
		return "''"
	}

	placeholders := make([]string, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("'%s'", id.String())
	}
	
	return strings.Join(placeholders, ",")
}
