package drivers

import (
	"database/sql"

	"github.com/mohamed-samir907/goquery/query"
)

var _ Driver = (*PostgreSQL)(nil)

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL(db *sql.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (d *PostgreSQL) Get(q query.SelectQuery) ([]map[string]any, error) {
	return nil, nil
}

func (d *PostgreSQL) Insert() {}
func (d *PostgreSQL) Update() {}
func (d *PostgreSQL) Delete() {}
