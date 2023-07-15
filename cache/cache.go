package cache

import (
	"context"

	red "github.com/dnevsky/http-products/cache/redis"
	"github.com/dnevsky/http-products/models"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

//go:generate mockgen -source=cache.go -destination=mocks/mock.go

type Product interface {
	GetWithOffsetFromJSON(ctx context.Context, offset, limit int) ([]models.Product, error)

	UpdateData(ctx context.Context, data []string) error
}

type Cache struct {
	Product
}

func NewCache(logger *zap.SugaredLogger, client *redis.Client) *Cache {
	return &Cache{
		Product: red.NewProductRedis(logger, client),
	}
}
