package repository

import (
	"context"

	"github.com/dnevsky/http-products/internal/repository/postgres"
	"github.com/dnevsky/http-products/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Product interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetAllWithOffset(offset, limit int) ([]models.Product, error)
}

type Repository struct {
	Product
}

func NewRepository(logger *zap.SugaredLogger, db *sqlx.DB) *Repository {
	return &Repository{
		Product: postgres.NewProductPostgres(logger, db),
	}
}
