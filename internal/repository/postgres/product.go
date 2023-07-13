package postgres

import (
	"context"
	"fmt"

	"github.com/dnevsky/http-products/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ProductPostgres struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
}

func NewProductPostgres(logger *zap.SugaredLogger, db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{
		logger: logger,
		db:     db,
	}
}

func (r *ProductPostgres) GetAllWithOffset(limit, offset int) ([]models.Product, error) {
	var products []models.Product

	// Было бы круто использовать курсор (SELECT * FROM %s WHERE id > $1 ORDER BY id ASC LIMIT $2), но со строками так не получится
	// Придется создавать индексы. Не стал этого делать, все равно функцию не вызываем.
	// А зачем тогда она тут? Проверял как база данных будет работать :)
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC OFFSET $1 LIMIT $2", productsTable)

	if err := r.db.Select(&products, query, offset, limit); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductPostgres) GetAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC", productsTable)

	if err := r.db.SelectContext(ctx, &products, query); err != nil {
		return nil, err
	}

	return products, nil
}
