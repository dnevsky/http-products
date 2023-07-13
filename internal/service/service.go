package service

import (
	"context"

	"github.com/dnevsky/http-products/cache"
	"github.com/dnevsky/http-products/models"
	"go.uber.org/zap"
)

type Product interface {
	GetAll(ctx context.Context, limit, offset int) ([]models.Product, error)
}

type Service struct {
	Product
}

func NewService(logger *zap.SugaredLogger, cache *cache.Cache) *Service {
	return &Service{
		Product: NewProductService(logger, cache),
	}
}
