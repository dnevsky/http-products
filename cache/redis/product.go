package redis

import (
	"context"
	"encoding/json"
	"errors"

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

func (c *ProductRedis) GetWithOffsetFromJSON(ctx context.Context, offset, limit int) ([]models.Product, error) {
	if offset < 0 || limit <= 0 {
		return nil, errors.New("cache: invalid offset or limit")
	}
	products := make([]models.Product, 0, limit)

	res, err := c.client.LRange(ctx, "products", int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		c.logger.Info(res, err)
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
	if len(data) == 0 {
		return errors.New("UpdateData: there's nothing to stuff")
	}

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
