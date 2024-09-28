package drivers

import (
	"github.com/mohamed-samir907/goquery/query"
)

type Driver interface {
	// Get executes the SELECT query and returns the result as a slice of maps.
	//
	// This method is useful for retrieving all rows from a SELECT query when you don't know
	// the exact structure of the result set in advance or when working with dynamic queries.
	//
	// Returns:
	//   - []map[string]any: A slice of maps, where each map represents a row from the query result.
	//     The keys in the map are column names, and the values are the corresponding data.
	//   - error: An error if any occurred during the query execution or result processing.
	//
	// Example:
	//
	//	rows, err := db.Table("users").
	//	    Where("age", ">", 18).
	//	    OrderBy("id", "DESC").
	//	    Limit(10).
	//	    Get()
	Get(q query.SelectQuery) ([]map[string]any, error)
	Insert()
	Update()
	Delete()
}

// prepareValuesAndScanArgs creates two slices to be used with the rows.Scan method.
//
// It takes a slice of column names and returns two slices:
// 1. values: A slice of empty interfaces to hold the scanned values.
// 2. scanArgs: A slice of pointers to the values slice, used as arguments for rows.Scan.
//
// This function is crucial for dynamically scanning rows with an unknown number of columns.
// It allows for flexible handling of different result set structures without hardcoding column names or types.
//
// Parameters:
//   - columns: A slice of strings representing the column names in the result set.
//
// Returns:
//   - []any: A slice of empty interfaces to store the scanned values.
//   - []any: A slice of pointers to the values, used as scan arguments.
//
// Example:
//
//	values, scanArgs := prepareValuesAndScanArgs(columns)
//	err := rows.Scan(scanArgs...)
//	// After scanning, the values can be accessed and type-asserted as needed.
func prepareValuesAndScanArgs(columns []string) ([]any, []any) {
	values := make([]any, len(columns))
	scanArgs := make([]any, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	return values, scanArgs
}

// convertRowToMap creates a map from the column names and values.
// It converts []byte values to strings.
//
// Parameters:
//   - columns: A slice of strings representing the column names in the result set.
//   - values: A slice of any type representing the values of the columns.
//
// Returns:
//   - map[string]any: A map where the keys are the column names and the values are the corresponding data.
//
// Example:
//
//	columns := []string{"id", "name", "email"}
//	values := []any{1, "John Doe", "john.doe@example.com"}
//	rowMap := convertRowToMap(columns, values)
//	// rowMap will be:
//	// map[id:1 name:John Doe email:john.doe@example.com]
func convertRowToMap(columns []string, values []any) map[string]any {
	rowMap := make(map[string]any, len(columns))
	for i, col := range columns {
		val := values[i]
		if b, ok := val.([]byte); ok {
			rowMap[col] = string(b)
		} else {
			rowMap[col] = val
		}
	}

	return rowMap
}
