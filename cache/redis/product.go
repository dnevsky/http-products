package redis

import (
	"context"
	"encoding/json"

	"github.com/dnevsky/http-products/models"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ProductRedis struct {
	logger *zap.SugaredLogger
	client *redis.Client
}

func NewProductRedis(logger *zap.SugaredLogger, client *redis.Client) *ProductRedis {
	return &ProductRedis{logger: logger, client: client}
}

func (c *ProductRedis) GetWithOffsetFromJSON(ctx context.Context, limit, offset int) ([]models.Product, error) {
	products := make([]models.Product, 0, limit)

	res, err := c.client.LRange(ctx, "products", int64(offset), int64(offset+limit)).Result()
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		var product models.Product

		err := json.Unmarshal([]byte(v), &product)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (c *ProductRedis) UpdateData(ctx context.Context, data []string) error {

	pipe := c.client.TxPipeline()
	defer pipe.Close()

	err := pipe.Del(ctx, "products").Err()
	if err != nil {
		return err
	}

	err = pipe.RPush(ctx, "products", data).Err()
	if err != nil {
		return err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
