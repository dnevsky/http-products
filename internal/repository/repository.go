package repository

import (
	"github.com/dnevsky/http-products/internal/repository/postgres"
	"github.com/dnevsky/http-products/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Product interface {
	GetAll() ([]models.Product, error)
	GetAllWithOffset(limit, offset int) ([]models.Product, error)
}

type Repository struct {
	Product
}

func NewRepository(logger *zap.SugaredLogger, db *sqlx.DB) *Repository {
	return &Repository{
		Product: postgres.NewProductPostgres(logger, db),
	}
}
