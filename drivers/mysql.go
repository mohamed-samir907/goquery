package drivers

import (
	"database/sql"

	"github.com/mohamed-samir907/goquery/query"
)

var _ Driver = (*MySQL)(nil)

type MySQL struct {
	db *sql.DB
}

func NewMySQL(db *sql.DB) *MySQL {
	return &MySQL{
		db: db,
	}
}

func (d *MySQL) Get(q query.SelectQuery) ([]map[string]any, error) {
	query, args := q.Build()

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values, scanArgs := prepareValuesAndScanArgs(columns)

	var results []map[string]any
	for rows.Next() {
		// scan args represent a pointers to each value of the values slice
		// which means the values slice will be filled by a db row after scan.
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		rowMap := convertRowToMap(columns, values)

		results = append(results, rowMap)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (d *MySQL) Insert() {}
func (d *MySQL) Update() {}
func (d *MySQL) Delete() {}
