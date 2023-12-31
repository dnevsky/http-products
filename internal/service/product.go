package service

import (
	"context"

	"github.com/dnevsky/http-products/cache"
	"github.com/dnevsky/http-products/models"
	"go.uber.org/zap"
)

type ProductService struct {
	logger *zap.SugaredLogger
	cache  *cache.Cache
}

func NewProductService(logger *zap.SugaredLogger, cache *cache.Cache) *ProductService {
	return &ProductService{logger: logger, cache: cache}
}

func (s *ProductService) GetAll(ctx context.Context, offset, limit int) ([]models.Product, error) {
	products, err := s.cache.Product.GetWithOffsetFromJSON(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}
