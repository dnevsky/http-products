package postgres

import (
	"github.com/jmoiron/sqlx"
)

const (
	productsTable = "products"
)

func NewPostgresDB(uri string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
