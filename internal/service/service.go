package service

import (
	"context"

	"github.com/dnevsky/http-products/cache"
	"github.com/dnevsky/http-products/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Product interface {
	GetAll(ctx context.Context, offset, limit int) ([]models.Product, error)
}

type Service struct {
	Product
}

func NewService(logger *zap.SugaredLogger, cache *cache.Cache) *Service {
	return &Service{
		Product: NewProductService(logger, cache),
	}
}
